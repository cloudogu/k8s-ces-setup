package data

import (
	"context"
	"fmt"
	k8sconf "github.com/cloudogu/k8s-registry-lib/config"
	k8sreg "github.com/cloudogu/k8s-registry-lib/repository"
	"k8s.io/client-go/kubernetes"

	appcontext "github.com/cloudogu/k8s-ces-setup/v2/app/context"
	"github.com/cloudogu/k8s-ces-setup/v2/app/validation"
)

type writeRegistryConfigEncryptedStep struct {
	configuration           *appcontext.SetupJsonConfiguration
	sensitiveDoguConfigRepo *k8sreg.DoguConfigRepository
}

// NewWriteRegistryConfigEncryptedStep create a new setup step which writes the registry config encrypted configuration into the cluster.
func NewWriteRegistryConfigEncryptedStep(configuration *appcontext.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *writeRegistryConfigEncryptedStep {
	return &writeRegistryConfigEncryptedStep{
		configuration:           configuration,
		sensitiveDoguConfigRepo: k8sreg.NewSensitiveDoguConfigRepository(clientSet.CoreV1().Secrets(namespace)),
	}
}

// GetStepDescription return the human-readable description of the step.
func (wrces *writeRegistryConfigEncryptedStep) GetStepDescription() string {
	return "Write registry config encrypted data to the registry"
}

// PerformSetupStep writes the registry config data into the registry
func (wrces *writeRegistryConfigEncryptedStep) PerformSetupStep(ctx context.Context) error {
	if wrces.configuration.RegistryConfigEncrypted == nil {
		wrces.configuration.RegistryConfigEncrypted = make(appcontext.CustomKeyValue)
	}
	// Add keys to make sure the sensitive configs are created
	if wrces.configuration.RegistryConfigEncrypted["ldap"] == nil {
		wrces.configuration.RegistryConfigEncrypted["ldap"] = map[string]any{}
	}
	if wrces.configuration.RegistryConfigEncrypted["cas"] == nil {
		wrces.configuration.RegistryConfigEncrypted["cas"] = map[string]any{}
	}
	if wrces.configuration.RegistryConfigEncrypted["ldap-mapper"] == nil {
		wrces.configuration.RegistryConfigEncrypted["ldap-mapper"] = map[string]any{}
	}

	for dogu, resultConfig := range wrces.configuration.RegistryConfigEncrypted {
		entries, err := k8sconf.MapToEntries(resultConfig)
		if err != nil {
			return fmt.Errorf("faild to map config for dogu '%s' to entries: %w", dogu, err)
		}
		doguConfig := k8sconf.CreateDoguConfig(k8sconf.SimpleDoguName(dogu), entries)

		switch dogu {
		case "cas":
			doguConfig, err = wrces.appendCasConfig(doguConfig)
			if err != nil {
				return fmt.Errorf("failed to append default config to cas: %w", err)
			}
		case "ldap":
			doguConfig, err = wrces.appendLdapConfig(doguConfig)
			if err != nil {
				return fmt.Errorf("failed to append default config to ldap: %w", err)
			}
		case "ldap-mapper":
			doguConfig, err = wrces.appendLdapMapperConfig(doguConfig)
			if err != nil {
				return fmt.Errorf("failed to append default config to ldap-mapper: %w", err)
			}
		}

		_, err = wrces.sensitiveDoguConfigRepo.Create(ctx, doguConfig)
		if err != nil {
			return fmt.Errorf("failed to create dogu config for '%s': %w", dogu, err)
		}
	}

	return nil
}

func (wrces *writeRegistryConfigEncryptedStep) appendLdapMapperConfig(doguConfig k8sconf.DoguConfig) (k8sconf.DoguConfig, error) {
	if wrces.configuration.UserBackend.DsType == validation.DsTypeEmbedded {
		return doguConfig, nil
	}

	var err error
	const ldapMapperDoguName = "ldap-mapper"
	if isDoguInstalled(wrces.configuration.Dogus.Install, ldapMapperDoguName) {
		doguConfig.Config, err = doguConfig.Set("backend/password", k8sconf.Value(wrces.configuration.UserBackend.Password))
		if err != nil {
			return doguConfig, err
		}
		doguConfig.Config, err = doguConfig.Set("backend/connection_dn", k8sconf.Value(wrces.configuration.UserBackend.ConnectionDN))

	}

	return doguConfig, err
}

func (wrces *writeRegistryConfigEncryptedStep) appendCasConfig(doguConfig k8sconf.DoguConfig) (k8sconf.DoguConfig, error) {
	var err error

	if wrces.configuration.UserBackend.DsType == validation.DsTypeExternal {
		doguConfig.Config, err = doguConfig.Set("password", k8sconf.Value(wrces.configuration.UserBackend.Password))
	}

	return doguConfig, err
}

func (wrces *writeRegistryConfigEncryptedStep) appendLdapConfig(doguConfig k8sconf.DoguConfig) (k8sconf.DoguConfig, error) {
	if wrces.configuration.UserBackend.DsType != validation.DsTypeEmbedded {
		return doguConfig, nil
	}

	var err error
	doguConfig.Config, err = doguConfig.Set("admin_password", k8sconf.Value(wrces.configuration.Admin.Password))

	return doguConfig, err
}
