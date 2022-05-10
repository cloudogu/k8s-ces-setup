package setup

import (
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const endpointPostStartSetup = "/api/v1/setup"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)

	router.POST(endpointPostStartSetup, func(ctx *gin.Context) {
		starter, err := NewStarter(setupContext)
		if err != nil {
			handleInternalServerError(ctx, err, "Failed to start setup")
			return
		}
		err = starter.StartSetup()
		if err != nil {
			handleInternalServerError(ctx, err, "Failed to start setup")
			return
		}
		ctx.Status(http.StatusOK)
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
