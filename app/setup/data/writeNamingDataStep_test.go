package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
)

func TestWriteNamingConfigStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new naming config step", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteNamingConfigStep(mockRegistry, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func TestWriteNamingConfigStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get naming config step description", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteNamingConfigStep(mockRegistry, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write naming configuration to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func Test_writeConfigToRegistryStep_writeNamingSection(t *testing.T) {
	t.Parallel()

	t.Run("fail on writing anything in the global scope", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)

		myStep := data.NewWriteNamingConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, globalRegistryMock, mockRegistry)
	})

	t.Run("fail on writing anything in the postfix scope", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Naming: context.Naming{
			Fqdn:        "myFqdn",
			Domain:      "myDomain",
			RelayHost:   "myRelayHost",
			MailAddress: "my@mail.address",
		}}

		doguPostfixContextMock := &mocks.ConfigurationContext{}
		doguPostfixContextMock.On("Set", "relayhost", "myRelayHost").Return(assert.AnError)

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", mock.Anything, mock.Anything).Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)
		mockRegistry.On("DoguConfig", "postfix").Return(doguPostfixContextMock)

		myStep := data.NewWriteNamingConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, doguPostfixContextMock, globalRegistryMock, mockRegistry)
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

		doguPostfixContextMock := &mocks.ConfigurationContext{}
		doguPostfixContextMock.On("Set", "relayhost", "myRelayHost").Return(nil)

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "fqdn", "myFqdn").Return(nil)
		globalRegistryMock.On("Set", "domain", "myDomain").Return(nil)
		globalRegistryMock.On("Set", "mail_address", "my@mail.address").Return(nil)
		globalRegistryMock.On("Set", "certificate/type", "self-signed").Return(nil)
		globalRegistryMock.On("Set", "certificate/server.crt", "myCertificate").Return(nil)
		globalRegistryMock.On("Set", "certificate/server.key", "myCertificateKey").Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)
		mockRegistry.On("DoguConfig", "postfix").Return(doguPostfixContextMock)

		myStep := data.NewWriteNamingConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, doguPostfixContextMock, globalRegistryMock, mockRegistry)
	})

	// TODO add tests for internal ip when addressed
}
