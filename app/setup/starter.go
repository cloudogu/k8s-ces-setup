package setup

import (
	gocontext "context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/cesregistry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

// SetupExecutor is uses to register all necessary steps and executes them
type SetupExecutor interface {
	// RegisterFQDNRetrieverStep registers the FQDN retriever step
	RegisterFQDNRetrieverStep()
	// RegisterSSLGenerationStep registers all ssl steps
	RegisterSSLGenerationStep() error
	// RegisterValidationStep registers all validation steps
	RegisterValidationStep() error
	// RegisterComponentSetupSteps adds all setup steps responsible to install vital components into the ecosystem.
	RegisterComponentSetupSteps() error
	// RegisterDataSetupSteps adds all setup steps responsible to read, write, or verify data needed by the setup.
	RegisterDataSetupSteps(registry.Registry) error
	// RegisterDoguInstallationSteps creates install steps for the dogu install list
	RegisterDoguInstallationSteps() error
	// PerformSetup starts the setup and executes all registered setup steps
	PerformSetup() (error, string)
}

// Starter is used to init and start the setup process
type Starter struct {
	EtcdRegistry  registry.Registry
	ClientSet     kubernetes.Interface
	ClusterConfig *rest.Config
	SetupContext  *context.SetupContext
	Namespace     string
	SetupExecutor SetupExecutor
}

// NewStarter creates a new setup starter struct which one inits registries and starts the setup process
func NewStarter(clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *context.SetupContextBuilder) (*Starter, error) {
	setupContext, err := setupContextBuilder.NewSetupContext(k8sClient)
	if err != nil {
		return nil, err
	}

	namespace := setupContext.AppConfig.TargetNamespace
	etcdRegistry, err := cesregistry.CreateEtcd(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	setupExecutor, err := NewExecutor(clusterConfig, k8sClient, setupContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create setup executor: %w", err)
	}

	return &Starter{
		EtcdRegistry:  etcdRegistry,
		ClientSet:     k8sClient,
		ClusterConfig: clusterConfig,
		SetupContext:  setupContext,
		Namespace:     namespace,
		SetupExecutor: setupExecutor,
	}, nil
}

// StartSetup creates necessary k8s config and client, register steps and executes them
func (s *Starter) StartSetup() error {
	err := setSetupState(s.ClientSet, s.Namespace, context.SetupStateInstalling)
	if err != nil {
		return err
	}

	err = registerSteps(s.SetupExecutor, s.EtcdRegistry, s.SetupContext)
	if err != nil {
		return err
	}

	err, errCausingAction := s.SetupExecutor.PerformSetup()
	if err != nil {
		return fmt.Errorf("error while performing setup [%s]: %w", errCausingAction, err)
	}

	err = setSetupState(s.ClientSet, s.Namespace, context.SetupStateInstalled)
	if err != nil {
		return err
	}

	return nil
}

func registerSteps(setupExecutor SetupExecutor, etcdRegistry registry.Registry, setupContext *context.SetupContext) error {
	if setupContext.SetupJsonConfiguration.Naming.Fqdn == "" {
		setupExecutor.RegisterFQDNRetrieverStep()
	}

	if setupContext.SetupJsonConfiguration.Naming.CertificateType == "selfsigned" {
		err := setupExecutor.RegisterSSLGenerationStep()
		if err != nil {
			return fmt.Errorf("failed to register ssl generation setup step: %w", err)
		}
	}

	err := setupExecutor.RegisterValidationStep()
	if err != nil {
		return fmt.Errorf("failed to register validation setup steps: %w", err)
	}

	err = setupExecutor.RegisterComponentSetupSteps()
	if err != nil {
		return fmt.Errorf("failed to register component setup steps: %w", err)
	}

	err = setupExecutor.RegisterDataSetupSteps(etcdRegistry)
	if err != nil {
		return fmt.Errorf("failed to register data setup steps: %w", err)
	}

	err = setupExecutor.RegisterDoguInstallationSteps()
	if err != nil {
		return fmt.Errorf("failed to register dogu installation steps: %w", err)
	}

	return nil
}

func setSetupState(clientSet kubernetes.Interface, namespace string, state string) error {
	stateCM, err := context.GetSetupStateConfigMap(clientSet, namespace)
	if err != nil {
		return fmt.Errorf("failed to get k8s-ces-setup configmap: %w", err)
	}

	if state == context.SetupStateInstalling {
		actualState := stateCM.Data[context.SetupStateKey]
		if actualState == context.SetupStateInstalling || actualState == context.SetupStateInstalled {
			return fmt.Errorf("setup is busy or already done")
		}
	}

	stateCM.Data[context.SetupStateKey] = state
	_, err = clientSet.CoreV1().ConfigMaps(namespace).Update(gocontext.Background(), stateCM, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update k8s-ces-setup configmap: %w", err)
	}

	return nil
}
