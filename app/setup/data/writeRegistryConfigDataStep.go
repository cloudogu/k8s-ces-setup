package data

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeRegistryConfigDataStep struct {
	Writer        RegistryWriter
	Configuration *context.SetupConfiguration
}

// NewWriteRegistryConfigDataStep create a new setup step which writes the naming configuration into the registry.
func NewWriteRegistryConfigDataStep(writer RegistryWriter, configuration *context.SetupConfiguration) *writeRegistryConfigDataStep {
	return &writeRegistryConfigDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wrcds *writeRegistryConfigDataStep) GetStepDescription() string {
	return "Write registry config data to the registry"
}

func (wrcds *writeRegistryConfigDataStep) PerformSetupStep() error {
	err := wrcds.Writer.WriteConfigToRegistry(wrcds.Configuration.RegistryConfig)
	if err != nil {
		return fmt.Errorf("failed to write registry config data to registry: %w", err)
	}

	return nil
}
