package context

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/sirupsen/logrus"
)

const SecretDoguRegistry = "k8s-dogu-operator-dogu-registry"
const SecretDockerRegistry = "k8s-dogu-operator-docker-registry"

const SetupConfigConfigmap = "k8s-ces-setup-config"
const SetupConfigConfigmapDevPath = "k8s/dev-resources/k8s-ces-setup.yaml"
const SetupStartUpConfigMap = "k8s-ces-setup-json"
const SetupStartUpConfigMapDevPath = "k8s/dev-resources/setup.json"
const SetupStateConfigMap = "k8s-setup-config"
const SetupStateKey = "state"
const SetupStateInstalled = "installed"
const SetupStateInstalling = "installing"
const EnvironmentVariableStage = "STAGE"
const StageDevelopment = "development"
const EnvironmentVariableTargetNamespace = "POD_NAMESPACE"

// SetupContext contains all context information provided by the setup.
type SetupContext struct {
	AppVersion           string              `yaml:"app_version"`
	AppConfig            *Config             `yaml:"app_config"`
	StartupConfiguration *SetupConfiguration `json:"startup_configuration"`
}

// SetupContextBuilder contains information to create a setup context
type SetupContextBuilder struct {
	version              string
	DevSetupConfigPath   string
	DevStartupConfigPath string
}

// NewSetupContextBuilder creates a new builder to create a setup context. Default dev resources paths are used.
func NewSetupContextBuilder(version string) *SetupContextBuilder {
	return &SetupContextBuilder{
		version:              version,
		DevSetupConfigPath:   SetupConfigConfigmapDevPath,
		DevStartupConfigPath: SetupStartUpConfigMapDevPath,
	}
}

// NewSetupContext creates a new setup context.
func (scb *SetupContextBuilder) NewSetupContext(clientSet kubernetes.Interface) (*SetupContext, error) {
	logrus.Print("Reading configuration file...")

	targetNamespace, err := GetEnvVar(EnvironmentVariableTargetNamespace)
	if err != nil {
		return nil, fmt.Errorf("could not read current namespace: %w", err)
	}

	config, setupJson, err := scb.getConfigurations(clientSet, targetNamespace)
	if err != nil {
		return nil, err
	}

	config.TargetNamespace = targetNamespace

	keyProvider := config.KeyProvider
	if keyProvider != "pkcs1v15" && config.KeyProvider != "oaesp" {
		return nil, fmt.Errorf("invalid key provider: %s", keyProvider)
	}

	configureLogger(config)

	return &SetupContext{
		AppVersion:           scb.version,
		AppConfig:            config,
		StartupConfiguration: setupJson,
	}, nil
}

func (scb *SetupContextBuilder) getConfigurations(clientSet kubernetes.Interface, targetNamespace string) (*Config, *SetupConfiguration, error) {
	if os.Getenv(EnvironmentVariableStage) == StageDevelopment {
		config, err := ReadConfigFromFile(scb.DevSetupConfigPath)
		if err != nil {
			return nil, nil, err
		}

		setupJson, err := ReadSetupConfigFromFile(scb.DevStartupConfigPath)
		return config, setupJson, err
	} else {
		config, err := ReadConfigFromCluster(clientSet, targetNamespace)
		if err != nil {
			return nil, nil, err
		}

		setupJson, err := ReadSetupConfigFromCluster(clientSet, targetNamespace)
		return config, setupJson, err
	}
}

// GetSetupStateConfigMap returns or creates if it does not exist the configmap map for presenting the state of the setup process
func GetSetupStateConfigMap(client kubernetes.Interface, namespace string) (*corev1.ConfigMap, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), SetupStateConfigMap, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		setupConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SetupStateConfigMap,
				Namespace: namespace,
				Labels:    map[string]string{"app.kubernetes.io/name": "k8s-ces-setup"},
			},
		}

		configMap, err = client.CoreV1().ConfigMaps(namespace).Create(context.Background(), setupConfigMap, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create configmap [%s]: %w", SetupStateConfigMap, err)
		}
	}

	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}

	return configMap, nil
}

// GetEnvVar returns an arbitrary environment variable; otherwise it returns an error
func GetEnvVar(name string) (string, error) {
	ns, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	return ns, nil
}

func configureLogger(config *Config) {
	logrus.SetLevel(config.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}
