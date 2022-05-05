package validation

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type validator struct {
	configuration       context.SetupConfiguration
	namingValidator     NamingValidator
	userBackenValidator UserBackendValidator
	adminValidator      AdminValidator
}

// NamingValidator is used to validate all the naming section of the setup configuration
type NamingValidator interface {
	ValidateNaming(naming context.Naming) error
}

// UserBackendValidator is used to validate all the user backend section of the setup configuration
type UserBackendValidator interface {
	ValidateUserBackend(backend context.UserBackend) error
}

// AdminValidator is used to validate all the admin section of the setup configuration
type AdminValidator interface {
	ValidateAdmin(admin context.User, dsType string) error
}

// NewStartupConfigurationValidator creates a new setup json validator
func NewStartupConfigurationValidator(configuration context.SetupConfiguration, namingValidator NamingValidator, userBackenValidator UserBackendValidator, adminValidator AdminValidator) *validator {
	return &validator{
		configuration:       configuration,
		namingValidator:     namingValidator,
		userBackenValidator: userBackenValidator,
		adminValidator:      adminValidator,
	}
}

// ValidateConfiguration checks the section naming, user backend and user from the setup.json configuration
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (v *validator) ValidateConfiguration() error {
	naming := v.configuration.Naming
	err := v.namingValidator.ValidateNaming(naming)
	if err != nil && naming.Completed {
		return fmt.Errorf("failed to validate naming section: %w", err)
	}

	userBackend := v.configuration.UserBackend
	err = v.userBackenValidator.ValidateUserBackend(userBackend)
	if err != nil && userBackend.Completed {
		return fmt.Errorf("failed to validate user backend section: %w", err)
	}

	admin := v.configuration.Admin
	err = v.adminValidator.ValidateAdmin(admin, userBackend.DsType)
	if err != nil && admin.Completed {
		return fmt.Errorf("failed to validate admin user section: %w", err)
	}

	return nil
}
