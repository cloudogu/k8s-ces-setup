package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-ces-setup/app/ssl"

	"github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/k8s-ces-setup/app/health"
	"github.com/cloudogu/k8s-ces-setup/app/setup"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// setupPort defines the port where the setup can be reached. Does not need to be configurable as the setup runs in
// a pod.
const setupPort = 8080

var (
	// Version of the application. Is automatically adapted at compile time.
	Version = "development"
)

type osExiter struct {
}

// Exit prints the actual error to stout and exits the application properly.
func (e *osExiter) Exit(err error) {
	logrus.Errorf("exiting setup because of error: %s", err.Error())
	os.Exit(1)
}

func main() {
	exiter := &osExiter{}

	router, err := createSetupRouter(context.NewSetupContextBuilder(Version))
	if err != nil {
		exiter.Exit(err)
	}

	err = router.Run(fmt.Sprintf(":%d", setupPort))
	if err != nil {
		exiter.Exit(err)
	}
}

func createSetupRouter(setupContextBuilder *context.SetupContextBuilder) (*gin.Engine, error) {
	logrus.Print("Starting k8s-ces-setup...")

	clusterConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load cluster configuration: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	setupContext, err := setupContextBuilder.NewSetupContext(clientSet)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Current Version: [%+v]", setupContext.AppVersion)

	if setupContext.StartupConfiguration.IsCompleted() {
		go func() {
			logrus.Info("Setup configuration is completed. Start setup...")
			starter, err := setup.NewStarter(clusterConfig, clientSet, setupContextBuilder)
			if err != nil {
				logrus.Error(err.Error())
			}
			err = starter.StartSetup()
			if err != nil {
				logrus.Error(err.Error())
			}
		}()
	}

	return createRouter(clusterConfig, clientSet, setupContext.AppConfig.TargetNamespace, setupContextBuilder), nil
}

func createRouter(clusterConfig *rest.Config, k8sClient kubernetes.Interface, namespace string, setupContextBuilder *context.SetupContextBuilder) *gin.Engine {
	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())

	setupAPI(router, clusterConfig, k8sClient, namespace, setupContextBuilder)
	return router
}

// SetupAPI configures the individual endpoints of the API
func setupAPI(router gin.IRoutes, clusterConfig *rest.Config, k8sClient kubernetes.Interface, namespace string, setupContextBuilder *context.SetupContextBuilder) {
	health.SetupAPI(router, Version)
	setup.SetupAPI(router, clusterConfig, k8sClient, setupContextBuilder)
	// TODO in Dogu-Operator verschieben
	ssl.SetupAPI(router, namespace)
}
