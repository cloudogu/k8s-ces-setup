package setup

import (
	"fmt"
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

const endpointPostStartSetup = "/api/v1/setup"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {

	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)

	router.POST(endpointPostStartSetup, func(context *gin.Context) {
		clusterConfig, err := ctrl.GetConfig()
		if err != nil {
			handleInternalServerError(context, err, "Load cluster configuration")
			return
		}

		setupExecutor, err := NewExecutor(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, err, "Creating setup executor")
			return
		}

		err = setupExecutor.RegisterComponentSetupSteps()
		if err != nil {
			handleInternalServerError(context, err, "Register component setup steps")
			return
		}

		err = setupExecutor.RegisterDataSetupSteps()
		if err != nil {
			handleInternalServerError(context, err, "Register data setup steps")
			return
		}

		err, errCausingAction := setupExecutor.PerformSetup()
		if err != nil {
			err2 := fmt.Errorf("error while initializing namespace for setup: %w", err)
			handleInternalServerError(context, err2, errCausingAction)
			return
		}

		context.Status(http.StatusOK)
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
