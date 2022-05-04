package data

import (
	"fmt"
	"github.com/cloudogu/cesapp/v4/registry/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewKeyProviderStep(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		// given
		globalConfigMock := &mocks.ConfigurationContext{}

		// when
		step := NewKeyProviderStep(globalConfigMock)

		// then
		require.NotNil(t, step)
	})
}

func Test_keyProviderSetterStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("get description", func(t *testing.T) {
		// given
		step := &keyProviderSetterStep{}

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, fmt.Sprintf("Set key provider %s", keyProvider), description)
	})
}

func Test_keyProviderSetterStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("successfull set key provider", func(t *testing.T) {
		// given
		globalConfigMock := &mocks.ConfigurationContext{}
		globalConfigMock.On("Set", "key_provider", keyProvider).Return(nil)
		step := &keyProviderSetterStep{globalConfig: globalConfigMock}

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, globalConfigMock)
	})

	t.Run("fail to set key provider", func(t *testing.T) {
		// given
		globalConfigMock := &mocks.ConfigurationContext{}
		globalConfigMock.On("Set", "key_provider", keyProvider).Return(assert.AnError)
		step := &keyProviderSetterStep{globalConfig: globalConfigMock}

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set key provider")
		mock.AssertExpectationsForObjects(t, globalConfigMock)
	})
}
