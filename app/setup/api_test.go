package setup

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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

func Test_readCredentialSourceNamespace(t *testing.T) {
	t.Run("should try to read credential source namespace from env var if original is empty", func(t *testing.T) {
		t.Setenv("CREDENTIAL_SOURCE_NAMESPACE", "nsfromenvvar")

		// when
		actual, err := readCredentialSourceNamespace("")

		// then
		require.NoError(t, err)
		assert.Equal(t, "nsfromenvvar", actual)
	})
	t.Run("should default to given credential source namespace when env var is unset", func(t *testing.T) {
		// when
		actual, err := readCredentialSourceNamespace("mynamespace")

		// then
		require.NoError(t, err)
		assert.Equal(t, "mynamespace", actual)
	})
	t.Run("should error when no namespace can be found either way", func(t *testing.T) {
		// when
		_, err := readCredentialSourceNamespace("")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read current namespace")
	})
}

func TestSetupAPI(t *testing.T) {
	t.Run("should fail SetupAPI during namespace creation with connection refused", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		SetupAPI(router, context.SetupContext{
			AppVersion: "1.2.3",
			AppConfig: context.Config{
				LogLevel:                  logrus.DebugLevel,
				TargetNamespace:           testTargetNamespaceName,
				CredentialSourceNamespace: "default",
				DoguOperatorURL:           "http://example.com/1.yaml",
				EtcdServerResourceURL:     "http://example.com/2.yaml",
				EtcdClientImageRepo:       "bitnami/etcd:3.5.2-debian-10-r0",
			},
		})
		w := httptest.NewRecorder()

		// when
		req, _ := http.NewRequest("POST", "/api/v1/setup", nil)
		router.ServeHTTP(w, req)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "HTTP 500: An error occurred during this action: Create new namespace myfavouritenamespace-1", w.Body.String())
	})
}
