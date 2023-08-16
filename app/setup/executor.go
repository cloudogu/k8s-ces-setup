package setup

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/helm"
	componentEcoSystem "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/remote"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/patch"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ExecutorStep describes a valid step in the setup.
type ExecutorStep interface {
	// GetStepDescription returns the description of the setup step. The Executor prints the description of every step
	// when executing the setup.
	GetStepDescription() string
	// PerformSetupStep is called for every registered step when executing the setup.
	PerformSetupStep(ctx context.Context) error
}

// Executor is responsible to perform the actual steps of the setup.
type Executor struct {
	// SetupContext contains information about the current context.
	SetupContext *appcontext.SetupContext
	// ClientSet is the actual k8s client responsible for the k8s API communication.
	ClientSet kubernetes.Interface
	// ClusterConfig is the current rest config used to create a kubernetes clients.
	ClusterConfig *rest.Config
	// Steps contains all necessary steps for the setup
	Steps []ExecutorStep
	// Registry is the dogu registry
	Registry remote.Registry
}

// NewExecutor creates a new setup executor with the given app configuration.
func NewExecutor(clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupCtx *appcontext.SetupContext) (*Executor, error) {
	credentials := &core.Credentials{
		Username: setupCtx.DoguRegistryConfiguration.Username,
		Password: setupCtx.DoguRegistryConfiguration.Password,
	}

	doguRegistry, err := remote.New(getRemoteConfig(setupCtx.DoguRegistryConfiguration.Endpoint, setupCtx.DoguRegistryConfiguration.URLSchema), credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote Registry: %w", err)
	}

	return &Executor{
		SetupContext:  setupCtx,
		ClientSet:     k8sClient,
		ClusterConfig: clusterConfig,
		Registry:      doguRegistry,
	}, nil
}

func getRemoteConfig(endpoint string, urlSchema string) *core.Remote {
	endpoint = strings.TrimSuffix(endpoint, "/")
	if urlSchema == "default" {
		endpoint = strings.TrimSuffix(endpoint, "dogus")
		endpoint = strings.TrimSuffix(endpoint, "/")
	}

	return &core.Remote{
		Endpoint:  endpoint,
		URLSchema: urlSchema,
		CacheDir:  "/tmp",
	}
}

// RegisterSetupStep adds a new step to the setup
func (e *Executor) RegisterSetupStep(step ExecutorStep) {
	logrus.Debugf("Register setup step [%s]", step.GetStepDescription())
	e.Steps = append(e.Steps, step)
}

// PerformSetup starts the setup and executes all registered setup steps
func (e *Executor) PerformSetup(ctx context.Context) (err error, errCausingAction string) {
	logrus.Print("Starting the setup process")

	for _, step := range e.Steps {
		logrus.Printf("Setup-Step: %s", step.GetStepDescription())

		err := step.PerformSetupStep(ctx)
		if err != nil {
			return fmt.Errorf("failed to perform step [%s]: %w", step.GetStepDescription(), err), step.GetStepDescription()
		}
	}

	return nil, ""
}

// RegisterComponentSetupSteps adds all setup steps responsible to install vital components into the ecosystem.
func (e *Executor) RegisterComponentSetupSteps() error {
	ociEndPoint, err := e.SetupContext.HelmRepositoryData.GetOciEndpoint()
	if err != nil {
		return fmt.Errorf("failed get OCI-endpoint of helm-repo: %w", err)
	}
	helmClient, err := helm.NewClient(e.SetupContext.AppConfig.TargetNamespace, ociEndPoint, appcontext.IsDevelopmentStage(e.SetupContext.Stage), logrus.StandardLogger().Infof)
	if err != nil {
		return fmt.Errorf("failed to create helm-client: %w", err)
	}

	namespace := e.SetupContext.AppConfig.TargetNamespace

	componentOpInstallerStep := component.NewComponentOperatorInstallerStep(e.SetupContext, helmClient)

	// create component steps
	componentSteps, componentWaitSteps, err := e.createComponentSteps(namespace)
	if err != nil {
		return fmt.Errorf("failed to create component-steps: %w", err)
	}

	createNodeMasterStep, err := component.NewNodeMasterCreationStep(e.ClusterConfig, namespace)
	if err != nil {
		return fmt.Errorf("failed to create node master file creation step: %w", err)
	}

	componentResourcePatchStep, err := createResourcePatchStep(patch.ComponentPhase, e.SetupContext.AppConfig.ResourcePatches, e.ClusterConfig, namespace)
	if err != nil {
		return fmt.Errorf("error while creating resource patch step for phase %s: %w", patch.ComponentPhase, err)
	}

	// Register steps
	e.RegisterSetupStep(createNodeMasterStep)
	e.RegisterSetupStep(componentOpInstallerStep)

	// install components
	for _, step := range componentSteps {
		e.RegisterSetupStep(step)
	}

	// wait for components to be installed
	for _, step := range componentWaitSteps {
		e.RegisterSetupStep(step)
	}

	e.RegisterSetupStep(component.NewEtcdClientInstallerStep(e.ClientSet, e.SetupContext))

	// Since this step should patch resources created in this phase, it should be executed last.
	e.RegisterSetupStep(componentResourcePatchStep)

	return nil
}

func (e *Executor) createComponentSteps(namespace string) (componentSteps []ExecutorStep, waitSteps []ExecutorStep, err error) {
	ecoSystemClient, err := componentEcoSystem.NewForConfig(e.ClusterConfig)
	if err != nil {
		return componentSteps, waitSteps, fmt.Errorf("failed to create K8s Component-EcoSystem client: %w", err)
	}
	componentsClient := ecoSystemClient.Components(namespace)

	componentSteps = make([]ExecutorStep, len(e.SetupContext.AppConfig.Components))
	waitSteps = make([]ExecutorStep, len(e.SetupContext.AppConfig.Components))
	index := 0
	for fullComponentName, version := range e.SetupContext.AppConfig.Components {
		componentName := fullComponentName[strings.LastIndex(fullComponentName, "/")+1:]
		componentNameSpace := fullComponentName[:strings.LastIndex(fullComponentName, "/")]

		componentSteps[index] = component.NewInstallComponentsStep(componentsClient, componentName, componentNameSpace, version, namespace)

		labelSelector := fmt.Sprintf("%s=%s", v1LabelK8sComponent, componentName)
		waitSteps[index] = component.NewWaitForComponentStep(componentsClient, labelSelector, namespace, component.DefaultComponentWaitTimeOut5Minutes)
		index++
	}

	return
}

func createResourcePatchStep(phase patch.Phase, patches []patch.ResourcePatch, clusterConfig *rest.Config, targetNamespace string) (*resourcePatchStep, error) {
	resourcePatchApplier, err := patch.NewApplier(clusterConfig, targetNamespace)
	if err != nil {
		return nil, err
	}

	resourcePatcher := patch.NewResourcePatcher(resourcePatchApplier)
	componentResourcePatchStep := NewResourcePatchStep(phase, resourcePatcher, patches)
	return componentResourcePatchStep, nil
}

// RegisterDataSetupSteps adds all setup steps responsible to read, write, or verify data needed by the setup.
func (e *Executor) RegisterDataSetupSteps(etcdRegistry registry.Registry) error {
	configWriter := data.NewRegistryConfigurationWriter(etcdRegistry)

	// register steps
	e.RegisterSetupStep(data.NewKeyProviderStep(configWriter, e.SetupContext.AppConfig.KeyProvider))
	e.RegisterSetupStep(data.NewInstanceSecretValidatorStep(e.ClientSet, e.SetupContext.AppConfig.TargetNamespace))
	e.RegisterSetupStep(data.NewWriteAdminDataStep(configWriter, e.SetupContext.SetupJsonConfiguration))
	e.RegisterSetupStep(data.NewWriteNamingDataStep(configWriter, e.SetupContext.SetupJsonConfiguration, e.ClientSet, e.SetupContext.AppConfig.TargetNamespace))
	e.RegisterSetupStep(data.NewWriteRegistryConfigEncryptedStep(e.SetupContext.SetupJsonConfiguration, e.ClientSet, e.SetupContext.AppConfig.TargetNamespace))
	e.RegisterSetupStep(data.NewWriteLdapDataStep(configWriter, e.SetupContext.SetupJsonConfiguration))
	e.RegisterSetupStep(data.NewWriteRegistryConfigDataStep(configWriter, e.SetupContext.SetupJsonConfiguration))
	e.RegisterSetupStep(data.NewWriteDoguDataStep(configWriter, e.SetupContext.SetupJsonConfiguration))

	return nil
}

// RegisterDoguInstallationSteps creates install steps for the dogu install list
func (e *Executor) RegisterDoguInstallationSteps() error {
	doguStepGenerator, err := NewDoguStepGenerator(e.ClientSet, e.ClusterConfig, e.SetupContext.SetupJsonConfiguration.Dogus, e.Registry, e.SetupContext.AppConfig.TargetNamespace)
	if err != nil {
		return fmt.Errorf("failed to generate dogu step generator: %w", err)
	}

	doguSteps, err := doguStepGenerator.GenerateSteps()
	if err != nil {
		return fmt.Errorf("could not register installation steps: %w", err)
	}

	for _, step := range doguSteps {
		e.RegisterSetupStep(step)
	}

	doguResourcePatchStep, err := createResourcePatchStep(patch.DoguPhase, e.SetupContext.AppConfig.ResourcePatches, e.ClusterConfig, e.SetupContext.AppConfig.TargetNamespace)
	if err != nil {
		return fmt.Errorf("failed to create resource patch step for phase %s: %w", patch.DoguPhase, err)
	}

	// Since this step should patch resources created in this phase, it should be executed last.
	e.RegisterSetupStep(doguResourcePatchStep)

	return nil
}

// RegisterLoadBalancerFQDNRetrieverSteps registers the steps for creating a loadbalancer retrieving the fqdn
func (e *Executor) RegisterLoadBalancerFQDNRetrieverSteps() error {
	namespace := e.SetupContext.AppConfig.TargetNamespace
	config := e.SetupContext.SetupJsonConfiguration
	e.RegisterSetupStep(data.NewCreateLoadBalancerStep(config, e.ClientSet, namespace))

	loadbalancerResourcePatchStep, err := createResourcePatchStep(
		patch.LoadbalancerPhase,
		e.SetupContext.AppConfig.ResourcePatches,
		e.ClusterConfig,
		namespace,
	)
	if err != nil {
		return fmt.Errorf("failed to create resource patch step for phase %s: %w", patch.LoadbalancerPhase, err)
	}

	// Since this step should patch resources created in this phase, it should be executed after creating the loadbalancer.
	e.RegisterSetupStep(loadbalancerResourcePatchStep)

	wantsLoadbalancerIpAddressAsFqdn := config.Naming.Fqdn == "" || config.Naming.Fqdn == "<<ip>>"
	if wantsLoadbalancerIpAddressAsFqdn {
		// Here we wait for an external IP address automagically or (after introducing the above patch) an internal IP address.
		// We ignore the case where the public IP address was already assigned but the patch should lead to another.
		e.RegisterSetupStep(data.NewFQDNRetrieverStep(config, e.ClientSet, namespace))
	}

	return nil
}

// RegisterValidationStep registers all validation steps
func (e *Executor) RegisterValidationStep() error {
	e.RegisterSetupStep(NewValidatorStep(e.Registry, e.SetupContext))
	return nil
}

// RegisterSSLGenerationStep registers all ssl steps
func (e *Executor) RegisterSSLGenerationStep() error {
	generationStep := data.NewGenerateSSLStep(e.SetupContext.SetupJsonConfiguration)
	e.RegisterSetupStep(generationStep)
	return nil
}
