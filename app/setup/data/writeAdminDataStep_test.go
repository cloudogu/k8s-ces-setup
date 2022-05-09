package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/k8s-ces-setup/app/validation"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
)

func TestWriteAdminConfigStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new admin config step", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func TestWriteAdminConfigStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get admin config step description", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write admin configuration to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func Test_writeConfigToRegistryStep_writeAdminSection(t *testing.T) {
	t.Parallel()

	t.Run("fail on setting admin_group in global config", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Admin: context.User{
			AdminGroup: "myAdminTestGroup",
		}}

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "admin_group", "myAdminTestGroup").Return(assert.AnError)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)

		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, globalRegistryMock, mockRegistry)
	})

	t.Run("successfully set values for external configuration", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{
			Admin: context.User{
				AdminGroup: "myAdminTestGroup",
			},
			UserBackend: context.UserBackend{DsType: validation.DsTypeExternal},
		}

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "admin_group", "myAdminTestGroup").Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)

		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, globalRegistryMock, mockRegistry)
	})

	t.Run("fail to set a value in the ldap dogu context", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{
			Admin: context.User{
				AdminGroup:  "myAdminTestGroup",
				Mail:        "myAdminMail",
				Username:    "myAdminUsername",
				AdminMember: true,
			},
			UserBackend: context.UserBackend{DsType: validation.DsTypeEmbedded},
		}

		doguLdapContextMock := &mocks.ConfigurationContext{}
		doguLdapContextMock.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "admin_group", "myAdminTestGroup").Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)
		mockRegistry.On("DoguConfig", "ldap").Return(doguLdapContextMock)

		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, doguLdapContextMock, globalRegistryMock, mockRegistry)
	})

	t.Run("fail to set a value in the ldap dogu context", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{
			Admin: context.User{
				AdminGroup:  "myAdminTestGroup",
				Mail:        "myAdminMail",
				Username:    "myAdminUsername",
				AdminMember: true,
			},
			UserBackend: context.UserBackend{DsType: validation.DsTypeEmbedded},
		}

		doguLdapContextMock := &mocks.ConfigurationContext{}
		doguLdapContextMock.On("Set", "admin_mail", "myAdminMail").Return(nil)
		doguLdapContextMock.On("Set", "admin_username", "myAdminUsername").Return(nil)
		doguLdapContextMock.On("Set", "admin_member", "true").Return(nil)

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "admin_group", "myAdminTestGroup").Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)
		mockRegistry.On("DoguConfig", "ldap").Return(doguLdapContextMock)

		myStep := data.NewWriteAdminConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, doguLdapContextMock, globalRegistryMock, mockRegistry)
	})
}
