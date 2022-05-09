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

func TestNewWriteDoguConfigStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new dogu config step", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteDoguConfigStep(mockRegistry, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func Test_writeDoguConfigStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get dogu config step description", func(t *testing.T) {
		// given
		mockRegistry := &mocks.Registry{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteDoguConfigStep(mockRegistry, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write dogu configuration to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}

func Test_writeDoguConfigStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail on setting default dogu in global config", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Dogus: context.Dogus{DefaultDogu: "myDefaultDogu"}}

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "default_dogu", "myDefaultDogu").Return(assert.AnError)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)

		myStep := data.NewWriteDoguConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, globalRegistryMock, mockRegistry)
	})

	t.Run("fail to set a value in the ldap dogu context", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Dogus: context.Dogus{DefaultDogu: "myDefaultDogu"}}

		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "default_dogu", "myDefaultDogu").Return(nil)

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GlobalConfig").Return(globalRegistryMock)

		myStep := data.NewWriteDoguConfigStep(mockRegistry, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, globalRegistryMock, mockRegistry)
	})
}
