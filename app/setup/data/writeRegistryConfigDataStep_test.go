package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
)

func TestNewWriteRegistryConfigDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new registry config data step", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &context.SetupJsonConfiguration{}

		// when
		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
	})
}

func Test_writeRegistryConfigDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get registry config data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &context.SetupJsonConfiguration{}
		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write registry config data to the registry", description)
	})
}

func Test_writeRegistryConfigDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupJsonConfiguration{}
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("successfully apply all registry config entries", func(t *testing.T) {
		// given
		registryConfig := context.CustomKeyValue{
			"_global": map[string]interface{}{
				"fqdn":                   "myFqdn",
				"domain":                 "myDomain",
				"certificate/type":       "self-signed",
				"certificate/server.crt": "myCertificate",
				"certificate/server.key": "myCertificateKey",
				"mail_address":           "my@mail.address",
			},
			"postfix": map[string]interface{}{
				"relayhost": "myRelayHost",
			},
		}
		testConfig := &context.SetupJsonConfiguration{RegistryConfig: registryConfig}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
