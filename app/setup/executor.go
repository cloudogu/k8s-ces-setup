package setup

import (
	"fmt"
	"time"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/cloudogu/k8s-ces-setup/app/setup/data"

	"k8s.io/client-go/rest"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

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
}

// NewExecutor creates a new setup executor with the given app configuration.
func NewExecutor(clusterConfig *rest.Config, context *context.SetupContext) (*Executor, error) {
	k8sClient, err := createKubernetesClient(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Executor{
		SetupContext:  context,
		ClientSet:     k8sClient,
		ClusterConfig: clusterConfig,
	}, nil
}

func createKubernetesClient(clusterConfig *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	return clientSet, nil
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
	etcdSrvInstallerStep, err := component.NewEtcdServerInstallerStep(e.ClusterConfig, e.SetupContext)
	if err != nil {
		return fmt.Errorf("failed to create new etcd server installer step: %w", err)
	}

	doguOpInstallerStep, err := component.NewDoguOperatorInstallerStep(e.ClusterConfig, e.SetupContext)
	if err != nil {
		return fmt.Errorf("failed to create new dogu operator installer step: %w", err)
	}

	serviceDisInstallerStep, err := component.NewServiceDiscoveryInstallerStep(e.ClusterConfig, e.SetupContext)
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
func (e *Executor) RegisterDataSetupSteps() error {
	// note: with introduction of the setup UI the instance secret may either come into play with a new instance
	// registration or it may already reside in the current namespace
	namespace := e.SetupContext.AppConfig.TargetNamespace
	e.RegisterSetupStep(data.NewInstanceSecretValidatorStep(e.ClientSet, namespace))

	etcdRegistry, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", namespace)},
	})

	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	configurationSetupStep := data.NewWriteConfigToRegistryStep(etcdRegistry, &e.SetupContext.StartupConfiguration)
	e.RegisterSetupStep(data.NewKeyProviderStep(etcdRegistry.GlobalConfig()))
	e.RegisterSetupStep(data.NeWriteSSL(e.SetupContext.StartupConfiguration, etcdRegistry.GlobalConfig()))
	e.RegisterSetupStep(configurationSetupStep)

	return nil
}

func (e *Executor) RegisterValidationStep() error {
	setupValidator, err := NewValidatorStep(e.ClientSet, e.SetupContext)
	if err != nil {
		return fmt.Errorf("failed to create new validation step: %w", err)
	}

	e.RegisterSetupStep(setupValidator)
	return nil
}

func (e *Executor) RegisterSSLGenerationStep() error {
	generationStep := data.NewGenerateSSL(&e.SetupContext.StartupConfiguration)
	e.RegisterSetupStep(generationStep)

	return nil
}
