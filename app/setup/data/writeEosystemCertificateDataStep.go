package data

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/v2/app/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const certificateSecretName = "ecosystem-certificate"

type writeEcosystemCertificateDataStep struct {
	secretClient  secretClient
	configuration *appcontext.SetupJsonConfiguration
}

func NewWriteEcosystemCertificateDataStep(secretClient secretClient, configuration *appcontext.SetupJsonConfiguration) *writeEcosystemCertificateDataStep {
	return &writeEcosystemCertificateDataStep{
		secretClient:  secretClient,
		configuration: configuration,
	}
}

// GetStepDescription return the human-readable description of the step.
func (wecs *writeEcosystemCertificateDataStep) GetStepDescription() string {
	return "Write ecosystem certificate data to a secret"
}

// PerformSetupStep writes the configured certificate data into a secret.
func (wecs *writeEcosystemCertificateDataStep) PerformSetupStep(ctx context.Context) error {
	certificateSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: certificateSecretName,
		},
		Data: map[string][]byte{
			v1.TLSCertKey:       []byte(wecs.configuration.Naming.Certificate),
			v1.TLSPrivateKeyKey: []byte(wecs.configuration.Naming.CertificateKey),
		},
	}

	_, err := wecs.secretClient.Create(ctx, certificateSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to write certificate to secret: %w", err)
	}

	return nil
}
