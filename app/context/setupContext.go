package context

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const SecretDoguRegistry = "k8s-dogu-operator-dogu-registry"
const SecretDockerRegistry = "k8s-dogu-operator-docker-registry"

const defaultSetupConfigJson = "/setup.json"

// SetupContext contains all context information provided by the setup.
type SetupContext struct {
	AppVersion           string             `yaml:"app_version"`
	AppConfig            Config             `yaml:"app_config"`
	StartupConfiguration SetupConfiguration `json:"startup_configuration"`
}

// NewSetupContext creates a new setup context.
func NewSetupContext(version string, configPath string) (*SetupContext, error) {
	logrus.Print("Reading configuration file...")

	targetNamespace, err := getEnvVar("POD_NAMESPACE")
	if err != nil {
		err2 := fmt.Errorf("could not read current namespace: %w", err)
		return nil, err2
	}

	config, err := ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	config.TargetNamespace = targetNamespace

	setupJsonPath := defaultSetupConfigJson
	if os.Getenv("STAGE") == "development" {
		setupJsonPath = "k8s/dev-resources/setup.json"
	}
	setupJson, err := ReadSetupConfig(setupJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get setup configuration: %w", err)
	}

	return &SetupContext{
		AppVersion:           version,
		AppConfig:            config,
		StartupConfiguration: setupJson,
	}, nil
}

// getEnvVar returns an arbitrary environment variable; otherwise it returns an error
func getEnvVar(name string) (string, error) {
	ns, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	return ns, nil
}
