package main

import (
	"fmt"
	"os"

	"github.com/cloudogu/k8s-ces-setup/app/health"
	"github.com/cloudogu/k8s-ces-setup/app/setup"

	"github.com/cloudogu/k8s-ces-setup/app/config"
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

// ApplicationExiter is responsible for exiting the application correctly.
type ApplicationExiter interface {
	// Exit exits the application and prints the actuator error to the console.
	Exit(err error)
}

type osExiter struct {
}

// Exit prints the actuator error to stout and exits the application properly.
func (e *osExiter) Exit(err error) {
	logrus.Errorf("exiting setup because of error: %s", err.Error())
	os.Exit(1)
}

func main() {
	exiter := &osExiter{}

	router := createSetupRouter(exiter, "k8s-ces-setup.yaml")

	err := router.Run(fmt.Sprintf(":%d", setupPort))
	if err != nil {
		exiter.Exit(err)
	}
}

func createSetupRouter(exiter ApplicationExiter, configFile string) *gin.Engine {
	logrus.Printf("Starting k8s-ces-setup...")

	logrus.Printf("Reading configuration file...")
	appConfig, err := config.ReadConfig(configFile)
	if err != nil {
		exiter.Exit(err)
	}

	config.AppVersion = config.Version(Version)
	configureLogger(appConfig)

	logrus.Debugf("Current Version: [%+v]", config.AppVersion)
	logrus.Debugf("Current configuration: [%+v]", appConfig)

	return createRouter(appConfig)
}

func createRouter(appConfig config.Config) *gin.Engine {
	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())

	setupAPI(router, appConfig)
	return router
}

func configureLogger(appConfig config.Config) {
	logrus.SetLevel(appConfig.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

// SetupAPI configures the individual endpoints of the API
func setupAPI(router gin.IRoutes, appConfig config.Config) {
	health.SetupAPI(router, appConfig)
	setup.SetupAPI(router, appConfig)
}
