package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
)

func TestNewWriteDoguDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new dogu data step", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeDoguDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get dogu data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write dogu data to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeDoguDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})

	t.Run("successfully write all dogu data to the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{Dogus: context.Dogus{
			DefaultDogu: "myDefaultDogu",
		}}

		registryConfig := context.CustomKeyValue{
			"_global": map[string]interface{}{
				"default_dogu": "myDefaultDogu",
			},
		}

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}
