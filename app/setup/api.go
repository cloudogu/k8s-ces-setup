package setup

import (
	"fmt"
	"net/http"

	"k8s.io/client-go/kubernetes"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

const endpointPostStartSetup = "/api/v1/setup"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {

	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)

	router.POST(endpointPostStartSetup, func(ctx *gin.Context) {
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

		setupExecutor, err := NewExecutor(clusterConfig, clientSet, setupContext)
		if err != nil {
			handleInternalServerError(ctx, err, "Creating setup executor")
			return
		}

		//if setupContext.StartupConfiguration.Naming.CertificateType == "selfsigned" {
		//	err = setupExecutor.RegisterSSLGenerationStep()
		//	if err != nil {
		//		handleInternalServerError(ctx, err, "Register ssl generation setup step")
		//		return
		//	}
		//}
		//
		//err = setupExecutor.RegisterValidationStep()
		//if err != nil {
		//	handleInternalServerError(ctx, err, "Register validation setup steps")
		//	return
		//}
		//
		//err = setupExecutor.RegisterComponentSetupSteps()
		//if err != nil {
		//	handleInternalServerError(ctx, err, "Register component setup steps")
		//	return
		//}

		err = setupExecutor.RegisterDataSetupSteps()
		if err != nil {
			handleInternalServerError(ctx, err, "Register data setup steps")
			return
		}

		err = setupExecutor.RegisterDoguInstallationSteps()
		if err != nil {
			handleInternalServerError(ctx, err, "Register dogu installation steps")
			return
		}

		err, errCausingAction := setupExecutor.PerformSetup()
		if err != nil {
			err2 := fmt.Errorf("error while initializing namespace for setup: %w", err)
			handleInternalServerError(ctx, err2, errCausingAction)
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
