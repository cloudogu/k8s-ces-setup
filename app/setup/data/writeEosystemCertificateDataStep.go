package data

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-registry-lib/repository"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const certificateSecretName = "ecosystem-certificate"
const certificateSecretPublicKey = "tls.crt"
const certificateSecretPrivateKey = "tls.key"

type writeEcosystemCertificateDataStep struct {
	SecretClient  repository.SecretClient
	Configuration *appcontext.SetupJsonConfiguration
}

func NewWriteEcosystemCertificateDataStep(secretClient repository.SecretClient, configuration *appcontext.SetupJsonConfiguration) *writeEcosystemCertificateDataStep {
	return &writeEcosystemCertificateDataStep{
		SecretClient:  secretClient,
		Configuration: configuration,
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
			certificateSecretPublicKey:  []byte(wecs.Configuration.Naming.Certificate),
			certificateSecretPrivateKey: []byte(wecs.Configuration.Naming.CertificateKey),
		},
	}

	_, err := wecs.SecretClient.Create(ctx, certificateSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to write certificate to secret: %w", err)
	}

	return nil
}
