package setup

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gocontext "context"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
	"k8s.io/client-go/kubernetes"
)

type setupValidatorStep struct {
	// Validator contains a setup configuration validator.
	Validator ConfigurationValidator `json:"validator"`
	// Validator contains a setup configuration validator.
	Configuration *context.SetupConfiguration `json:"configuration"`
}

// ConfigurationValidator is responsible to validate the setup configuration to prevent inconsistent state after a setup.
type ConfigurationValidator interface {
	ValidateConfiguration(configuration *context.SetupConfiguration) error
}

// NewValidatorStep creates a new setup step to validate the setup configuration.
func NewValidatorStep(client kubernetes.Interface, setupCtx *context.SetupContext) (*setupValidatorStep, error) {
	registrySecret, err := client.CoreV1().Secrets(setupCtx.AppConfig.TargetNamespace).Get(gocontext.Background(), context.SecretDoguRegistry, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret [%s]: %w", context.SecretDoguRegistry, err)
	}

	validator, err := validation.NewStartupConfigurationValidator(registrySecret)
	if err != nil {
		return nil, err
	}

	return &setupValidatorStep{Validator: validator, Configuration: &setupCtx.StartupConfiguration}, nil
}

// GetStepDescription return the human-readable description of the step.
func (svs *setupValidatorStep) GetStepDescription() string {
	return "Validating the setup configuration"
}

// PerformSetupStep validates the setup configuration.
func (svs *setupValidatorStep) PerformSetupStep() error {
	return svs.Validator.ValidateConfiguration(svs.Configuration)
}
