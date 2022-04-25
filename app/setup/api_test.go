package setup

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupAPI(t *testing.T) {
	t.Run("should fail SetupAPI during namespace creation with connection refused", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		SetupAPI(router, context.SetupContext{
			AppVersion: "1.2.3",
			AppConfig: context.Config{
				LogLevel:              logrus.DebugLevel,
				TargetNamespace:       testTargetNamespaceName,
				DoguOperatorURL:       "http://example.com/1.yaml",
				EtcdServerResourceURL: "http://example.com/2.yaml",
				EtcdClientImageRepo:   "bitnami/etcd:3.5.2-debian-10-r0",
			},
		})
		w := httptest.NewRecorder()

		// when
		req, _ := http.NewRequest("POST", "/api/v1/setup", nil)
		router.ServeHTTP(w, req)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		// depending on whether a kube config exists (local dev) or not (CI) errors may be returned at different phases
		assert.Regexp(t, "HTTP 500: An error occurred during this action: (Validate target namespace myfavouritenamespace-1|Load cluster configuration)", w.Body.String())
	})
}
