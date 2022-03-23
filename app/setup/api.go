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

type fileClient interface {
	// Get retrieves a file identified by its URL and returns the contents.
	Get(url string) ([]byte, error)
}

type k8sClient interface {
	// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster's namespace.
	Apply(yamlResources []byte, namespace string) error
}

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

		setupExecutor.RegisterSetupStep(newNamespaceCreator(setupExecutor.ClientSet, config.Namespace))
		setupExecutor.RegisterSetupStep(newEtcdServerInstallerStep(clusterConfig, setupContext))
		setupExecutor.RegisterSetupStep(newEtcdClientInstallerStep(clusterConfig, setupContext))
		setupExecutor.RegisterSetupStep(newDoguOperatorInstallerStep(clusterConfig, setupContext))

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
