package context

// SetupContext contains all context information provided by the setup.
type SetupContext struct {
	AppVersion string `yaml:"app_version"`
	AppConfig  Config `yaml:"app_config"`
}

// NewSetupContext creates a new setup context.
func NewSetupContext(version string, configPath string, targetNamespace string) (SetupContext, error) {
	config, err := ReadConfig(configPath)
	if err != nil {
		return SetupContext{}, err
	}

	config.TargetNamespace = targetNamespace

	return SetupContext{
		AppVersion: version,
		AppConfig:  config,
	}, nil
}
