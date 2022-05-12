package context

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSetupContext(t *testing.T) {
	t.Run("should return new context", func(t *testing.T) {
		// given
		t.Setenv("POD_NAMESPACE", "myTestNamespace")

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

func Test_getEnvVar(t *testing.T) {
	_ = os.Unsetenv("POD_NAMESPACE")

	t.Run("successfully query env var namespace", func(t *testing.T) {
		// given
		t.Setenv("POD_NAMESPACE", "myTestNamespace")

		// when
		ns, err := getEnvVar("POD_NAMESPACE")

		// then
		require.NoError(t, err)

		assert.Equal(t, "myTestNamespace", ns)
	})

	t.Run("failed to query env var namespace", func(t *testing.T) {
		// when
		_, err := getEnvVar("POD_NAMESPACE")

		// then
		require.Error(t, err)
	})
}
