package data_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appcontext "github.com/cloudogu/k8s-ces-setup/v4/app/context"
	"github.com/cloudogu/k8s-ces-setup/v4/app/setup/data"
)

func TestNewWriteDoguDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new dogu data step", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &appcontext.SetupJsonConfiguration{}

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
		testConfig := &appcontext.SetupJsonConfiguration{}
		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write dogu data to the registry", description)
	})
}

func Test_writeDoguDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("successfully write all dogu data to the registry", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{Dogus: appcontext.Dogus{
			DefaultDogu: "myDefaultDogu",
		}}

		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"default_dogu": "myDefaultDogu",
			},
		}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		myStep := data.NewWriteDoguDataStep(mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
