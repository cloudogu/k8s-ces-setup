package data

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/ssl"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
)

const (
	expireDays = 365
)

type generateSSLStep struct {
	config       *appcontext.SetupJsonConfiguration
	SslGenerator SSLGenerator
}

// SSLGenerator is used to generate a self-signed certificate for a specific fqdn and domain
type SSLGenerator interface {
	GenerateSelfSignedCert(fqdn string, domain string, certExpireDays int, country string,
		province string, locality string, altDNSNames []string) (string, string, error)
}

// NewGenerateSSLStep creates a new setup step which on generates ssl certificates
func NewGenerateSSLStep(config *appcontext.SetupJsonConfiguration) *generateSSLStep {
	generator := ssl.NewSSLGenerator()
	return &generateSSLStep{config: config, SslGenerator: generator}
}

// GetStepDescription return the human-readable description of the step
func (gss *generateSSLStep) GetStepDescription() string {
	return fmt.Sprintf("Generate SSL certificate and key")
}

// PerformSetupStep either generates a certificate if necessary and writes it to the setup configuration
func (gss *generateSSLStep) PerformSetupStep(context.Context) error {
	naming := &gss.config.Naming
	// Generation not needed
	if naming.CertificateType == "external" {
		return nil
	}
	cert, key, err := gss.SslGenerator.GenerateSelfSignedCert(naming.Fqdn, naming.Domain, expireDays, ssl.Country, ssl.Province, ssl.Locality, []string{})
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate and key: %w", err)
	}

	naming.Certificate = cert
	naming.CertificateKey = key

	return nil
}
