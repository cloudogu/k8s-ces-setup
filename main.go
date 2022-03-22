package main

import (
	"fmt"
	"os"

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

// applicationExiter is responsible for exiting the application correctly.
type applicationExiter interface {
	// Exit exits the application and prints the actuator error to the console.
	Exit(err error)
}

type osExiter struct {
}

// Exit prints the actual error to stout and exits the application properly.
func (e *osExiter) Exit(err error) {
	logrus.Errorf("exiting setup because of error: %s", err.Error())
	os.Exit(1)
}

func main() {
	exiter := &osExiter{}

	configFile := "k8s-ces-setup.yaml"
	if os.Getenv("STAGE") == "development" {
		configFile = "k8s/dev-resources/k8s-ces-setup.yaml"
	}
	router := createSetupRouter(exiter, configFile)

	err := router.Run(fmt.Sprintf(":%d", setupPort))
	if err != nil {
		exiter.Exit(err)
	}
}

func createSetupRouter(exiter applicationExiter, configFile string) *gin.Engine {
	logrus.Print("Starting k8s-ces-setup...")

	logrus.Print("Reading configuration file...")
	setupContext, err := context.NewSetupContext(Version, configFile)
	if err != nil {
		exiter.Exit(err)
	}

	configureLogger(setupContext.AppConfig)

	logrus.Debugf("Current Version: [%+v]", setupContext.AppVersion)
	logrus.Debugf("Current context: [%+v]", setupContext)

	return createRouter(setupContext)
}

func createRouter(setupContext context.SetupContext) *gin.Engine {
	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())

	setupAPI(router, setupContext)
	return router
}

func configureLogger(appConfig context.Config) {
	logrus.SetLevel(appConfig.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

// SetupAPI configures the individual endpoints of the API
func setupAPI(router gin.IRoutes, context context.SetupContext) {
	health.SetupAPI(router, context)
	setup.SetupAPI(router, context)
}
