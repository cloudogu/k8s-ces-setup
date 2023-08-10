package validation

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	remoteMocks "github.com/cloudogu/cesapp-lib/remote/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSetupJsonConfigurationValidator(t *testing.T) {
	t.Run("successful creating validator", func(t *testing.T) {
		// when
		secret := &v1.Secret{}
		secret.StringData = make(map[string]string)
		secret.StringData["username"] = "user"
		secret.StringData["password"] = "password"
		secret.StringData["endpoint"] = "endpoint"
		mockRegistry := &remoteMocks.Registry{}
		validator := NewSetupJsonConfigurationValidator(mockRegistry)

		// then
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
		configuration := &context.SetupJsonConfiguration{Dogus: context.Dogus{Completed: true}, Naming: context.Naming{Completed: true}, UserBackend: context.UserBackend{Completed: true}, Admin: context.User{Completed: true}}
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.EXPECT().ValidateDogus(mock.Anything).Return(nil)
		namingValidatorMock := NewMockNamingValidator(t)
		namingValidatorMock.EXPECT().ValidateNaming(mock.Anything).Return(nil)
		userBackendValidatorMock := NewMockUserBackendValidator(t)
		userBackendValidatorMock.EXPECT().ValidateUserBackend(mock.Anything).Return(nil)
		adminValidatorMock := NewMockAdminValidator(t)
		adminValidatorMock.EXPECT().ValidateAdmin(mock.Anything, mock.Anything).Return(nil)
		registryConfigEncryptedValidatorMock := NewMockRegistryConfigEncryptedValidator(t)
		registryConfigEncryptedValidatorMock.EXPECT().ValidateRegistryConfigEncrypted(mock.Anything).Return(nil)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock
		sut.namingValidator = namingValidatorMock
		sut.userBackenValidator = userBackendValidatorMock
		sut.adminValidator = adminValidatorMock
		sut.registryConfigEncryptedValidator = registryConfigEncryptedValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("error during dogu validation", func(t *testing.T) {
		// given
		configuration := &context.SetupJsonConfiguration{Dogus: context.Dogus{Completed: true}}
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.On("ValidateDogus", mock.Anything).Return(assert.AnError)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate dogu section")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("error during naming validation", func(t *testing.T) {
		// given
		configuration := &context.SetupJsonConfiguration{Naming: context.Naming{Completed: true}}
		namingValidatorMock := NewMockNamingValidator(t)
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.EXPECT().ValidateDogus(mock.Anything).Return(nil)
		namingValidatorMock.EXPECT().ValidateNaming(mock.Anything).Return(assert.AnError)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock
		sut.namingValidator = namingValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate naming section")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("error during user backend validation", func(t *testing.T) {
		// given
		configuration := &context.SetupJsonConfiguration{UserBackend: context.UserBackend{Completed: true}}
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.EXPECT().ValidateDogus(mock.Anything).Return(nil)
		namingValidatorMock := NewMockNamingValidator(t)
		namingValidatorMock.EXPECT().ValidateNaming(mock.Anything).Return(nil)
		userBackendValidatorMock := NewMockUserBackendValidator(t)
		userBackendValidatorMock.EXPECT().ValidateUserBackend(mock.Anything).Return(assert.AnError)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock
		sut.namingValidator = namingValidatorMock
		sut.userBackenValidator = userBackendValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate user backend section")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("error during admin user validation", func(t *testing.T) {
		// given
		configuration := &context.SetupJsonConfiguration{Admin: context.User{Completed: true}}
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.EXPECT().ValidateDogus(mock.Anything).Return(nil)
		namingValidatorMock := NewMockNamingValidator(t)
		namingValidatorMock.EXPECT().ValidateNaming(mock.Anything).Return(nil)
		userBackendValidatorMock := NewMockUserBackendValidator(t)
		userBackendValidatorMock.EXPECT().ValidateUserBackend(mock.Anything).Return(nil)
		adminValidatorMock := NewMockAdminValidator(t)
		adminValidatorMock.EXPECT().ValidateAdmin(mock.Anything, mock.Anything).Return(assert.AnError)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock
		sut.namingValidator = namingValidatorMock
		sut.userBackenValidator = userBackendValidatorMock
		sut.adminValidator = adminValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate admin user section")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("error during registry config encrypted validation", func(t *testing.T) {
		// given
		configuration := &context.SetupJsonConfiguration{Admin: context.User{Completed: true}}
		doguValidatorMock := NewMockDoguValidator(t)
		doguValidatorMock.EXPECT().ValidateDogus(mock.Anything).Return(nil)
		namingValidatorMock := NewMockNamingValidator(t)
		namingValidatorMock.EXPECT().ValidateNaming(mock.Anything).Return(nil)
		userBackendValidatorMock := NewMockUserBackendValidator(t)
		userBackendValidatorMock.EXPECT().ValidateUserBackend(mock.Anything).Return(nil)
		adminValidatorMock := NewMockAdminValidator(t)
		adminValidatorMock.EXPECT().ValidateAdmin(mock.Anything, mock.Anything).Return(nil)
		registryConfigEncryptedValidatorMock := NewMockRegistryConfigEncryptedValidator(t)
		registryConfigEncryptedValidatorMock.EXPECT().ValidateRegistryConfigEncrypted(mock.Anything).Return(assert.AnError)
		mockRegistry := &remoteMocks.Registry{}
		sut := NewSetupJsonConfigurationValidator(mockRegistry)
		sut.doguValidator = doguValidatorMock
		sut.namingValidator = namingValidatorMock
		sut.userBackenValidator = userBackendValidatorMock
		sut.adminValidator = adminValidatorMock
		sut.registryConfigEncryptedValidator = registryConfigEncryptedValidatorMock

		// when
		err := sut.Validate(configuration)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate registry config encrypted section")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}
