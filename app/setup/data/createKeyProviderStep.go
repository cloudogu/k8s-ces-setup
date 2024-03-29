package data

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
)

const (
	// DefaultKeyProvider is the standard with which keys and certificates get created in the Cloudogu EcoSystem.
	DefaultKeyProvider = "pkcs1v15"
)

// KeyProviderSetterStep writes the key provider in the registry
type KeyProviderSetterStep struct {
	Writer      RegistryWriter
	KeyProvider string
}

// NewKeyProviderStep create a new setup step which on sets the key provider
func NewKeyProviderStep(writer RegistryWriter, keyProvider string) *KeyProviderSetterStep {
	return &KeyProviderSetterStep{
		Writer:      writer,
		KeyProvider: keyProvider,
	}
}

// GetStepDescription return the human-readable description of the step
func (kps *KeyProviderSetterStep) GetStepDescription() string {
	return fmt.Sprintf("Set key provider %s", kps.KeyProvider)
}

// PerformSetupStep sets the key provider in the global config
func (kps *KeyProviderSetterStep) PerformSetupStep(context.Context) error {
	if kps.KeyProvider == "" {
		kps.KeyProvider = DefaultKeyProvider
	}

	keyProviderConfig := appcontext.CustomKeyValue{
		"_global": map[string]interface{}{
			"key_provider": kps.KeyProvider,
		},
	}

	err := kps.Writer.WriteConfigToRegistry(keyProviderConfig)
	if err != nil {
		return fmt.Errorf("failed to set key provider: %w", err)
	}

	return nil
}
