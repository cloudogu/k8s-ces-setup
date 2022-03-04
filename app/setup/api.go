package setup

import (
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const EndpointPostStartSetup = "/api/v1/setup"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, appConfig config.Config) error {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, EndpointPostStartSetup)
	router.POST(EndpointPostStartSetup, func(context *gin.Context) {
		setupExecutor, err := NewExecutor(appConfig)
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = setupExecutor.performSetup()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		context.Status(http.StatusOK)
	})
	return nil
}
