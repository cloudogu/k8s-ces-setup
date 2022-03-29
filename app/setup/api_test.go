package setup

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getEnvVar(t *testing.T) {
	t.Run("successfully query env var namespace", func(t *testing.T) {
		// given
		t.Setenv("CREDENTIAL_SOURCE_NAMESPACE", "myTestNamespace")

		// when
		ns, err := getEnvVar("CREDENTIAL_SOURCE_NAMESPACE")

		// then
		require.NoError(t, err)

		assert.Equal(t, "myTestNamespace", ns)
	})

	t.Run("failed to query env var namespace", func(t *testing.T) {
		// when
		_, err := getEnvVar("CREDENTIAL_SOURCE_NAMESPACE")

		// then
		require.Error(t, err)
	})
}
