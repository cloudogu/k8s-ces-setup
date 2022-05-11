package data

import (
	"fmt"
	"strings"

	"github.com/cloudogu/k8s-ces-setup/app/validation"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

type writeLdapDataStep struct {
	Writer        RegistryWriter
	Configuration *context.SetupConfiguration
}

// NewWriteLdapDataStep create a new setup step which writes the ldap configuration into the registry.
func NewWriteLdapDataStep(writer RegistryWriter, configuration *context.SetupConfiguration) *writeLdapDataStep {
	return &writeLdapDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wlds *writeLdapDataStep) GetStepDescription() string {
	return "Write ldap data to the registry"
}

// PerformSetupStep writes the configured ldap data into the registry
func (wlds *writeLdapDataStep) PerformSetupStep() error {
	registryConfig := context.CustomKeyValue{}
	registryConfig["cas"] = wlds.getCasEntriesAsMap()

	if isDoguInstalled(wlds.Configuration.Dogus.Install, "ldap-mapper") {
		registryConfig["ldap-mapper"] = wlds.getLdapMapperEntriesAsMap()
	}

	err := wlds.Writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write ldap data to registry: %w", err)
	}

	return nil
}

func (wlds *writeLdapDataStep) getCasEntriesAsMap() map[string]interface{} {
	ldapInCasOptions := map[string]string{
		"ds_type":              wlds.Configuration.UserBackend.DsType,
		"server":               wlds.Configuration.UserBackend.Server,
		"attribute_id":         wlds.Configuration.UserBackend.AttributeID,
		"attribute_given_name": wlds.Configuration.UserBackend.AttributeGivenName,
		"attribute_fullname":   wlds.Configuration.UserBackend.AttributeFullname,
		"attribute_mail":       wlds.Configuration.UserBackend.AttributeMail,
		"attribute_group":      wlds.Configuration.UserBackend.AttributeGroup,
		"group_base_dn":        wlds.Configuration.UserBackend.GroupBaseDN,
		"group_search_filter":  wlds.Configuration.UserBackend.GroupSearchFilter,
		"group_attribute_name": wlds.Configuration.UserBackend.GroupAttributeName,
		"base_dn":              wlds.Configuration.UserBackend.BaseDN,
		"search_filter":        wlds.Configuration.UserBackend.SearchFilter,
		"connection_dn":        wlds.Configuration.UserBackend.ConnectionDN,
		"host":                 wlds.Configuration.UserBackend.Host,
		"port":                 wlds.Configuration.UserBackend.Port,
	}

	if wlds.Configuration.UserBackend.Encryption != "" {
		ldapInCasOptions["encryption"] = wlds.Configuration.UserBackend.Encryption
	}

	// TODO save encrypted password in cas ldap-configuration
	//log.Debug("set encrypted password in cas ldap-configuration")
	//casPublicKeyString, err := doguConfig.Get("public.pem")
	//if err == nil {
	//	keyProvider, err := keys.NewKeyProviderFromContext(cesappCtx)
	//	if err != nil {
	//		return errors.Wrap(err, "could not create key provider")
	//	}
	//
	//	casPublicKey, err := keyProvider.ReadPublicKeyFromString(casPublicKeyString)
	//	if err != nil {
	//		return errors.Wrap(err, "could not get public key from public.pem")
	//	}
	//
	//	passwordEnc, err := casPublicKey.Encrypt(conf.UserBackend.Password)
	//	if err != nil {
	//		return errors.Wrap(err, "could not encrypt password")
	//	}
	//
	//	err = doguConfig.Set("ldap/password", passwordEnc)
	//	if err != nil {
	//		return errors.Wrap(err, "could not set password")
	//	}
	//} else {
	//	// maybe cas is not installed, continue execution
	//	log.Warning("error while trying to get public.pem from cas: " + err.Error())
	//}

	return map[string]interface{}{"ldap": ldapInCasOptions}
}

func (wlds *writeLdapDataStep) getLdapMapperEntriesAsMap() map[string]interface{} {
	ldapMapperRegistryConfig := map[string]interface{}{
		"backend": map[string]string{
			"type": wlds.Configuration.UserBackend.DsType,
			"host": wlds.Configuration.UserBackend.Host,
			"port": wlds.Configuration.UserBackend.Port,
		},
	}

	if wlds.Configuration.UserBackend.DsType != validation.DsTypeEmbedded {
		ldapMapperRegistryConfig["mapping"] = map[string]interface{}{
			"user": map[string]string{
				"base_dn":       wlds.Configuration.UserBackend.BaseDN,
				"search_filter": wlds.Configuration.UserBackend.SearchFilter,
				"id":            wlds.Configuration.UserBackend.AttributeID,
				"given_name":    wlds.Configuration.UserBackend.AttributeGivenName,
				"surname":       wlds.Configuration.UserBackend.AttributeSurname,
				"full_name":     wlds.Configuration.UserBackend.AttributeFullname,
				"mail":          wlds.Configuration.UserBackend.AttributeMail,
				"group":         wlds.Configuration.UserBackend.AttributeGroup,
			},
			"group": map[string]string{
				"base_dn":       wlds.Configuration.UserBackend.GroupBaseDN,
				"search_filter": wlds.Configuration.UserBackend.GroupSearchFilter,
				"name":          wlds.Configuration.UserBackend.GroupAttributeName,
				"description":   wlds.Configuration.UserBackend.GroupAttributeDescription,
				"member":        wlds.Configuration.UserBackend.GroupAttributeMember,
				"encryption":    wlds.Configuration.UserBackend.Encryption,
				"server":        wlds.Configuration.UserBackend.Server,
			},
		}
	}

	// TODO set encrypted password in ldap-mapper ldap-configuraion
	//	log.Debug("set encrypted password in ldap-mapper ldap-configuraion")
	//	ldapMapperPublicKeyString, err := doguConfig.Get("public.pem")
	//	if err == nil {
	//		keyProvider, err := keys.NewKeyProviderFromContext(cesappCtx)
	//		if err != nil {
	//			return errors.Wrap(err, "could not create key provider")
	//		}
	//
	//		ldapMapperPublicKey, err := keyProvider.ReadPublicKeyFromString(ldapMapperPublicKeyString)
	//		if err != nil {
	//			return errors.Wrap(err, "could not get public key from public.pem")
	//		}
	//
	//		passwordEnc, err := ldapMapperPublicKey.Encrypt(conf.UserBackend.Password)
	//		if err != nil {
	//			return errors.Wrap(err, "could not encrypt password")
	//		}
	//
	//		connectionDnEnc, err := ldapMapperPublicKey.Encrypt(conf.UserBackend.ConnectionDN)
	//		if err != nil {
	//			return errors.Wrap(err, "could not encrypt password")
	//		}
	//
	//		err = doguConfig.Set("backend/password", passwordEnc)
	//		if err != nil {
	//			return errors.Wrap(err, "could not set password")
	//		}
	//
	//		err = doguConfig.Set("backend/connection_dn", connectionDnEnc)
	//		if err != nil {
	//			return errors.Wrap(err, "could not set username")
	//		}
	//	} else {
	//		log.Warning("error while trying to get public.pem from ldap-mapper: " + err.Error())
	//	}

	return ldapMapperRegistryConfig
}

func isDoguInstalled(dogus []string, doguName string) bool {
	for _, doguNameWithVersion := range dogus {
		if strings.Contains(doguNameWithVersion, doguName) {
			return true
		}
	}
	return false
}
