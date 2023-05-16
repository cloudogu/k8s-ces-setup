package main

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createSetupRouter(t *testing.T) {
	_ = os.Unsetenv(context.EnvironmentVariableTargetNamespace)

	t.Run("Startup without error", func(t *testing.T) {
		// given
		oldGetConfig := ctrl.GetConfig
		defer func() { ctrl.GetConfig = oldGetConfig }()
		ctrl.GetConfig = func() (*rest.Config, error) {
			return &rest.Config{}, nil
		}

		t.Setenv(context.EnvironmentVariableTargetNamespace, "myTestNamespace")
		t.Setenv(context.EnvironmentVariableStage, "development")
		contextBuilder := &context.SetupContextBuilder{}
		contextBuilder.DevSetupConfigPath = "testdata/k8s-ces-setup-testdata.yaml"
		contextBuilder.DevStartupConfigPath = "testdata/testSetup.json.yaml"
		contextBuilder.DevDoguRegistrySecretPath = "testdata/testRegistrySecret.yaml"

		// when
		router, err := createSetupRouter(contextBuilder)

		// then
		require.NoError(t, err)
		assert.NotNil(t, router)
	})

	t.Run("Startup with error while creating client set", func(t *testing.T) {
		// given
		oldGetConfig := ctrl.GetConfig
		defer func() { ctrl.GetConfig = oldGetConfig }()
		ctrl.GetConfig = func() (*rest.Config, error) {
			return &rest.Config{
				AuthProvider: &api.AuthProviderConfig{},
				ExecProvider: &api.ExecConfig{},
			}, nil
		}

		t.Setenv(context.EnvironmentVariableTargetNamespace, "myTestNamespace")
		t.Setenv(context.EnvironmentVariableStage, "development")
		contextBuilder := &context.SetupContextBuilder{}
		contextBuilder.DevSetupConfigPath = "testdata/k8s-ces-setup-testdata.yaml"
		contextBuilder.DevStartupConfigPath = "testdata/testSetup.json.yaml"
		contextBuilder.DevDoguRegistrySecretPath = "testdata/testRegistrySecret.yaml"

		// when
		_, err := createSetupRouter(contextBuilder)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot create kubernetes client")
	})

	t.Run("Startup with error while creating setup context", func(t *testing.T) {
		// given
		oldGetConfig := ctrl.GetConfig
		defer func() { ctrl.GetConfig = oldGetConfig }()
		ctrl.GetConfig = func() (*rest.Config, error) {
			return &rest.Config{}, nil
		}

		_ = os.Unsetenv(context.EnvironmentVariableTargetNamespace)
		t.Setenv(context.EnvironmentVariableStage, "development")
		contextBuilder := &context.SetupContextBuilder{}
		contextBuilder.DevSetupConfigPath = "testdata/k8s-ces-setup-testdata.yaml"
		contextBuilder.DevStartupConfigPath = "testdata/testSetup.json.yaml"
		contextBuilder.DevDoguRegistrySecretPath = "testdata/testRegistrySecret.yaml"

		// when
		_, err := createSetupRouter(contextBuilder)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not read current namespace: POD_NAMESPACE must be set")
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
