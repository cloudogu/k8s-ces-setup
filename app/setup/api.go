package setup

import (
	"fmt"
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const endpointPostStartSetup = "/api/v1/setup"
const etcdClientVersion = "1.2.3"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext context.SetupContext) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(context *gin.Context) {
		client, err := createKubernetesClient()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		setupExecutor := NewExecutor(client)

		setupExecutor.RegisterSetupStep(newEtcdClientInstallerStep(setupExecutor.ClientSet, etcdClientVersion))
		setupExecutor.RegisterSetupStep(newNamespaceCreator(setupExecutor.ClientSet, setupContext.AppConfig.Namespace))
		setupExecutor.RegisterSetupStep(newEtcdInstallerStep(setupExecutor.ClientSet, setupContext.AppConfig.EtcdServerVersion))
		setupExecutor.RegisterSetupStep(newDoguOperatorInstallerStep(setupExecutor.ClientSet, setupContext.AppConfig.DoguOperatorVersion))

		err = setupExecutor.PerformSetup()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		context.Status(http.StatusOK)
	})
}

func createKubernetesClient() (*kubernetes.Clientset, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot load in cluster configuration: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes configuration: %w", err)
	}
	return clientSet, nil
}
