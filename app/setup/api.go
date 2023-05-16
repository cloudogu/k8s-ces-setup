package setup

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"k8s.io/client-go/rest"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

const endpointPostStartSetup = "/api/v1/setup"

type ginRoutes interface {
	gin.IRoutes
}

// SetupAPI setups the REST API for configuration information
func SetupAPI(router ginRoutes, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *context.SetupContextBuilder) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(ctx *gin.Context) {
		startSetup(ctx, clusterConfig, k8sClient, setupContextBuilder)
	})
}

func handleInternalServerError(ginCtx *gin.Context, err error, causingAction string) {
	logrus.Error(err.Error())
	ginCtx.String(http.StatusInternalServerError, "HTTP %d: An error occurred during this action: %s",
		http.StatusInternalServerError, causingAction)
	ginCtx.Writer.WriteHeaderNow()
	ginCtx.Abort()
	_ = ginCtx.Error(err)
}

func startSetup(ctx *gin.Context, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *context.SetupContextBuilder) {
	starter, err := NewStarter(clusterConfig, k8sClient, setupContextBuilder)
	if err != nil {
		handleInternalServerError(ctx, err, "Failed to create setup starter")
		return
	}

	err = starter.StartSetup()
	if err != nil {
		handleInternalServerError(ctx, err, "Failed to start setup")
		return
	}

	ctx.Status(http.StatusOK)
}
