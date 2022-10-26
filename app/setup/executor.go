package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-apply-lib/apply"
	k8sv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"strings"
	"time"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/remote"

	"github.com/cloudogu/k8s-ces-setup/app/setup/data"

	"k8s.io/client-go/rest"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

const K8sSetupFieldManagerName = "k8s-ces-setup"

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
	endpoint = strings.TrimSuffix(endpoint, "dogus")
	endpoint = strings.TrimSuffix(endpoint, "/")

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

// RegisterComponentSetupSteps adds all setups steps responsible to install vital components into the ecosystem.
func (e *Executor) RegisterComponentSetupSteps() error {
	k8sApplyClient, scheme, err := apply.New(e.ClusterConfig, K8sSetupFieldManagerName)
	if err != nil {
		return fmt.Errorf("failed to create k8s apply client: %w", err)
	}
	err = k8sv1.AddToScheme(scheme)
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

	e.RegisterSetupStep(createNodeMasterStep)
	e.RegisterSetupStep(etcdSrvInstallerStep)
	e.RegisterSetupStep(component.NewWaitForPodStep(e.ClientSet, "statefulset.kubernetes.io/pod-name=etcd-0", namespace, time.Second*300))
	e.RegisterSetupStep(component.NewEtcdClientInstallerStep(e.ClientSet, e.SetupContext))
	e.RegisterSetupStep(doguOpInstallerStep)
	e.RegisterSetupStep(serviceDisInstallerStep)

	return nil
}

// RegisterDataSetupSteps adds all setups steps responsible to read, write, or verify data needed by the setup.
func (e *Executor) RegisterDataSetupSteps(etcdRegistry registry.Registry) error {
	configWriter := data.NewRegistryConfigurationWriter(etcdRegistry)

	// register steps
	e.RegisterSetupStep(data.NewKeyProviderStep(configWriter, e.SetupContext.AppConfig.KeyProvider))
	e.RegisterSetupStep(data.NewInstanceSecretValidatorStep(e.ClientSet, e.SetupContext.AppConfig.TargetNamespace))
	e.RegisterSetupStep(data.NewWriteAdminDataStep(configWriter, e.SetupContext.StartupConfiguration))
	e.RegisterSetupStep(data.NewWriteNamingDataStep(configWriter, e.SetupContext.StartupConfiguration))
	e.RegisterSetupStep(data.NewWriteRegistryConfigEncryptedStep(e.SetupContext.StartupConfiguration, e.ClientSet, e.SetupContext.AppConfig.TargetNamespace))
	e.RegisterSetupStep(data.NewWriteLdapDataStep(configWriter, e.SetupContext.StartupConfiguration))
	e.RegisterSetupStep(data.NewWriteRegistryConfigDataStep(configWriter, e.SetupContext.StartupConfiguration))
	e.RegisterSetupStep(data.NewWriteDoguDataStep(configWriter, e.SetupContext.StartupConfiguration))

	return nil
}

// RegisterDoguInstallationSteps creates install steps for the dogu install list
func (e *Executor) RegisterDoguInstallationSteps() error {
	doguStepGenerator, err := NewDoguStepGenerator(e.ClientSet, e.ClusterConfig, e.SetupContext.StartupConfiguration.Dogus, e.Registry, e.SetupContext.AppConfig.TargetNamespace)
	if err != nil {
		return fmt.Errorf("failed to generate dogu step generator: %w", err)
	}

	doguSteps := doguStepGenerator.GenerateSteps()
	for _, step := range doguSteps {
		e.RegisterSetupStep(step)
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
	generationStep := data.NewGenerateSSLStep(e.SetupContext.StartupConfiguration)
	e.RegisterSetupStep(generationStep)
	return nil
}
