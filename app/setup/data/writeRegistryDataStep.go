package data

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeConfigToRegistryStep struct {
	Registry      registry.Registry
	Configuration *context.SetupConfiguration
}

// NewWriteConfigToRegistryStep create a new setup step which writes the essential configuration into the registry.
func NewWriteConfigToRegistryStep(registry registry.Registry, configuration *context.SetupConfiguration) *writeConfigToRegistryStep {
	return &writeConfigToRegistryStep{Registry: registry, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wctrs *writeConfigToRegistryStep) GetStepDescription() string {
	return "Write configuration to the registry"
}

// PerformSetupStep writes the essential configuration into the registry.
func (wctrs *writeConfigToRegistryStep) PerformSetupStep() error {
	err := wctrs.writeNamingSection()
	if err != nil {
		return err
	}

	err = wctrs.writeAdminSection()
	if err != nil {
		return err
	}

	err = wctrs.writeDoguSection()
	if err != nil {
		return err
	}

	err = wctrs.writeUserBackendSection()
	if err != nil {
		return err
	}

	err = wctrs.writeRegistryConfigSection()
	if err != nil {
		return err
	}
	return nil
}

// TODO: Currently the internalIp is not handled as it originally creates a host entry on the vm
func (wctrs *writeConfigToRegistryStep) writeNamingSection() error {
	globalMap := map[string]string{
		"fqdn":                   wctrs.Configuration.Naming.Fqdn,
		"domain":                 wctrs.Configuration.Naming.Domain,
		"certificate/type":       wctrs.Configuration.Naming.Domain,
		"certificate/server.crt": wctrs.Configuration.Naming.Certificate,
		"certificate/server.key": wctrs.Configuration.Naming.CertificateKey,
		"mail_address":           wctrs.Configuration.Naming.MailAddress,
	}

	for key, value := range globalMap {
		err := wctrs.Registry.GlobalConfig().Set(key, value)
		if err != nil {
			return fmt.Errorf("failed to set key [%s] in registry: %w", key, err)
		}
	}

	err := wctrs.Registry.DoguConfig("postfix").Set("relayhost", wctrs.Configuration.Naming.RelayHost)
	if err != nil {
		return fmt.Errorf("could not set relayhost dogu: %w", err)
	}

	return nil
}

func (wctrs *writeConfigToRegistryStep) writeAdminSection() error {

	return nil
}

func (wctrs *writeConfigToRegistryStep) writeDoguSection() error {

	return nil
}

func (wctrs *writeConfigToRegistryStep) writeUserBackendSection() error {

	return nil
}

func (wctrs *writeConfigToRegistryStep) writeRegistryConfigSection() error {

	return nil
}
