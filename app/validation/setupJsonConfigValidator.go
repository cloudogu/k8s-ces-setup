package validation

import (
	ctx "context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type validator struct {
	namingValidator                  NamingValidator
	userBackenValidator              UserBackendValidator
	adminValidator                   AdminValidator
	doguValidator                    DoguValidator
	registryConfigEncryptedValidator RegistryConfigEncryptedValidator
}

// NewSetupJsonConfigurationValidator creates a new setup json validator
func NewSetupJsonConfigurationValidator(repository cescommons.RemoteDoguDescriptorRepository) *validator {
	doguValidator := NewDoguValidator(repository)

	return &validator{
		namingValidator:                  NewNamingValidator(),
		userBackenValidator:              NewUserBackendValidator(),
		adminValidator:                   NewAdminValidator(),
		doguValidator:                    doguValidator,
		registryConfigEncryptedValidator: NewRegistryConfigEncryptedValidator(),
	}
}

// Validate checks the section naming, user backend and user from the setup.json configuration
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (v *validator) Validate(ctx ctx.Context, configuration *context.SetupJsonConfiguration) error {
	dogus := configuration.Dogus
	err := v.doguValidator.ValidateDogus(ctx, dogus)
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

	err = v.registryConfigEncryptedValidator.ValidateRegistryConfigEncrypted(configuration)
	if err != nil {
		return fmt.Errorf("failed to validate registry config encrypted section: %w", err)
	}

	return nil
}
