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

const (
	secretNameDoguRegistry  = "k8s-dogu-operator-dogu-registry"
	secretNameImageRegistry = "k8s-dogu-operator-docker-registry"
)

type fileClient interface {
	// Get retrieves a file identified by its URL and returns the contents.
	Get(url string) ([]byte, error)
}

type k8sClient interface {
	// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster's namespace.
	Apply(yamlResources []byte, namespace string) error
}

// SetupAPI setups the REST API for configuration information
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {

	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostStartSetup)

	router.POST(endpointPostStartSetup, func(context *gin.Context) {
		clusterConfig, err := ctrl.GetConfig()
		if err != nil {
			handleInternalServerError(context, err, "Load cluster configuration")
			return
		}

		client, err := createKubernetesClient(clusterConfig)
		if err != nil {
			handleInternalServerError(context, err, "Error while creating kubernetes client for the setup API")
			return
		}

		appConfig := setupContext.AppConfig
		// note: with introduction of the setup UI the instance secret may either come into play with a new instance
		// registration or it may already reside in the current namespace
		err = newInstanceSecretValidator(client, appConfig.TargetNamespace).validate()
		if err != nil {
			handleInternalServerError(context, err, "Validate target namespace "+appConfig.TargetNamespace)
			return
		}

		setupExecutor := NewExecutor(client)
		etcdSrvInstallerStep, err := newEtcdServerInstallerStep(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, err, "Create registry step")
			return
		}
		doguOpInstallerStep, err := newDoguOperatorInstallerStep(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, err, "Create dogu operator step")
			return
		}
		serviceDisInstallerStep, err := newServiceDiscoveryInstallerStep(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, err, "Create service discovery step")
			return
		}

		setupExecutor.RegisterSetupStep(etcdSrvInstallerStep)
		setupExecutor.RegisterSetupStep(newEtcdClientInstallerStep(setupExecutor.ClientSet, setupContext))
		setupExecutor.RegisterSetupStep(doguOpInstallerStep)
		setupExecutor.RegisterSetupStep(serviceDisInstallerStep)

		err, errCausingAction := setupExecutor.PerformSetup()
		if err != nil {
			err2 := fmt.Errorf("error while initializing namespace for setup: %w", err)
			handleInternalServerError(context, err2, errCausingAction)
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

func handleInternalServerError(ginCtx *gin.Context, err error, causingAction string) {
	logrus.Error(err.Error())
	ginCtx.String(http.StatusInternalServerError, "HTTP %d: An error occurred during this action: %s",
		http.StatusInternalServerError, causingAction)
	ginCtx.Writer.WriteHeaderNow()
	ginCtx.Abort()
	_ = ginCtx.Error(err)
}
