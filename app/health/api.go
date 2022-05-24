package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const endpointGetHealth = "/api/v1/health"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, version string) {
	requestHandler := newRequestHandler(version)
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodGet, endpointGetHealth)
	router.GET(endpointGetHealth, requestHandler.getHealth)
}
