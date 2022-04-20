package context

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSetupContext(t *testing.T) {
	t.Run("should return new context", func(t *testing.T) {
		// when
		actual, err := NewSetupContext("1.2.3", "./testdata/testConfig.yaml")

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", actual.AppVersion)
		assert.NotEmpty(t, "1.2.3", actual.AppConfig)
	})
	t.Run("should error for not found config file", func(t *testing.T) {
		// when
		_, err := NewSetupContext("1.2.3", "/nothing/here")

		// then
		require.Error(t, err)
	})
}
