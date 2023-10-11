package context

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"github.com/sirupsen/logrus"

	"github.com/cloudogu/k8s-ces-setup/app/patch"
)

// ComponentAttributes defines the component to install.
type ComponentAttributes struct {
	// Version specifies the component version.
	Version string `json:"version" yaml:"version"`
	// HelmRepositoryNamespace specifies the component namespace inside the originating helm repository, f. i. "k8s" or "official".
	HelmRepositoryNamespace string `json:"helmRepositoryNamespace" yaml:"helmRepositoryNamespace"`
	// DeployNamespace specifies the cluster namespace where the helm chart should be deployed to. If left empty, the currently selected cluster namespace will be used by the component operator.
	// +optional
	DeployNamespace string `json:"deployNamespace" yaml:"deployNamespace"`
}

// Config contains the common configuration for the setup
type Config struct {
	// LogLevel sets the log level for the app
	LogLevel *logrus.Level `json:"log_level" yaml:"log_level"`
	// TargetNamespace represents the namespace that is created for the ecosystem
	TargetNamespace string `json:"target_namespace" yaml:"target_namespace"`
	// ComponentOperatorCrdChart sets the Helm-Chart which controls the installation of the component-operator CRD into the current cluster.
	ComponentOperatorCrdChart string `json:"component_operator_crd_chart" yaml:"component_operator_crd_chart"`
	// ComponentOperatorChart sets the Helm-Chart which controls the installation of the component-operator into the current cluster.
	ComponentOperatorChart string `json:"component_operator_chart" yaml:"component_operator_chart"`
	// Components sets the List of Components that should be installed by the setup
	Components map[string]ComponentAttributes `json:"components" yaml:"components"`
	// EtcdServerResourceURL sets the K8s resource URL which controls the installation of the etcd server into the current cluster.
	EtcdClientImageRepo string `json:"etcd_client_image_repo" yaml:"etcd_client_image_repo"`
	// KeyProvider sets the key provider used to encrypt etcd values
	KeyProvider string `json:"key_provider" yaml:"key_provider"`
	// ResourcePatches contains json patches for kubernetes resources to be applied on certain phases of the setup process.
	ResourcePatches []patch.ResourcePatch `json:"resource_patches" yaml:"resource_patches"`
}

// ReadConfigFromCluster reads the setup config from the cluster state
func ReadConfigFromCluster(ctx context.Context, client kubernetes.Interface, namespace string) (*Config, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, SetupConfigConfigmap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get setup configuration from cluster: %w", err)
	}

	config := &Config{}
	stringData := configMap.Data["k8s-ces-setup.yaml"]
	err = yaml.Unmarshal([]byte(stringData), config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration from configmap: %w", err)
	}

	return config, nil
}

// ReadConfigFromFile reads the application configuration from a configuration file.
func ReadConfigFromFile(path string) (*Config, error) {
	config := &Config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, fmt.Errorf("could not find configuration at %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read configuration %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration %s: %w", path, err)
	}

	return config, nil
}
