package data

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/ssl"
)

type WriteSSLStep struct {
	Config    *context.SetupConfiguration
	SSLWriter SSLWriter
}

type SSLWriter interface {
	WriteCertificate(certType string, cert string, key string) error
}

// NewWriteSSLStep creates a new setup step which on writes the certificate to the global Config
func NewWriteSSLStep(config *context.SetupConfiguration, globalConfig registry.ConfigurationContext) *WriteSSLStep {
	sslWriter := ssl.NewSSLWriter(globalConfig)
	return &WriteSSLStep{Config: config, SSLWriter: sslWriter}
}

// GetStepDescription return the human-readable description of the step
func (gss *WriteSSLStep) GetStepDescription() string {
	return fmt.Sprintf("Write SSL certificate and key")
}

// PerformSetupStep writes the certificate from startup config to the global config
func (gss *WriteSSLStep) PerformSetupStep() error {
	naming := gss.Config.Naming
	err := gss.SSLWriter.WriteCertificate(naming.CertificateType, naming.Certificate, naming.CertificateKey)
	if err != nil {
		return err
	}

	return nil
}
