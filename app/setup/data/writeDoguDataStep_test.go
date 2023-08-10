package data_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
)

func TestNewWriteDoguDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new dogu data step", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &context.SetupJsonConfiguration{}

		// when
		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
	})
}

func Test_writeDoguDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get dogu data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &context.SetupJsonConfiguration{}
		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write dogu data to the registry", description)
	})
}

func Test_writeDoguDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupJsonConfiguration{}
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("successfully write all dogu data to the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupJsonConfiguration{Dogus: context.Dogus{
			DefaultDogu: "myDefaultDogu",
		}}

		registryConfig := context.CustomKeyValue{
			"_global": map[string]interface{}{
				"default_dogu": "myDefaultDogu",
			},
		}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
