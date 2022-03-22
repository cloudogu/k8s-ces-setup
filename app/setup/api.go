package setup

import (
	"fmt"
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

const endpointPostStartSetup = "/api/v1/setup"
const etcdClientVersion = "1.2.3"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext context.SetupContext) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(context *gin.Context) {
		clusterConfig, err := ctrl.GetConfig()
		if err != nil {
			logrus.Errorf("cannot load in cluster configuration: %s", err.Error())
			return
		}

		client, err := createKubernetesClient(clusterConfig)
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		setupExecutor := NewExecutor(client)
		config := setupContext.AppConfig

		//setupExecutor.RegisterSetupStep(newNamespaceCreator(setupExecutor.ClientSet, config.Namespace))
		//setupExecutor.RegisterSetupStep(newEtcdInstallerStep(clusterConfig, config.EtcdServerVersion))
		//setupExecutor.RegisterSetupStep(newEtcdClientInstallerStep(sclusterConfig, etcdClientVersion))
		setupExecutor.RegisterSetupStep(newDoguOperatorInstallerStep(clusterConfig, config.DoguOperatorURL, config.DoguOperatorVersion))

		err = setupExecutor.PerformSetup()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		context.Status(http.StatusOK)
	})
}

func createKubernetesClient(clusterConfig *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes configuration: %w", err)
	}

	return clientSet, nil
}
