package setup

import (
	"context"
	"fmt"
	componentEcoSystem "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	componentHelm "github.com/cloudogu/k8s-component-operator/pkg/helm"
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

const (
	longhornComponentName       = "k8s-longhorn"
	certManagerComponentName    = "k8s-cert-manager"
	certManagerCrdComponentName = "k8s-cert-manager-crd"
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
	helmClient, err := componentHelm.NewClient(e.SetupContext.AppConfig.TargetNamespace, e.SetupContext.HelmRepositoryData, appcontext.IsDevelopmentStage(e.SetupContext.Stage), logrus.StandardLogger().Infof)
	if err != nil {
		return fmt.Errorf("failed to create helm-client: %w", err)
	}

	namespace := e.SetupContext.AppConfig.TargetNamespace

	ecoSystemClient, err := componentEcoSystem.NewForConfig(e.ClusterConfig)
	if err != nil {
		return fmt.Errorf("failed to create K8s Component-EcoSystem client: %w", err)
	}
	componentsClient := ecoSystemClient.Components(namespace)

	certManagerInstallerSteps := e.createCertManagerSteps(helmClient)

	componentOpInstallerSteps, err := e.createComponentOperatorSteps(helmClient, componentsClient)
	if err != nil {
		return err
	}

	longhornComponentSteps := e.createLonghornSteps(componentsClient)

	componentSteps, componentWaitSteps := e.createComponentSteps(componentsClient)

	createNodeMasterStep, err := component.NewNodeMasterCreationStep(e.ClusterConfig, namespace)
	if err != nil {
		return fmt.Errorf("failed to create node master file creation step: %w", err)
	}

	componentResourcePatchStep, err := createResourcePatchStep(patch.ComponentPhase, e.SetupContext.AppConfig.ResourcePatches, e.ClusterConfig, namespace)
	if err != nil {
		return fmt.Errorf("error while creating resource patch step for phase %s: %w", patch.ComponentPhase, err)
	}

	e.RegisterSetupStep(createNodeMasterStep)

	for _, step := range certManagerInstallerSteps {
		e.RegisterSetupStep(step)
	}

	for _, step := range componentOpInstallerSteps {
		e.RegisterSetupStep(step)
	}

	// Install and wait for longhorn before other component installation steps because the component operator can't handle the optional relation between longhorn and e.g. etcd.
	// These steps may be empty if longhorn is not part of the component list.
	for _, step := range longhornComponentSteps {
		e.RegisterSetupStep(step)
	}

	for _, step := range componentSteps {
		e.RegisterSetupStep(step)
	}

	for _, step := range componentWaitSteps {
		e.RegisterSetupStep(step)
	}

	e.RegisterSetupStep(component.NewEtcdClientInstallerStep(e.ClientSet, e.SetupContext))

	// Since this step should patch resources created in this phase, it should be executed last.
	e.RegisterSetupStep(componentResourcePatchStep)

	return nil
}

func (e *Executor) createComponentOperatorSteps(helmClient *componentHelm.Client, componentClient componentEcoSystem.ComponentInterface) ([]ExecutorStep, error) {
	var result []ExecutorStep
	namespace := e.SetupContext.AppConfig.TargetNamespace

	result = append(result, component.NewCesComponentChartInstallerStep(namespace, e.SetupContext.AppConfig.ComponentOperatorCrdChart, helmClient))
	result = append(result, component.NewCesComponentChartInstallerStep(namespace, e.SetupContext.AppConfig.ComponentOperatorChart, helmClient))
	operatorComponentSteps, err := e.appendComponentStepsForComponentOperator(componentClient)
	if err != nil {
		return nil, err
	}
	result = append(result, operatorComponentSteps...)

	return result, nil
}

func (e *Executor) appendComponentStepsForComponentOperator(componentClient componentEcoSystem.ComponentInterface) ([]ExecutorStep, error) {
	var result []ExecutorStep
	namespace := e.SetupContext.AppConfig.TargetNamespace

	stepsCrdChart, err := e.createComponentStepsByString(componentClient, e.SetupContext.AppConfig.ComponentOperatorCrdChart, namespace)
	if err != nil {
		return nil, err
	}

	stepsChart, err := e.createComponentStepsByString(componentClient, e.SetupContext.AppConfig.ComponentOperatorChart, namespace)
	if err != nil {
		return nil, err
	}

	result = append(result, stepsCrdChart...)
	result = append(result, stepsChart...)

	return result, nil
}

func (e *Executor) createComponentStepsByString(componentClient componentEcoSystem.ComponentInterface, chartStr string, namespace string) ([]ExecutorStep, error) {
	var result []ExecutorStep

	fullChartName, chartVersion, err := component.SplitChartString(chartStr)
	if err != nil {
		return nil, fmt.Errorf("failed to split chart string %s: %w", chartStr, err)
	}
	helmNamespace, name := component.SplitHelmNamespaceFromChartString(fullChartName)

	result = append(result, component.NewInstallComponentsStep(componentClient, name, helmNamespace, chartVersion, namespace, namespace))
	result = append(result, component.NewWaitForComponentStep(componentClient, createComponentLabelSelector(name), namespace, component.DefaultComponentWaitTimeOut5Minutes))

	return result, nil
}

func (e *Executor) createCertManagerSteps(helmClient *componentHelm.Client) []ExecutorStep {
	var result []ExecutorStep

	result = e.appendComponentChartInstallStepToSlice(helmClient, certManagerCrdComponentName, result)
	result = e.appendComponentChartInstallStepToSlice(helmClient, certManagerComponentName, result)

	return result
}

func (e *Executor) appendComponentChartInstallStepToSlice(helmClient *componentHelm.Client, name string, steps []ExecutorStep) []ExecutorStep {
	if step := e.createComponentChartInstallStepFromComponentList(name, helmClient); step != nil {
		steps = append(steps, step)
	}

	return steps
}

func (e *Executor) createComponentChartInstallStepFromComponentList(name string, helmClient *componentHelm.Client) ExecutorStep {
	components := e.SetupContext.AppConfig.Components

	if c, containsComponentChart := components[name]; containsComponentChart {
		namespace := e.SetupContext.AppConfig.TargetNamespace
		if c.DeployNamespace != "" {
			namespace = c.DeployNamespace
		}
		chartUrl := fmt.Sprintf("%s/%s:%s", c.HelmRepositoryNamespace, name, c.Version)

		return component.NewCesComponentChartInstallerStep(namespace, chartUrl, helmClient)
	}

	return nil
}

func (e *Executor) createLonghornSteps(componentsClient componentEcoSystem.ComponentInterface) []ExecutorStep {
	var result []ExecutorStep
	components := e.SetupContext.AppConfig.Components
	namespace := e.SetupContext.AppConfig.TargetNamespace

	_, containsLonghorn := components[longhornComponentName]

	if containsLonghorn {
		helmRepoNamespace := components[longhornComponentName].HelmRepositoryNamespace
		version := components[longhornComponentName].Version
		deployNamespace := components[longhornComponentName].DeployNamespace
		installStep := component.NewInstallComponentsStep(componentsClient, longhornComponentName, helmRepoNamespace, version, namespace, deployNamespace)
		selector := createComponentLabelSelector(longhornComponentName)
		waitStep := component.NewWaitForComponentStep(componentsClient, selector, namespace, component.DefaultComponentWaitTimeOut5Minutes)
		result = append(result, installStep)
		result = append(result, waitStep)
		delete(components, longhornComponentName)
	}

	return result
}

func createComponentLabelSelector(name string) string {
	return fmt.Sprintf("%s=%s", v1LabelK8sComponent, name)
}

func (e *Executor) createComponentSteps(componentsClient componentEcoSystem.ComponentInterface) ([]ExecutorStep, []ExecutorStep) {
	namespace := e.SetupContext.AppConfig.TargetNamespace
	var componentSteps []ExecutorStep
	var waitSteps []ExecutorStep

	for componentName, componentAttributes := range e.SetupContext.AppConfig.Components {
		componentSteps = append(componentSteps, component.NewInstallComponentsStep(componentsClient, componentName, componentAttributes.HelmRepositoryNamespace, componentAttributes.Version, namespace, componentAttributes.DeployNamespace))
		waitSteps = append(waitSteps, component.NewWaitForComponentStep(componentsClient, createComponentLabelSelector(componentName), namespace, component.DefaultComponentWaitTimeOut5Minutes))
	}

	return componentSteps, waitSteps
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
