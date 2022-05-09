package setup

import (
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
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
func NewValidatorStep(registry remote.Registry, setupCtx *context.SetupContext) *setupValidatorStep {
	validator := validation.NewStartupConfigurationValidator(registry)

	return &setupValidatorStep{Validator: validator, Configuration: &setupCtx.StartupConfiguration}
}

// GetStepDescription return the human-readable description of the step.
func (svs *setupValidatorStep) GetStepDescription() string {
	return "Validating the setup configuration"
}

// PerformSetupStep validates the setup configuration.
func (svs *setupValidatorStep) PerformSetupStep() error {
	return svs.Validator.ValidateConfiguration(svs.Configuration)
}
