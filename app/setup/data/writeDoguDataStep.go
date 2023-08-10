package data

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeDoguDataStep struct {
	Writer        RegistryWriter
	Configuration *context.SetupJsonConfiguration
}

// NewWriteDoguDataStep create a new setup step which writes the dogu data into the registry.
func NewWriteDoguDataStep(writer RegistryWriter, configuration *context.SetupJsonConfiguration) *writeDoguDataStep {
	return &writeDoguDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wdds *writeDoguDataStep) GetStepDescription() string {
	return "Write dogu data to the registry"
}

// PerformSetupStep writes the configured dogu data into the registry
func (wdds *writeDoguDataStep) PerformSetupStep() error {
	registryConfig := context.CustomKeyValue{
		"_global": map[string]interface{}{
			"default_dogu": wdds.Configuration.Dogus.DefaultDogu,
		},
	}

	err := wdds.Writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write dogu data to registry: %w", err)
	}

	return nil
}
