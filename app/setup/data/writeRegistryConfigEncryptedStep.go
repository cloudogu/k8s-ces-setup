package data

import (
	gocontext "context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
)

type writeRegistryConfigEncryptedStep struct {
	configuration *context.SetupJsonConfiguration
	clientSet     kubernetes.Interface
	namespace     string
	Writer        MapWriter
}

// MapWriter is responsible to write entries into a map[string]map[string]string{}.
type MapWriter interface {
	WriteConfigToStringDataMap(registryConfig context.CustomKeyValue) (map[string]map[string]string, error)
}

// NewWriteRegistryConfigEncryptedStep create a new setup step which writes the registry config encrypted configuration into the cluster.
func NewWriteRegistryConfigEncryptedStep(configuration *context.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *writeRegistryConfigEncryptedStep {
	return &writeRegistryConfigEncryptedStep{configuration: configuration, clientSet: clientSet, namespace: namespace, Writer: &stringDataConfigurationWriter{}}
}

// GetStepDescription return the human-readable description of the step.
func (wrces *writeRegistryConfigEncryptedStep) GetStepDescription() string {
	return "Write registry config encrypted data to the registry"
}

// PerformSetupStep writes the registry config data into the registry
func (wrces *writeRegistryConfigEncryptedStep) PerformSetupStep() error {
	resultConfigs, err := wrces.Writer.WriteConfigToStringDataMap(wrces.configuration.RegistryConfigEncrypted)
	if err != nil {
		return fmt.Errorf("failed to write registry config encrypted: %w", err)
	}

	// append edge cases
	wrces.appendLdapConfig(resultConfigs)
	wrces.appendCasConfig(resultConfigs)
	wrces.appendLdapMapperConfig(resultConfigs)

	// write secrets
	for dogu, resultConfig := range resultConfigs {
		err := wrces.createRegistryConfigEncryptedSecret(dogu, resultConfig)
		if err != nil {
			return fmt.Errorf("failed create %s-secrets: %w", dogu, err)
		}
	}

	return nil
}

func (wrces *writeRegistryConfigEncryptedStep) appendLdapMapperConfig(resultConfigs map[string]map[string]string) {
	if wrces.configuration.UserBackend.DsType == validation.DsTypeEmbedded {
		return
	}

	const ldapMapperDoguName = "ldap-mapper"
	if isDoguInstalled(wrces.configuration.Dogus.Install, ldapMapperDoguName) {
		if resultConfigs[ldapMapperDoguName] == nil {
			resultConfigs[ldapMapperDoguName] = map[string]string{"backend.password": wrces.configuration.UserBackend.Password,
				"backend.connection_dn": wrces.configuration.UserBackend.ConnectionDN}
		} else {
			resultConfigs[ldapMapperDoguName]["backend.password"] = wrces.configuration.UserBackend.Password
			resultConfigs[ldapMapperDoguName]["backend.connection_dn"] = wrces.configuration.UserBackend.ConnectionDN
		}
	}
}

func (wrces *writeRegistryConfigEncryptedStep) appendCasConfig(resultConfigs map[string]map[string]string) {
	if wrces.configuration.UserBackend.DsType == validation.DsTypeExternal {
		if resultConfigs["cas"] == nil {
			resultConfigs["cas"] = map[string]string{"password": wrces.configuration.UserBackend.Password}
		} else {
			resultConfigs["cas"]["password"] = wrces.configuration.UserBackend.Password
		}
	}
}

func (wrces *writeRegistryConfigEncryptedStep) appendLdapConfig(resultConfigs map[string]map[string]string) {
	if wrces.configuration.UserBackend.DsType != validation.DsTypeEmbedded {
		return
	}

	if resultConfigs["ldap"] == nil {
		resultConfigs["ldap"] = map[string]string{"admin_password": wrces.configuration.Admin.Password}
	} else {
		resultConfigs["ldap"]["admin_password"] = wrces.configuration.Admin.Password
	}
}

func (wrces *writeRegistryConfigEncryptedStep) createRegistryConfigEncryptedSecret(dogu string, stringData map[string]string) error {
	secretName := dogu + "-secrets"
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: wrces.namespace}, StringData: stringData}

	_, err := wrces.clientSet.CoreV1().Secrets(wrces.namespace).Create(gocontext.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create secret %s: %w", secretName, err)
	}

	return nil
}

// stringDataConfigurationWriter writes a configuration into a map used to set in secrets.
type stringDataConfigurationWriter struct{}

// NewStringDataConfigurationWriter creates a new instance of a map string data configuration write used for
// registry config encrypted
func NewStringDataConfigurationWriter() *stringDataConfigurationWriter {
	return &stringDataConfigurationWriter{}
}

// WriteConfigToStringDataMap write the given registry config to a map. It uses the delimiter '.' because the keys
// from the secret do not allow '/' in their data keys
func (mcw *stringDataConfigurationWriter) WriteConfigToStringDataMap(registryConfig context.CustomKeyValue) (map[string]map[string]string, error) {
	resultConfigs := map[string]map[string]string{}

	// build a string map for every key in registry config encrypted
	for config, entryMap := range registryConfig {
		resultConfig := map[string]string{}
		for key, value := range entryMap {
			contextWriter := &configWriter{delimiter: ".", write: func(field string, value string) error {
				resultConfig[field] = value
				return nil
			}}
			err := contextWriter.handleEntry(key, value)
			if err != nil {
				return nil, fmt.Errorf("failed to write %s config to map: %w", config, err)
			}
		}
		resultConfigs[config] = resultConfig
	}

	return resultConfigs, nil
}
