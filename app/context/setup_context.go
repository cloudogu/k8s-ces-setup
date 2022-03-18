package context

import "k8s.io/client-go/rest"

// SetupContext contains all context information provided by the setup.
type SetupContext struct {
	AppVersion    string `yaml:"app_version"`
	AppConfig     Config `yaml:"app_config"`
	K8sRestConfig rest.Config
}

// NewSetupContext creates a new setup context.
func NewSetupContext(version string, configPath string) (SetupContext, error) {
	config, err := ReadConfig(configPath)
	if err != nil {
		return SetupContext{}, err
	}

	return SetupContext{
		AppVersion: version,
		AppConfig:  config,
	}, nil
}
