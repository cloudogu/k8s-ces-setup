package data

import (
	"fmt"
	"github.com/cloudogu/cesapp/v4/registry"
)

const (
	keyProvider = "pkcs1v15"
)

type keyProviderSetterStep struct {
	globalConfig registry.ConfigurationContext
}

// NewKeyProviderStep create a new setup step which on sets the key provider
func NewKeyProviderStep(globalConfig registry.ConfigurationContext) *keyProviderSetterStep {
	return &keyProviderSetterStep{globalConfig: globalConfig}
}

// GetStepDescription return the human-readable description of the step
func (kps *keyProviderSetterStep) GetStepDescription() string {
	return fmt.Sprintf("Set key provider %s", keyProvider)
}

// PerformSetupStep sets the key provider in the global config
func (kps *keyProviderSetterStep) PerformSetupStep() error {
	err := kps.globalConfig.Set("key_provider", keyProvider)
	if err != nil {
		return fmt.Errorf("failed to set key provider: %w", err)
	}

	return nil
}
