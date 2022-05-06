package data

import (
	"github.com/cloudogu/cesapp-lib/registry"
)

type writeConfigToRegistryStep struct {
	globalConfig registry.ConfigurationContext
}

// NewWriteConfigToRegistryStep create a new setup step which writes the essential configuration into the registry.
func NewWriteConfigToRegistryStep(globalConfig registry.ConfigurationContext) *keyProviderSetterStep {
	return &keyProviderSetterStep{globalConfig: globalConfig}
}

// GetStepDescription return the human-readable description of the step.
func (kps *writeConfigToRegistryStep) GetStepDescription() string {
	return "Write configuration to the registry"
}

// PerformSetupStep writes the essential configuration into the registry.
func (kps *writeConfigToRegistryStep) PerformSetupStep() error {

	return nil
}
