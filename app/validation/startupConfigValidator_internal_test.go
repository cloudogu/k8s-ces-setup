package validation

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewStartupConfigurationValidator(t *testing.T) {
	t.Run("successfull creating validator", func(t *testing.T) {
		// when
		secret := &v1.Secret{}
		secret.StringData = make(map[string]string)
		secret.StringData["username"] = "user"
		secret.StringData["password"] = "password"
		secret.StringData["endpoint"] = "endpoint"
		validator, err := NewStartupConfigurationValidator(secret)

		// then
		require.NoError(t, err)
		require.NotNil(t, validator)
	})
}

func Test_validator_ValidateConfiguration(t *testing.T) {
	t.Parallel()

	secret := &v1.Secret{}
	secret.StringData = make(map[string]string)
	secret.StringData["username"] = "user"
	secret.StringData["password"] = "password"
	secret.StringData["endpoint"] = "endpoint"

	t.Run("successful validation", func(t *testing.T) {
		// given
		configuration := &context.SetupConfiguration{Dogus: context.Dogus{Completed: true}, Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: true}, Admin: context.User{Completed: true}}
		doguValidatorMock := &mocks.DoguValidator{}
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(nil)
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(nil)
		adminValidatorMock := &mocks.AdminValidator{}
		adminValidatorMock.On("ValidateAdmin", mock.Anything, mock.Anything).Return(nil)
		validator, err := NewStartupConfigurationValidator(secret)
		require.NoError(t, err)
		validator.doguValidator = doguValidatorMock
		validator.namingValidator = namingValidatorMock
		validator.userBackenValidator = userBackendValidatorMock
		validator.adminValidator = adminValidatorMock

		// when
		err = validator.ValidateConfiguration(configuration)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})

	t.Run("error during dogu validation", func(t *testing.T) {
		// given
		configuration := &context.SetupConfiguration{Dogus: context.Dogus{Completed: true}}
		doguValidatorMock := &mocks.DoguValidator{}
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(assert.AnError)
		validator, err := NewStartupConfigurationValidator(secret)
		require.NoError(t, err)
		validator.doguValidator = doguValidatorMock

		// when
		err = validator.ValidateConfiguration(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate dogu section")
		mock.AssertExpectationsForObjects(t, doguValidatorMock)
	})

	t.Run("error during naming validation", func(t *testing.T) {
		// given
		configuration := &context.SetupConfiguration{Naming: context.Naming{Completed: true}}
		namingValidatorMock := &mocks.NamingValidator{}
		doguValidatorMock := &mocks.DoguValidator{}
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(nil)
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(assert.AnError)
		validator, err := NewStartupConfigurationValidator(secret)
		require.NoError(t, err)
		validator.doguValidator = doguValidatorMock
		validator.namingValidator = namingValidatorMock

		// when
		err = validator.ValidateConfiguration(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate naming section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, doguValidatorMock)
	})

	t.Run("error during user backend validation", func(t *testing.T) {
		// given
		configuration := &context.SetupConfiguration{UserBackend: context.UserBackend{Completed: true}}
		doguValidatorMock := &mocks.DoguValidator{}
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(nil)
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(assert.AnError)
		validator, err := NewStartupConfigurationValidator(secret)
		require.NoError(t, err)
		validator.doguValidator = doguValidatorMock
		validator.namingValidator = namingValidatorMock
		validator.userBackenValidator = userBackendValidatorMock

		// when
		err = validator.ValidateConfiguration(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate user backend section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, doguValidatorMock, userBackendValidatorMock)
	})

	t.Run("error during admin user validation", func(t *testing.T) {
		// given
		configuration := &context.SetupConfiguration{Admin: context.User{Completed: true}}
		doguValidatorMock := &mocks.DoguValidator{}
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(nil)
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(nil)
		adminValidatorMock := &mocks.AdminValidator{}
		adminValidatorMock.On("ValidateAdmin", mock.Anything, mock.Anything).Return(assert.AnError)
		validator, err := NewStartupConfigurationValidator(secret)
		require.NoError(t, err)
		validator.doguValidator = doguValidatorMock
		validator.namingValidator = namingValidatorMock
		validator.userBackenValidator = userBackendValidatorMock
		validator.adminValidator = adminValidatorMock

		// when
		err = validator.ValidateConfiguration(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate admin user section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, doguValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})
}
