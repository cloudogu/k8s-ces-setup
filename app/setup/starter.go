package setup

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	k8sreg "github.com/cloudogu/k8s-registry-lib/registry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
)

// SetupExecutor is uses to register all necessary steps and executes them
type SetupExecutor interface {
	// RegisterLoadBalancerFQDNRetrieverSteps registers the FQDN retriever step
	RegisterLoadBalancerFQDNRetrieverSteps() error
	// RegisterSSLGenerationStep registers all ssl steps
	RegisterSSLGenerationStep() error
	// RegisterValidationStep registers all validation steps
	RegisterValidationStep() error
	// RegisterComponentSetupSteps adds all setup steps responsible to install vital components into the ecosystem.
	RegisterComponentSetupSteps() error
	// RegisterDataSetupSteps adds all setup steps responsible to read, write, or verify data needed by the setup.
	RegisterDataSetupSteps(globalConfig data.ConfigurationRegistry, doguConfigProvider data.DoguConfigurationRegistryProvider[data.ConfigurationRegistry]) error
	// RegisterDoguInstallationSteps creates install steps for the dogu install list
	RegisterDoguInstallationSteps() error
	// PerformSetup starts the setup and executes all registered setup steps
	PerformSetup(ctx context.Context) (error, string)
}

// Starter is used to init and start the setup process
type Starter struct {
	globalRegistry data.ConfigurationRegistry
	doguReg        data.DoguConfigurationRegistryProvider[data.ConfigurationRegistry]
	ClientSet      kubernetes.Interface
	ClusterConfig  *rest.Config
	SetupContext   *appcontext.SetupContext
	Namespace      string
	SetupExecutor  SetupExecutor
}

// NewStarter creates a new setup starter struct which one inits registries and starts the setup process
func NewStarter(ctx context.Context, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *appcontext.SetupContextBuilder) (*Starter, error) {
	setupContext, err := setupContextBuilder.NewSetupContext(ctx, k8sClient)
	if err != nil {
		return nil, err
	}

	namespace := setupContext.AppConfig.TargetNamespace
	cmClient := k8sClient.CoreV1().ConfigMaps(namespace)
	doguReg := data.NewDoguConfigRegistryProvider[data.ConfigurationRegistry](cmClient)
	globalReg, err := k8sreg.NewGlobalConfigRegistry(ctx, cmClient)
	if err != nil {
		return nil, fmt.Errorf("failed to setup global registry: %w", err)
	}

	setupExecutor, err := NewExecutor(clusterConfig, k8sClient, setupContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create setup executor: %w", err)
	}

	return &Starter{
		globalRegistry: globalReg,
		doguReg:        doguReg,
		ClientSet:      k8sClient,
		ClusterConfig:  clusterConfig,
		SetupContext:   setupContext,
		Namespace:      namespace,
		SetupExecutor:  setupExecutor,
	}, nil
}

// StartSetup creates necessary k8s config and client, register steps and executes them
func (s *Starter) StartSetup(ctx context.Context) error {
	err := setSetupState(ctx, s.ClientSet, s.Namespace, appcontext.SetupStateInstalling)
	if err != nil {
		return err
	}

	err = registerSteps(s.SetupExecutor, s.globalRegistry, s.doguReg, s.SetupContext)
	if err != nil {
		return err
	}

	err, errCausingAction := s.SetupExecutor.PerformSetup(ctx)
	if err != nil {
		return fmt.Errorf("error while performing setup [%s]: %w", errCausingAction, err)
	}

	err = setSetupState(ctx, s.ClientSet, s.Namespace, appcontext.SetupStateInstalled)
	if err != nil {
		return err
	}

	return nil
}

func registerSteps(setupExecutor SetupExecutor, globalConfig data.ConfigurationRegistry, doguConfigProvider data.DoguConfigurationRegistryProvider[data.ConfigurationRegistry], setupContext *appcontext.SetupContext) error {
	err := setupExecutor.RegisterLoadBalancerFQDNRetrieverSteps()
	if err != nil {
		return fmt.Errorf("failed to register steps for creating loadbalancer and retrieving its ip as fqdn: %w", err)
	}

	if setupContext.SetupJsonConfiguration.Naming.CertificateType == "selfsigned" {
		err := setupExecutor.RegisterSSLGenerationStep()
		if err != nil {
			return fmt.Errorf("failed to register ssl generation setup step: %w", err)
		}
	}

	err = setupExecutor.RegisterValidationStep()
	if err != nil {
		return fmt.Errorf("failed to register validation setup steps: %w", err)
	}

	err = setupExecutor.RegisterComponentSetupSteps()
	if err != nil {
		return fmt.Errorf("failed to register component setup steps: %w", err)
	}

	err = setupExecutor.RegisterDataSetupSteps(globalConfig, doguConfigProvider)
	if err != nil {
		return fmt.Errorf("failed to register data setup steps: %w", err)
	}

	err = setupExecutor.RegisterDoguInstallationSteps()
	if err != nil {
		return fmt.Errorf("failed to register dogu installation steps: %w", err)
	}

	return nil
}

func setSetupState(ctx context.Context, clientSet kubernetes.Interface, namespace string, state string) error {
	stateCM, err := appcontext.GetSetupStateConfigMap(ctx, clientSet, namespace)
	if err != nil {
		return fmt.Errorf("failed to get k8s-ces-setup configmap: %w", err)
	}

	if state == appcontext.SetupStateInstalling {
		actualState := stateCM.Data[appcontext.SetupStateKey]
		if actualState == appcontext.SetupStateInstalling || actualState == appcontext.SetupStateInstalled {
			return fmt.Errorf("setup is busy or already done")
		}
	}

	stateCM.Data[appcontext.SetupStateKey] = state
	_, err = clientSet.CoreV1().ConfigMaps(namespace).Update(ctx, stateCM, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update k8s-ces-setup configmap: %w", err)
	}

	return nil
}
