package data

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeSSLStep struct {
	config       context.SetupConfiguration
	globalConfig registry.ConfigurationContext
}

// NeWriteSSL create a new setup step which on writes the certificate to the global config
func NeWriteSSL(config context.SetupConfiguration, globalConfig registry.ConfigurationContext) *writeSSLStep {
	return &writeSSLStep{config: config, globalConfig: globalConfig}
}

// GetStepDescription return the human-readable description of the step
func (gss *writeSSLStep) GetStepDescription() string {
	return fmt.Sprintf("Write SSL certificate and key")
}

// PerformSetupStep writes the certificate from startup config to the global config
func (gss *writeSSLStep) PerformSetupStep() error {
	err := gss.globalConfig.Set("certificate/type", gss.config.Naming.CertificateType)
	if err != nil {
		return fmt.Errorf("failed to set certificate type: %w", err)
	}

	err = gss.globalConfig.Set("certificate/server.crt", gss.config.Naming.Certificate)
	if err != nil {
		return fmt.Errorf("failed to set certificate: %w", err)
	}

	err = gss.globalConfig.Set("certificate/server.key", gss.config.Naming.CertificateKey)
	if err != nil {
		return fmt.Errorf("failed to set certificate key: %w", err)
	}

	return nil
}
