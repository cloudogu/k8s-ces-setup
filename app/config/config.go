package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Version string

var AppVersion Version = "0.0.0"

// Config contains the common configuration for the application.
type Config struct {
	LogLevel            logrus.Level `yaml:"log_level"`
	Namespace           string       `yaml:"namespace"`
	DoguOperatorVersion string       `yaml:"dogu_operator_version"`
	EtcdServerVersion   string       `yaml:"etcd_server_version"`
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
