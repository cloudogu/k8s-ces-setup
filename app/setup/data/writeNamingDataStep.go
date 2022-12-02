package data

import (
	gocontext "context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeNamingDataStep struct {
	writer        RegistryWriter
	configuration *context.SetupConfiguration
	client        rest.Interface
	namespace     string
}

// NewWriteNamingDataStep create a new setup step which writes the naming data into the registry.
func NewWriteNamingDataStep(writer RegistryWriter, configuration *context.SetupConfiguration) *writeNamingDataStep {
	return &writeNamingDataStep{writer: writer, configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wnds *writeNamingDataStep) GetStepDescription() string {
	return "Write naming data to the registry"
}

// PerformSetupStep writes the configured naming data into the registry
func (wnds *writeNamingDataStep) PerformSetupStep() error {
	registryConfig := context.CustomKeyValue{
		"_global": map[string]interface{}{
			"fqdn":                   wnds.configuration.Naming.Fqdn,
			"domain":                 wnds.configuration.Naming.Domain,
			"certificate/type":       wnds.configuration.Naming.CertificateType,
			"certificate/server.crt": wnds.configuration.Naming.Certificate,
			"certificate/server.key": wnds.configuration.Naming.CertificateKey,
			"mail_address":           wnds.configuration.Naming.MailAddress,
		},
		"postfix": map[string]interface{}{
			"relayhost": wnds.configuration.Naming.RelayHost,
		},
	}

	err := wnds.writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write naming data to registry: %w", err)
	}

	secret := &v1.Secret{ObjectMeta: controllerruntime.ObjectMeta{Name: "ecosystem-certificate", Namespace: wnds.namespace}}
	secret.Data[v1.TLSCertKey] = []byte(wnds.configuration.Naming.Certificate)
	secret.Data[v1.TLSPrivateKeyKey] = []byte(wnds.configuration.Naming.CertificateKey)
	err = wnds.client.Post().Namespace(wnds.namespace).Resource("secret").Body(secret).Do(gocontext.Background()).Error()
	if err != nil {
		return fmt.Errorf("failed to create ecosystem certificate secret: %w", err)
	}

	return nil
}
