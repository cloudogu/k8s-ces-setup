package context

import (
	"fmt"
	"io/ioutil"
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
	// EtcdServerResourceURL sets the K8s resource URL which controls the installation of the etcd server into the current cluster.
	EtcdServerResourceURL string `yaml:"etcd_server_url"`
	// EtcdServerResourceURL sets the K8s resource URL which controls the installation of the etcd server into the current cluster.
	EtcdClientImageRepo string `yaml:"etcd_client_image_repo"`
}

// ReadConfig reads the application configuration from a configuration file.
func ReadConfig(path string) (Config, error) {
	config := Config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, fmt.Errorf("could not find configuration at %s", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read configuration %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration %s: %w", path, err)
	}

	return config, nil
}
