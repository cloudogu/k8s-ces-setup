package data_test

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
)

func TestNewWriteNamingDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new naming data step", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}
		fakeClient := fake.NewSimpleClientset()

		// when
		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeNamingDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get naming data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}
		fakeClient := fake.NewSimpleClientset()

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write naming data to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeNamingDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	fakeClient := fake.NewSimpleClientset()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})

	t.Run("successfully apply all naming entries", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Naming: context.Naming{
			Fqdn:            "myFqdn",
			Domain:          "myDomain",
			MailAddress:     "my@mail.address",
			CertificateType: "self-signed",
			Certificate:     "myCertificate",
			CertificateKey:  "myCertificateKey",
			RelayHost:       "myRelayHost",
		}}

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

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})

	// TODO add tests for internal ip when addressed
}
