package data

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"strconv"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
)

const tlsSecretName = "ecosystem-certificate"

type writeNamingDataStep struct {
	writer        RegistryWriter
	configuration *appcontext.SetupJsonConfiguration
	clientSet     kubernetes.Interface
	namespace     string
}

// NewWriteNamingDataStep create a new setup step which writes the naming data into the registry.
func NewWriteNamingDataStep(writer RegistryWriter, configuration *appcontext.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *writeNamingDataStep {
	return &writeNamingDataStep{writer: writer, configuration: configuration, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step.
func (wnds *writeNamingDataStep) GetStepDescription() string {
	return "Write naming data to the registry"
}

// PerformSetupStep writes the configured naming data into the registry
func (wnds *writeNamingDataStep) PerformSetupStep(ctx context.Context) error {
	registryConfig := appcontext.CustomKeyValue{
		"_global": map[string]interface{}{
			"fqdn":                wnds.configuration.Naming.Fqdn,
			"domain":              wnds.configuration.Naming.Domain,
			"certificate/type":    wnds.configuration.Naming.CertificateType,
			"mail_address":        wnds.configuration.Naming.MailAddress,
			"k8s/use_internal_ip": strconv.FormatBool(wnds.configuration.Naming.UseInternalIp),
			"k8s/internal_ip":     wnds.configuration.Naming.InternalIp,
		},
		"postfix": map[string]interface{}{
			"relayhost": wnds.configuration.Naming.RelayHost,
		},
	}

	err := wnds.writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write naming data to registry: %w", err)
	}

	return nil
}
