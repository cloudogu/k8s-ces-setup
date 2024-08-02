package data_test

import (
	"context"
	"testing"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewWriteAdminDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new admin data step", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &appcontext.SetupJsonConfiguration{}

		// when
		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
	})
}

func Test_writeAdminDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get admin data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &appcontext.SetupJsonConfiguration{}
		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write admin data to the registry", description)
	})
}

func Test_writeAdminDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("successfully set values for external ldap data", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{
			UserBackend: appcontext.UserBackend{DsType: validation.DsTypeExternal},
			Admin:       appcontext.User{AdminGroup: "myTestAdminGroup"},
		}

		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"admin_group": "myTestAdminGroup",
			},
		}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("successfully set values for embedded ldap data", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{
			Admin: appcontext.User{
				AdminGroup:  "myAdminTestGroup",
				Mail:        "myAdminMail",
				Username:    "myAdminUsername",
				AdminMember: true,
			},
			UserBackend: appcontext.UserBackend{DsType: validation.DsTypeEmbedded},
		}

		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"admin_group": "myAdminTestGroup",
			},
			"ldap": map[string]interface{}{
				"admin_mail":     "myAdminMail",
				"admin_username": "myAdminUsername",
				"admin_member":   "true",
			},
		}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		myStep := data.NewWriteAdminDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
