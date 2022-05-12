package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/k8s-ces-setup/app/validation"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
)

func TestNewWriteAdminDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new admin data step", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeAdminDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get admin data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write admin data to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeAdminDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})

	t.Run("successfully set values for external ldap data", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{
			UserBackend: context.UserBackend{DsType: validation.DsTypeExternal},
			Admin:       context.User{AdminGroup: "myTestAdminGroup"},
		}

		registryConfig := context.CustomKeyValue{
			"_global": map[string]interface{}{
				"admin_group": "myTestAdminGroup",
			},
		}

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})

	t.Run("successfully set values for embedded ldap data", func(t *testing.T) {
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

		registryConfig := context.CustomKeyValue{
			"_global": map[string]interface{}{
				"admin_group": "myAdminTestGroup",
			},
			"ldap": map[string]interface{}{
				"admin_mail":     "myAdminMail",
				"admin_username": "myAdminUsername",
				"admin_member":   "true",
			},
		}

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}
