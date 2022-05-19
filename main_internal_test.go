package main

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createSetupRouter(t *testing.T) {
	_ = os.Unsetenv("POD_NAMESPACE")

	t.Run("Startup without error", func(t *testing.T) {
		// given
		oldGetConfig := ctrl.GetConfig
		defer func() { ctrl.GetConfig = oldGetConfig }()
		ctrl.GetConfig = func() (*rest.Config, error) {
			return &rest.Config{}, nil
		}

		t.Setenv("POD_NAMESPACE", "myTestNamespace")
		t.Setenv("STAGE", "development")
		contextBuilder := &context.SetupContextBuilder{}
		contextBuilder.DevSetupConfigPath = "testdata/k8s-ces-setup-testdata.yaml"
		contextBuilder.DevStartupConfigPath = "testdata/testSetup.json.yaml"

		// when
		router, err := createSetupRouter(contextBuilder)

		// then
		require.NoError(t, err)
		assert.NotNil(t, router)
	})

	t.Run("Startup with config error", func(t *testing.T) {
		// given
		oldGetConfig := ctrl.GetConfig
		defer func() { ctrl.GetConfig = oldGetConfig }()
		ctrl.GetConfig = func() (*rest.Config, error) {
			return nil, assert.AnError
		}

		contextBuilder := &context.SetupContextBuilder{}

		// when
		_, err := createSetupRouter(contextBuilder)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load cluster configuration")
	})
}
