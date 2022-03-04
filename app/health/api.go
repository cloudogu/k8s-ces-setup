package health

import (
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const EndpointGetHealth = "/api/v1/health"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, _ config.Config) error {
	requestHandler := newRequestHandler()
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodGet, EndpointGetHealth)
	router.GET(EndpointGetHealth, requestHandler.getHealth)
	return nil
}
