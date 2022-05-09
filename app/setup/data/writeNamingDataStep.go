package data

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeNamingDataStep struct {
	Writer        RegistryWriter
	Configuration *context.SetupConfiguration
}

// NewWriteNamingDataStep create a new setup step which writes the naming data into the registry.
func NewWriteNamingDataStep(writer RegistryWriter, configuration *context.SetupConfiguration) *writeNamingDataStep {
	return &writeNamingDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wnds *writeNamingDataStep) GetStepDescription() string {
	return "Write naming data to the registry"
}

func (wnds *writeNamingDataStep) PerformSetupStep() error {
	registryConfig := context.CustomKeyValue{
		"_global": map[string]interface{}{
			"fqdn":                   wnds.Configuration.Naming.Fqdn,
			"domain":                 wnds.Configuration.Naming.Domain,
			"certificate/type":       wnds.Configuration.Naming.CertificateType,
			"certificate/server.crt": wnds.Configuration.Naming.Certificate,
			"certificate/server.key": wnds.Configuration.Naming.CertificateKey,
			"mail_address":           wnds.Configuration.Naming.MailAddress,
		},
		"postfix": map[string]interface{}{
			"relayhost": wnds.Configuration.Naming.RelayHost,
		},
	}

	err := wnds.Writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write naming data to registry: %w", err)
	}

	return nil
}
