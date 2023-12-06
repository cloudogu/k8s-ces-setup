package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
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
	Version           = "development"
	doNotLogEndpoints = []string{health.EndpointGetHealth}
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

	router, err := createSetupRouter(appcontext.NewSetupContextBuilder(Version))
	if err != nil {
		exiter.Exit(err)
	}

	err = router.Run(fmt.Sprintf(":%d", setupPort))
	if err != nil {
		exiter.Exit(err)
	}
}

func createSetupRouter(setupContextBuilder *appcontext.SetupContextBuilder) (*gin.Engine, error) {
	logrus.Print("Starting k8s-ces-setup...")
	logrus.Print("Starting k8s-ces-setup...")
	logrus.Print("Starting k8s-ces-setup...")
	logrus.Print("Starting k8s-ces-setup...")
	logrus.Print("Starting k8s-ces-setup...")
	logrus.Print("Starting k8s-ces-setup...")

	clusterConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load cluster configuration: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	ctx := context.Background()

	setupContext, err := setupContextBuilder.NewSetupContext(ctx, clientSet)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Current Version: [%+v]", setupContext.AppVersion)

	if setupContext.SetupJsonConfiguration.IsCompleted() {
		go func() {
			logrus.Info("Setup configuration is completed. Start setup...")
			starter, err := setup.NewStarter(ctx, clusterConfig, clientSet, setupContextBuilder)
			if err != nil {
				logrus.Error(err.Error())
			}
			err = starter.StartSetup(ctx)
			if err != nil {
				logrus.Error(err.Error())
			}
		}()
	}

	return createRouter(ctx, clusterConfig, clientSet, setupContextBuilder), nil
}

func createRouter(ctx context.Context, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *appcontext.SetupContextBuilder) *gin.Engine {
	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger(), doNotLogEndpoints...), gin.Recovery())

	setupAPI(ctx, router, clusterConfig, k8sClient, setupContextBuilder)
	return router
}

// SetupAPI configures the individual endpoints of the API
func setupAPI(ctx context.Context, router gin.IRoutes, clusterConfig *rest.Config, k8sClient kubernetes.Interface, setupContextBuilder *appcontext.SetupContextBuilder) {
	health.SetupAPI(router, Version)
	setup.SetupAPI(ctx, router, clusterConfig, k8sClient, setupContextBuilder)
}
