package health

import (
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const endpointGetHealth = "/api/v1/health"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, context context.SetupContext) {
	requestHandler := newRequestHandler(context.AppVersion)
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodGet, endpointGetHealth)
	router.GET(endpointGetHealth, requestHandler.getHealth)
}
