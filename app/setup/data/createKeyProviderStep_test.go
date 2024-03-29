package data_test

import (
	"context"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

var testCtx = context.Background()

func TestNewKeyProviderStep(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		// given
		writerMock := data.NewMockRegistryWriter(t)

		// when
		step := data.NewKeyProviderStep(writerMock, "lalelu")

		// then
		require.NotNil(t, step)
		require.Equal(t, "lalelu", step.KeyProvider)
		require.Equal(t, writerMock, step.Writer)
	})
}

func Test_keyProviderSetterStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("get description", func(t *testing.T) {
		// given
		step := &data.KeyProviderSetterStep{KeyProvider: "key"}

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, "Set key provider key", description)
	})
}

func Test_keyProviderSetterStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("successfull set key provider", func(t *testing.T) {
		// given
		writerMock := data.NewMockRegistryWriter(t)
		keyProviderConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"key_provider": "key",
			},
		}
		writerMock.EXPECT().WriteConfigToRegistry(keyProviderConfig).Return(nil)
		step := data.KeyProviderSetterStep{KeyProvider: "key", Writer: writerMock}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("use default key provider", func(t *testing.T) {
		// given
		writerMock := data.NewMockRegistryWriter(t)
		keyProviderConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"key_provider": data.DefaultKeyProvider,
			},
		}
		writerMock.EXPECT().WriteConfigToRegistry(keyProviderConfig).Return(nil)
		step := data.KeyProviderSetterStep{Writer: writerMock}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("fail to set key provider", func(t *testing.T) {
		// given
		writerMock := data.NewMockRegistryWriter(t)
		writerMock.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)
		step := data.KeyProviderSetterStep{Writer: writerMock}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set key provider")
	})
}
