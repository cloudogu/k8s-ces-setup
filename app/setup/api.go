package setup

import (
	"fmt"
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const endpointPostStartSetup = "/api/v1/setup"

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, appConfig config.Config) error {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)
	router.POST(endpointPostStartSetup, func(context *gin.Context) {
		client, err := createKubernetesClient()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		setupExecutor := NewExecutor(client, appConfig)
		setupExecutor.RegisterSetupStep(newNamespaceCreator(setupExecutor.ClientSet, setupExecutor.Config.Namespace))

		err = setupExecutor.PerformSetup()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		context.Status(http.StatusOK)
	})
	return nil
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
