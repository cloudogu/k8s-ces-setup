package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const EndpointGetHealth = "/api/v1/health"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, version string) {
	requestHandler := newRequestHandler(version)
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodGet, EndpointGetHealth)
	router.GET(EndpointGetHealth, requestHandler.getHealth)
}
