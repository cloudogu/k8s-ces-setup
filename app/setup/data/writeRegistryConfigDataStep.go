package data

import (
	"context"
	"fmt"

	appcontext "github.com/cloudogu/k8s-ces-setup/v2/app/context"
)

type writeRegistryConfigDataStep struct {
	Writer        RegistryWriter
	Configuration *appcontext.SetupJsonConfiguration
}

// NewWriteRegistryConfigDataStep create a new setup step which writes the registry config configuration into the registry.
func NewWriteRegistryConfigDataStep(writer RegistryWriter, configuration *appcontext.SetupJsonConfiguration) *writeRegistryConfigDataStep {
	return &writeRegistryConfigDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wrcds *writeRegistryConfigDataStep) GetStepDescription() string {
	return "Write registry config data to the registry"
}

// PerformSetupStep writes the registry config data into the registry
func (wrcds *writeRegistryConfigDataStep) PerformSetupStep(context.Context) error {
	err := wrcds.Writer.WriteConfigToRegistry(wrcds.Configuration.RegistryConfig)
	if err != nil {
		return fmt.Errorf("failed to write registry config data to registry: %w", err)
	}

	return nil
}
