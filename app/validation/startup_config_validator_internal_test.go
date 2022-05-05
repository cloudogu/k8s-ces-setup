package validation

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStartupConfigurationValidator(t *testing.T) {
	t.Run("successfull creating validator", func(t *testing.T) {
		// when
		validator := NewStartupConfigurationValidator(context.SetupConfiguration{}, nil, nil, nil)

		// then
		require.NotNil(t, validator)
	})
}

func Test_validator_ValidateConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("successful validation", func(t *testing.T) {
		// given
		configuration := context.SetupConfiguration{Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: true}, Admin: context.User{Completed: true}}
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(nil)
		adminValidatorMock := &mocks.AdminValidator{}
		adminValidatorMock.On("ValidateAdmin", mock.Anything, mock.Anything).Return(nil)
		validator := NewStartupConfigurationValidator(configuration, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)

		// when
		err := validator.ValidateConfiguration()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})

	t.Run("successful validation with errors because sections aren't completed", func(t *testing.T) {
		// given
		configuration := context.SetupConfiguration{Naming: context.Naming{Completed: false}, UserBackend: context.UserBackend{Completed: false}, Admin: context.User{Completed: false}}
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(assert.AnError)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(assert.AnError)
		adminValidatorMock := &mocks.AdminValidator{}
		adminValidatorMock.On("ValidateAdmin", mock.Anything, mock.Anything).Return(assert.AnError)
		validator := NewStartupConfigurationValidator(configuration, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)

		// when
		err := validator.ValidateConfiguration()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})

	t.Run("error during naming validation", func(t *testing.T) {
		// given
		configuration := context.SetupConfiguration{Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: false}, Admin: context.User{Completed: false}}
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(assert.AnError)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		adminValidatorMock := &mocks.AdminValidator{}
		validator := NewStartupConfigurationValidator(configuration, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)

		// when
		err := validator.ValidateConfiguration()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate naming section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})

	t.Run("error during user backend validation", func(t *testing.T) {
		// given
		configuration := context.SetupConfiguration{Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: true}, Admin: context.User{Completed: false}}
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(assert.AnError)
		adminValidatorMock := &mocks.AdminValidator{}
		validator := NewStartupConfigurationValidator(configuration, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)

		// when
		err := validator.ValidateConfiguration()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate user backend section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})

	t.Run("error during admin user validation", func(t *testing.T) {
		// given
		configuration := context.SetupConfiguration{Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: true}, Admin: context.User{Completed: true}}
		namingValidatorMock := &mocks.NamingValidator{}
		namingValidatorMock.On("ValidateNaming", mock.Anything).Return(nil)
		userBackendValidatorMock := &mocks.UserBackendValidator{}
		userBackendValidatorMock.On("ValidateUserBackend", mock.Anything).Return(nil)
		adminValidatorMock := &mocks.AdminValidator{}
		adminValidatorMock.On("ValidateAdmin", mock.Anything, mock.Anything).Return(assert.AnError)
		validator := NewStartupConfigurationValidator(configuration, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)

		// when
		err := validator.ValidateConfiguration()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate admin user section")
		mock.AssertExpectationsForObjects(t, namingValidatorMock, userBackendValidatorMock, adminValidatorMock)
	})
}
