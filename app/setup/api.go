package setup

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

const endpointPostStartSetup = "/api/v1/setup"

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
			handleInternalServerError(context, errors.Wrap(err, "cannot load in cluster configuration"))
			return
		}

		client, err := createKubernetesClient(clusterConfig)
		if err != nil {
			handleInternalServerError(context, errors.Wrap(err, "error while creating kubernetes client for the setup API"))
			return
		}

		appConfig := setupContext.AppConfig
		credentialSourceNamespace, err := readCredentialSourceNamespace(appConfig.CredentialSourceNamespace)
		if err != nil {
			handleInternalServerError(context, errors.Wrap(err, "error while setting up the setup API"))
			return
		}

		etcdSrvInstallerStep, err := newEtcdServerInstallerStep(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, errors.Wrap(err, "error while setting up the setup API"))
			return
		}
		doguOpInstallerStep, err := newDoguOperatorInstallerStep(clusterConfig, setupContext)
		if err != nil {
			handleInternalServerError(context, errors.Wrap(err, "error while setting up the setup API"))
			return
		}

		setupExecutor := NewExecutor(client)
		setupExecutor.RegisterSetupStep(newNamespaceCreator(setupExecutor.ClientSet, appConfig.TargetNamespace))
		// maybe we should even transport the pure credential pair instead of the meta-namespace?
		setupExecutor.RegisterSetupStep(newSecretCreator(setupExecutor.ClientSet, appConfig.TargetNamespace, credentialSourceNamespace))
		setupExecutor.RegisterSetupStep(etcdSrvInstallerStep)
		setupExecutor.RegisterSetupStep(newEtcdClientInstallerStep(setupExecutor.ClientSet, setupContext))
		setupExecutor.RegisterSetupStep(doguOpInstallerStep)

		err = setupExecutor.PerformSetup()
		if err != nil {
			handleInternalServerError(context, errors.Wrap(err, "error while setting up the setup API"))
			return
		}

		context.Status(http.StatusOK)
	})
}

func readCredentialSourceNamespace(credSourceNamespaceFromConfig string) (string, error) {
	if credSourceNamespaceFromConfig != "" {
		return credSourceNamespaceFromConfig, nil
	}

	credentialSourceNamespace, err := getEnvVar("CREDENTIAL_SOURCE_NAMESPACE")
	if err != nil {
		return "", errors.Wrap(err, "failed to read current namespace from CREDENTIAL_SOURCE_NAMESPACE")
	}

	return credentialSourceNamespace, err
}

func createKubernetesClient(clusterConfig *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes configuration: %w", err)
	}

	return clientSet, nil
}

// getEnvVar returns the namespace the operator should be watching for changes
func getEnvVar(name string) (string, error) {
	ns, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	return ns, nil
}

func handleInternalServerError(ginCtx *gin.Context, err error) {
	logrus.Error(err.Error())
	_ = ginCtx.AbortWithError(http.StatusInternalServerError, err)
}
