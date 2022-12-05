package data

import (
	gocontext "context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

const tlsSecretName = "ecosystem-certificate"

type writeNamingDataStep struct {
	writer        RegistryWriter
	configuration *context.SetupConfiguration
	clientSet     kubernetes.Interface
	namespace     string
}

// NewWriteNamingDataStep create a new setup step which writes the naming data into the registry.
func NewWriteNamingDataStep(writer RegistryWriter, configuration *context.SetupConfiguration, clientSet kubernetes.Interface, namespace string) *writeNamingDataStep {
	return &writeNamingDataStep{writer: writer, configuration: configuration, clientSet: clientSet, namespace: namespace}
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

	secret := &v1.Secret{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name:      tlsSecretName,
			Namespace: wnds.namespace,
		},
		Data: map[string][]byte{
			v1.TLSCertKey:       []byte(wnds.configuration.Naming.Certificate),
			v1.TLSPrivateKeyKey: []byte(wnds.configuration.Naming.CertificateKey),
		},
	}

	_, err = wnds.clientSet.CoreV1().Secrets(wnds.namespace).Create(gocontext.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ecosystem certificate secret: %w", err)
	}

	return nil
}
