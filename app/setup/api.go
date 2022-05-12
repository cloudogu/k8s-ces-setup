package setup

import (
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const endpointPostStartSetup = "/api/v1/setup"
const endpointPostFinishSetup = "/api/v1/finish"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostFinishSetup)
	router.POST(endpointPostFinishSetup, func(c *gin.Context) {
		finishSetup(c, setupContext)
	})

	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(ctx *gin.Context) {
		startSetup(ctx, setupContext)
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

func startSetup(ctx *gin.Context, setupCtx *context.SetupContext) {
	starter, err := NewStarter(setupCtx)
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

func finishSetup(ctx *gin.Context, setupCtx *context.SetupContext) {
	clusterConfig, err := ctrl.GetConfig()
	if err != nil {
		handleInternalServerError(ctx, err, "Load cluster configuration")
		return
	}

	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		handleInternalServerError(ctx, err, "Cannot create kubernetes client")
		return
	}

	setupFinalizer := NewFinisher(clientSet, setupCtx.AppConfig.TargetNamespace)
	err = setupFinalizer.FinishSetup()
	if err != nil {
		handleInternalServerError(ctx, err, "Cannot finalize setup")
		return
	}

	ctx.Status(http.StatusOK)
}
