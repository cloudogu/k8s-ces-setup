package setup

import (
	"context"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
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
func SetupAPI(ctx context.Context, router ginRoutes, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *appcontext.SetupContextBuilder) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(ginCtx *gin.Context) {
		startSetup(ctx, ginCtx, clusterConfig, k8sClient, setupContextBuilder)
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

func startSetup(ctx context.Context, ginCtx *gin.Context, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *appcontext.SetupContextBuilder) {
	starter, err := NewStarter(ctx, clusterConfig, k8sClient, setupContextBuilder)
	if err != nil {
		handleInternalServerError(ginCtx, err, "Failed to create setup starter")
		return
	}

	err = starter.StartSetup(ctx)
	if err != nil {
		handleInternalServerError(ginCtx, err, "Failed to start setup")
		return
	}

	ginCtx.Status(http.StatusOK)
}
