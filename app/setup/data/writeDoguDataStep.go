package data

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeDoguConfigStep struct {
	Registry      registry.Registry
	Configuration *context.SetupConfiguration
}

// NewWriteDoguConfigStep create a new setup step which writes the dogu configuration into the registry.
func NewWriteDoguConfigStep(registry registry.Registry, configuration *context.SetupConfiguration) *writeDoguConfigStep {
	return &writeDoguConfigStep{Registry: registry, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wdcs *writeDoguConfigStep) GetStepDescription() string {
	return "Write dogu configuration to the registry"
}

func (wdcs *writeDoguConfigStep) PerformSetupStep() error {
	err := wdcs.Registry.GlobalConfig().Set("default_dogu", wdcs.Configuration.Dogus.DefaultDogu)
	if err != nil {
		return fmt.Errorf("could not set default dogu: %w", err)
	}

	return nil
}
