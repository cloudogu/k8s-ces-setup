package validation

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type validator struct {
	namingValidator     NamingValidator
	userBackenValidator UserBackendValidator
	adminValidator      AdminValidator
	doguValidator       DoguValidator
}

// NamingValidator is used to validate the naming section of the setup configuration
type NamingValidator interface {
	ValidateNaming(naming context.Naming) error
}

// UserBackendValidator is used to validate the user backend section of the setup configuration
type UserBackendValidator interface {
	ValidateUserBackend(backend context.UserBackend) error
}

// AdminValidator is used to validate the admin section of the setup configuration
type AdminValidator interface {
	ValidateAdmin(admin context.User, dsType string) error
}

// DoguValidator is used to validate the dogu section of the setup configuration
type DoguValidator interface {
	ValidateDogus(dogus context.Dogus) error
}

// NewStartupConfigurationValidator creates a new setup json validator
func NewStartupConfigurationValidator(secret *v1.Secret) (*validator, error) {
	doguValidator, err := NewDoguValidator(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create dogu validator: %w", err)
	}

	return &validator{
		namingValidator:     NewNamingValidator(),
		userBackenValidator: NewUserBackendValidator(),
		adminValidator:      NewAdminValidator(),
		doguValidator:       doguValidator,
	}, nil
}

// ValidateConfiguration checks the section naming, user backend and user from the setup.json configuration
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (v *validator) ValidateConfiguration(configuration *context.SetupConfiguration) error {
	dogus := configuration.Dogus
	err := v.doguValidator.ValidateDogus(dogus)
	if err != nil {
		return fmt.Errorf("failed to validate dogu section: %w", err)
	}

	naming := configuration.Naming
	err = v.namingValidator.ValidateNaming(naming)
	if err != nil {
		return fmt.Errorf("failed to validate naming section: %w", err)
	}

	userBackend := configuration.UserBackend
	err = v.userBackenValidator.ValidateUserBackend(userBackend)
	if err != nil {
		return fmt.Errorf("failed to validate user backend section: %w", err)
	}

	admin := configuration.Admin
	err = v.adminValidator.ValidateAdmin(admin, userBackend.DsType)
	if err != nil {
		return fmt.Errorf("failed to validate admin user section: %w", err)
	}

	return nil
}
