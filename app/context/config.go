package context

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/patch"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// Config contains the common configuration for the setup
type Config struct {
	// LogLevel sets the log level for the app
	LogLevel logrus.Level `yaml:"log_level"`
	// TargetNamespace represents the namespace that is created for the ecosystem
	TargetNamespace string `yaml:"target_namespace"`
	// DoguOperatorResourceURL sets the K8s resource URL which controls the installation of the operator into the current cluster.
	DoguOperatorURL string `yaml:"dogu_operator_url"`
	// ServiceDiscoveryURL sets the K8s resource URL which controls the installation of the service discovery into the current cluster.
	ServiceDiscoveryURL string `yaml:"service_discovery_url"`
	// EtcdServerResourceURL sets the K8s resource URL which controls the installation of the etcd server into the current cluster.
	EtcdServerResourceURL string `yaml:"etcd_server_url"`
	// EtcdServerResourceURL sets the K8s resource URL which controls the installation of the etcd server into the current cluster.
	EtcdClientImageRepo string `yaml:"etcd_client_image_repo"`
	// KeyProvider sets the key provider used to encrypt etcd values
	KeyProvider string `yaml:"keyProvider"`
	// ResourcePatches contains json patches for kubernetes resources to be applied on certain phases of the setup process.
	ResourcePatches []patch.ResourcePatch `yaml:"resource_patches"`
}

// ReadConfigFromCluster reads the setup config from the cluster state
func ReadConfigFromCluster(client kubernetes.Interface, namespace string) (*Config, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), SetupConfigConfigmap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get setup configuration from cluster: %w", err)
	}

	config := &Config{}
	stringData := configMap.Data["k8s-ces-setup.yaml"]
	err = yaml.Unmarshal([]byte(stringData), config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarschal configuration from configmap: %w", err)
	}

	return config, nil
}

// ReadConfigFromFile reads the application configuration from a configuration file.
func ReadConfigFromFile(path string) (*Config, error) {
	config := &Config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, fmt.Errorf("could not find configuration at %s", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read configuration %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration %s: %w", path, err)
	}

	return config, nil
}
