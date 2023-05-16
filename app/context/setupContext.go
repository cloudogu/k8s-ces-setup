package context

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-apply-lib/apply"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/sirupsen/logrus"
)

const (
	// SecretDoguRegistry is the name of the secret containing the dogu registry credentials.
	SecretDoguRegistry = "k8s-dogu-operator-dogu-registry"
	// SecretDoguRegistryDevPath is the path to the secret containing the dogu registry credentials. This is used for development.
	SecretDoguRegistryDevPath = "k8s/dev-resources/dogu-registry-secret.yaml"
	// SecretDockerRegistry is the name of the secret containing the docker registry credentials.
	SecretDockerRegistry = "k8s-dogu-operator-docker-registry"
	// SetupConfigConfigmap is the name of the config map containing the setup config.
	SetupConfigConfigmap = "k8s-ces-setup-config"
	// SetupConfigConfigmapDevPath is the path to the config map containing the setup config. This is used for development.
	SetupConfigConfigmapDevPath = "k8s/dev-resources/k8s-ces-setup.yaml"
	// SetupStartUpConfigMap is the name of the config map containing the setup.json.
	SetupStartUpConfigMap = "k8s-ces-setup-json"
	// SetupStartUpConfigMapDevPath is the path to the config map containing the setup.json. This is used for development.
	SetupStartUpConfigMapDevPath = "k8s/dev-resources/setup.json"
	// SetupStateConfigMap is the name of the config map containing the setup state.
	SetupStateConfigMap = "k8s-setup-config"
	// SetupStateKey is the key by which the setup state can be referenced.
	SetupStateKey = "state"
	// SetupStateInstalled means the setup installed the Cloudogu EcoSystem successfully.
	SetupStateInstalled = "installed"
	// SetupStateInstalling means the setup is currently installing the Cloudogu EcoSystem.
	SetupStateInstalling = "installing"
	// EnvironmentVariableStage is the name of the environment variable by which the development stage can be set.
	EnvironmentVariableStage = "STAGE"
	// StageDevelopment is the value that EnvironmentVariableStage must have in order to start the setup in development mode.
	StageDevelopment = "development"
	// EnvironmentVariableTargetNamespace is the name of the environment variable which discerns where the setup should deploy the Cloudogu EcoSystem.
	EnvironmentVariableTargetNamespace = "POD_NAMESPACE"
)

// SetupContext contains all context information provided by the setup.
type SetupContext struct {
	AppVersion                string              `yaml:"app_version"`
	AppConfig                 *Config             `yaml:"app_config"`
	StartupConfiguration      *SetupConfiguration `json:"startup_configuration"`
	DoguRegistryConfiguration *DoguRegistrySecret
}

// SetupContextBuilder contains information to create a setup context
type SetupContextBuilder struct {
	version                   string
	DevSetupConfigPath        string
	DevStartupConfigPath      string
	DevDoguRegistrySecretPath string
}

// NewSetupContextBuilder creates a new builder to create a setup context. Default dev resources paths are used.
func NewSetupContextBuilder(version string) *SetupContextBuilder {
	return &SetupContextBuilder{
		version:                   version,
		DevSetupConfigPath:        SetupConfigConfigmapDevPath,
		DevStartupConfigPath:      SetupStartUpConfigMapDevPath,
		DevDoguRegistrySecretPath: SecretDoguRegistryDevPath,
	}
}

// NewSetupContext creates a new setup context.
func (scb *SetupContextBuilder) NewSetupContext(clientSet kubernetes.Interface) (*SetupContext, error) {
	logrus.Print("Reading configuration file...")

	targetNamespace, err := GetEnvVar(EnvironmentVariableTargetNamespace)
	if err != nil {
		return nil, fmt.Errorf("could not read current namespace: %w", err)
	}

	config, setupJson, doguRegistrySecret, err := scb.getConfigurations(clientSet, targetNamespace)
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
		AppVersion:                scb.version,
		AppConfig:                 config,
		StartupConfiguration:      setupJson,
		DoguRegistryConfiguration: doguRegistrySecret,
	}, nil
}

func (scb *SetupContextBuilder) getConfigurations(clientSet kubernetes.Interface, targetNamespace string) (*Config, *SetupConfiguration, *DoguRegistrySecret, error) {
	if os.Getenv(EnvironmentVariableStage) == StageDevelopment {
		config, err := ReadConfigFromFile(scb.DevSetupConfigPath)
		if err != nil {
			return nil, nil, nil, err
		}

		setupJson, err := ReadSetupConfigFromFile(scb.DevStartupConfigPath)
		if err != nil {
			return nil, nil, nil, err
		}

		doguRegistrySecret, err := ReadDoguRegistrySecretFromFile(scb.DevDoguRegistrySecretPath)
		if err != nil {
			return nil, nil, nil, err
		}

		return config, setupJson, doguRegistrySecret, nil
	} else {
		config, err := ReadConfigFromCluster(clientSet, targetNamespace)
		if err != nil {
			return nil, nil, nil, err
		}

		setupJson, err := ReadSetupConfigFromCluster(clientSet, targetNamespace)
		if err != nil {
			return nil, nil, nil, err
		}

		doguRegistrySecret, err := ReadDoguRegistrySecretFromCluster(clientSet, targetNamespace)

		return config, setupJson, doguRegistrySecret, err
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
				Labels:    map[string]string{"app": "ces", "app.kubernetes.io/name": "k8s-ces-setup"},
			},
		}

		configMap, err = client.CoreV1().ConfigMaps(namespace).Create(context.Background(), setupConfigMap, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create configmap [%s]: %w", SetupStateConfigMap, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get configmap [%s]: %w", SetupConfigConfigmap, err)
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

	core.GetLogger = func() core.Logger {
		return logrus.StandardLogger()
	}

	apply.GetLogger = func() apply.Logger {
		return logrus.StandardLogger()
	}
}
