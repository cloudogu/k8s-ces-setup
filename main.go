package main

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/api"
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// SetupPort defines the port where the setup can be reached. Does not need to be configurable as the setup runs in
// a pod
const SetupPort = 8080

var (
	// Version of the application. Is automatically adapted at compile time.
	Version = "development"
)

func main() {
	logrus.Printf("Starting k8s-ces-setup...")

	logrus.Printf("Reading configuration file...")
	appConfig, err := config.ReadConfig("k8s-ces-setup.yaml")
	if err != nil {
		panic(err)
	}
	config.AppVersion = config.Version(Version)

	err = configureLogger(appConfig)
	if err != nil {
		panic(err)
	}

	logrus.Debugf("Current Version: [%+v]", config.AppVersion)
	logrus.Debugf("Current configuration: [%+v]", appConfig)

	router := setupRouter(appConfig)

	err = router.Run(fmt.Sprintf(":%d", SetupPort))
	if err != nil {
		panic(err)
	}
}

func setupRouter(appConfig config.Config) *gin.Engine {
	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())

	err := api.SetupAPI(router, appConfig)
	if err != nil {
		panic(err)
	}
	return router
}

func configureLogger(appConfig config.Config) error {
	logrus.SetLevel(appConfig.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	return nil
}
