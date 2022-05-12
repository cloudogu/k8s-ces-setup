package data_test

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"

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
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeRegistryConfigDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get registry config data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write registry config data to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeRegistryConfigDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
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
		testConfig := &context.SetupConfiguration{RegistryConfig: registryConfig}

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteRegistryConfigDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
