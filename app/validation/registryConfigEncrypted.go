package validation

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"strings"
)

type registryConfigEncryptedValidator struct {
}

// NewRegistryConfigEncryptedValidator creates a new validator for the registryConfigEncrypted region of the setup configuration.
func NewRegistryConfigEncryptedValidator() *registryConfigEncryptedValidator {
	return &registryConfigEncryptedValidator{}
}

// ValidateRegistryConfigEncrypted check whether the registryConfigEncrypted section has invalid dogu keys
func (rcev *registryConfigEncryptedValidator) ValidateRegistryConfigEncrypted(config *context.SetupJsonConfiguration) error {
	for key := range config.RegistryConfigEncrypted {
		keyFound := false
		for _, dogu := range config.Dogus.Install {
			if strings.Contains(dogu, key) {
				keyFound = true
				break
			}
		}

		if !keyFound {
			return fmt.Errorf("key %s does not exist in dogu install list", key)
		}
	}

	return nil
}
