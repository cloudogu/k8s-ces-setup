package setup

import (
	"context"
	"errors"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/patch"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
)

type setupValidatorStep struct {
	setupJsonValidator         setupJsonConfigurationValidator
	resourcePatchValidator     resourcePatchConfigurationValidator
	setupJsonConfiguration     *appcontext.SetupJsonConfiguration
	resourcePatchConfiguration []patch.ResourcePatch
}

// setupJsonConfigurationValidator is responsible to validate the Cloudogu EcoSystem setup JSON configuration to prevent inconsistent state after a setup.
type setupJsonConfigurationValidator interface {
	Validate(ctx context.Context, setupJson *appcontext.SetupJsonConfiguration) error
}

// resourcePatchConfigurationValidator is responsible to validate the setup resource patch configuration to prevent inconsistent state after a setup.
type resourcePatchConfigurationValidator interface {
	Validate(resourcePatchConfig []patch.ResourcePatch) error
}

// NewValidatorStep creates a new setup step to validate the setup configuration.
func NewValidatorStep(repository cescommons.RemoteDoguDescriptorRepository, setupCtx *appcontext.SetupContext) *setupValidatorStep {
	setupJsonValidator := validation.NewSetupJsonConfigurationValidator(repository)
	resourcePatchValidator := validation.NewResourcePatchConfigurationValidator()

	return &setupValidatorStep{
		setupJsonValidator:         setupJsonValidator,
		resourcePatchValidator:     resourcePatchValidator,
		setupJsonConfiguration:     setupCtx.SetupJsonConfiguration,
		resourcePatchConfiguration: setupCtx.AppConfig.ResourcePatches,
	}
}

// GetStepDescription return the human-readable description of the step.
func (svs *setupValidatorStep) GetStepDescription() string {
	return "Validating the setup configuration"
}

// PerformSetupStep validates the setup configuration.
func (svs *setupValidatorStep) PerformSetupStep(ctx context.Context) error {
	var errs []error

	errs = append(errs, svs.resourcePatchValidator.Validate(svs.resourcePatchConfiguration))
	errs = append(errs, svs.setupJsonValidator.Validate(ctx, svs.setupJsonConfiguration))

	return errors.Join(errs...)
}
