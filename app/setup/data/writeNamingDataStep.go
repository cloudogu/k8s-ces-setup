package data

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeNamingConfigStep struct {
	Registry      registry.Registry
	Configuration *context.SetupConfiguration
}

// NewWriteNamingConfigStep create a new setup step which writes the naming configuration into the registry.
func NewWriteNamingConfigStep(registry registry.Registry, configuration *context.SetupConfiguration) *writeNamingConfigStep {
	return &writeNamingConfigStep{Registry: registry, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wncs *writeNamingConfigStep) GetStepDescription() string {
	return "Write naming configuration to the registry"
}

func (wncs *writeNamingConfigStep) PerformSetupStep() error {
	globalMap := map[string]string{
		"fqdn":                   wncs.Configuration.Naming.Fqdn,
		"domain":                 wncs.Configuration.Naming.Domain,
		"certificate/type":       wncs.Configuration.Naming.CertificateType,
		"certificate/server.crt": wncs.Configuration.Naming.Certificate,
		"certificate/server.key": wncs.Configuration.Naming.CertificateKey,
		"mail_address":           wncs.Configuration.Naming.MailAddress,
	}

	err := writeMapToContext(wncs.Registry.GlobalConfig(), "_global", globalMap)
	if err != nil {
		return err
	}

	err = wncs.Registry.DoguConfig("postfix").Set("relayhost", wncs.Configuration.Naming.RelayHost)
	if err != nil {
		return fmt.Errorf("could not set relayhost dogu: %w", err)
	}

	return nil
}
