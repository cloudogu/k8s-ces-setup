package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Version string

var AppVersion Version = "0.0.0"

// Config contains the common configuration for the setup
type Config struct {
	// LogLevel sets the log level for the app
	LogLevel logrus.Level `yaml:"logLevel"`
	// Namespace represents the namespace that is created for the ecosystem
	Namespace string `yaml:"namespace"`
	// DoguOperatorVersion contains the link to the installed dogu operator version
	DoguOperatorVersion string `yaml:"doguOperatorVersion"`
	// EtcdServerVersion contains the link to the installed etcd server version
	EtcdServerVersion string `yaml:"etcdServerVersion"`
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
