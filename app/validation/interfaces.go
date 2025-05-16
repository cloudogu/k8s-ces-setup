package validation

import (
	ctx "context"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-ces-setup/v4/app/context"
)

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
	ValidateDogus(ctx ctx.Context, dogus context.Dogus) error
}

// RegistryConfigEncryptedValidator is used to validate the registry config encrypted section of the setup configuration
type RegistryConfigEncryptedValidator interface {
	ValidateRegistryConfigEncrypted(config *context.SetupJsonConfiguration) error
}

type remoteDoguDescriptorRepository interface {
	cescommons.RemoteDoguDescriptorRepository
}
