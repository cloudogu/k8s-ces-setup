package context_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/testSetupJson.json
var myTestSetupJson []byte

func TestReadSetupConfig(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		// when
		c, err := context.ReadSetupConfig("testdata/testSetupJson.json")
		var expectedSetupJson context.SetupConfiguration
		err = json.Unmarshal(myTestSetupJson, &expectedSetupJson)
		require.NoError(t, err)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedSetupJson, c)
	})

	t.Run("config does not exists -> empty config returned", func(t *testing.T) {
		// when
		emptyConfig, err := context.ReadSetupConfig("testdata/doesnotexist.json")

		// then
		assert.NoError(t, err)
		assert.Equal(t, context.SetupConfiguration{}, emptyConfig)
	})

	t.Run("fail on invalid file content", func(t *testing.T) {
		// when
		_, err := context.ReadSetupConfig("testdata/invalidConfig.yaml")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal")
	})
}

func TestSetupConfiguration_IsCompleted(t *testing.T) {
	// given
	setupJSON := context.SetupConfiguration{}

	// when+then
	assert.False(t, setupJSON.IsCompleted())

	// when+then
	setupJSON.Naming.Completed = true
	assert.False(t, setupJSON.IsCompleted())

	// when+then
	setupJSON.Dogus.Completed = true
	assert.False(t, setupJSON.IsCompleted())

	// when+then
	setupJSON.Admin.Completed = true
	assert.False(t, setupJSON.IsCompleted())

	// when+then
	setupJSON.UserBackend.Completed = true
	assert.True(t, setupJSON.IsCompleted())
}
