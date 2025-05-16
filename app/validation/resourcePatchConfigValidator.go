package validation

import (
	"errors"
	"github.com/cloudogu/k8s-ces-setup/v4/app/patch"
)

type resourcePatchValidator struct{}

// NewResourcePatchConfigurationValidator creates a new validator.
func NewResourcePatchConfigurationValidator() *resourcePatchValidator {
	return &resourcePatchValidator{}
}

// Validate checks JSON resource patch configurations for configuration errors.
func (r *resourcePatchValidator) Validate(resourcePatchConfig []patch.ResourcePatch) error {
	var errs []error

	for _, resourcePatch := range resourcePatchConfig {
		errs = append(errs, resourcePatch.Validate())
	}

	return errors.Join(errs...)
}
