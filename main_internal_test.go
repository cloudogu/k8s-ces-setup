package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createSetupRouter(t *testing.T) {
	_ = os.Unsetenv("POD_NAMESPACE")

	t.Run("Startup without error", func(t *testing.T) {
		// given
		t.Setenv("POD_NAMESPACE", "myTestNamespace")

		// when
		router, err := createSetupRouter("testdata/k8s-ces-setup-testdata.yaml")

		//then
		require.NoError(t, err)
		assert.NotNil(t, router)
	})

	t.Run("Startup error", func(t *testing.T) {
		// when
		router, err := createSetupRouter("not-a-config")

		//then
		require.Error(t, err)
		assert.Nil(t, router)
		assert.Contains(t, err.Error(), "could not read current namespace")
	})
}
