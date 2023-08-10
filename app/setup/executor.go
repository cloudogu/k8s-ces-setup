package setup

import (
	"fmt"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/remote"

	"k8s.io/client-go/rest"

	"github.com/cloudogu/k8s-apply-lib/apply"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/patch"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	k8sv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/scheme"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

const k8sSetupFieldManagerName = "k8s-ces-setup"

// ExecutorStep describes a valid step in the setup.
type ExecutorStep interface {
	// GetStepDescription returns the description of the setup step. The Executor prints the description of every step
	// when executing the setup.
	GetStepDescription() string
	// PerformSetupStep is called for every registered step when executing the setup.
	PerformSetupStep() error
}

// Executor is responsible to perform the actual steps of the setup.
type Executor struct {
	// SetupContext contains information about the current context.
	SetupContext *context.SetupContext `json:"setup_context"`
	// ClientSet is the actual k8s client responsible for the k8s API communication.
	ClientSet kubernetes.Interface `json:"client_set"`
	// ClusterConfig is the current rest config used to create a kubernetes clients.
	ClusterConfig *rest.Config `json:"cluster_config"`
	// Steps contains all necessary steps for the setup
	Steps []ExecutorStep `json:"steps"`
	// Registry is the dogu registry
	Registry remote.Registry `json:"registry"`
}

// NewExecutor creates a new setup executor with the given app configuration.
func NewExecutor(clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupCtx *context.SetupContext) (*Executor, error) {
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
func (e *Executor) PerformSetup() (err error, errCausingAction string) {
	logrus.Print("Starting the setup process")

	for _, step := range e.Steps {
		logrus.Printf("Setup-Step: %s", step.GetStepDescription())

		err := step.PerformSetupStep()
		if err != nil {
			return fmt.Errorf("failed to perform step [%s]: %w", step.GetStepDescription(), err), step.GetStepDescription()
		}
	}

	return nil, ""
}

// RegisterComponentSetupSteps adds all setup steps responsible to install vital components into the ecosystem.
func (e *Executor) RegisterComponentSetupSteps() error {
	k8sApplyClient, scheme1, err := apply.New(e.ClusterConfig, k8sSetupFieldManagerName)
	if err != nil {
		return fmt.Errorf("failed to create k8s apply client: %w", err)
	}
	err = k8sv1.AddToScheme(scheme1)
	if err != nil {
		return fmt.Errorf("failed add applier scheme to dogu CRD scheme handling: %w", err)
	}

	etcdSrvInstallerStep, err := component.NewEtcdServerInstallerStep(e.SetupContext, k8sApplyClient)
	if err != nil {
		return fmt.Errorf("failed to create new etcd server installer step: %w", err)
	}

	doguOpInstallerStep, err := component.NewDoguOperatorInstallerStep(e.SetupContext, k8sApplyClient)
	if err != nil {
		return fmt.Errorf("failed to create new dogu operator installer step: %w", err)
	}

	serviceDisInstallerStep, err := component.NewServiceDiscoveryInstallerStep(e.SetupContext, k8sApplyClient)
	if err != nil {
		return fmt.Errorf("failed to create new service discovery installer step: %w", err)
	}

	namespace := e.SetupContext.AppConfig.TargetNamespace
	createNodeMasterStep, err := component.NewNodeMasterCreationStep(e.ClusterConfig, namespace)
	if err != nil {
		return fmt.Errorf("failed to create node master file creation step: %w", err)
	}

	componentResourcePatchStep, err := createResourcePatchStep(patch.ComponentPhase, e.SetupContext.AppConfig.ResourcePatches, e.ClusterConfig, e.SetupContext.AppConfig.TargetNamespace)
	if err != nil {
		return fmt.Errorf("failed to create resource patch step for phase %s: %w", patch.ComponentPhase, err)
	}

	e.RegisterSetupStep(createNodeMasterStep)
	e.RegisterSetupStep(etcdSrvInstallerStep)
	e.RegisterSetupStep(component.NewWaitForPodStep(e.ClientSet, "statefulset.kubernetes.io/pod-name=etcd-0", namespace, component.PodTimeoutInSeconds()))
	e.RegisterSetupStep(component.NewEtcdClientInstallerStep(e.ClientSet, e.SetupContext))
	e.RegisterSetupStep(doguOpInstallerStep)
	e.RegisterSetupStep(serviceDisInstallerStep)
	// Since this step should patch resources created in this phase, it should be executed last.
	e.RegisterSetupStep(componentResourcePatchStep)

	return nil
}

func createResourcePatchStep(phase patch.Phase, patches []patch.ResourcePatch, clusterConfig *rest.Config, targetNamespace string) (*resourcePatchStep, error) {
	resourcePatchApplier, err := patch.NewApplier(clusterConfig, []scheme.Builder{*k8sv1.SchemeBuilder}, targetNamespace)
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

// RegisterFQDNRetrieverStep registers the steps for retrieving the fqdn
func (e *Executor) RegisterFQDNRetrieverStep() {
	namespace := e.SetupContext.AppConfig.TargetNamespace
	config := e.SetupContext.SetupJsonConfiguration
	e.RegisterSetupStep(data.NewFQDNRetrieverStep(config, e.ClientSet, namespace))
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
