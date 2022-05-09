package data

import (
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
	err := wctrs.writeDoguSection()
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

func (wctrs *writeConfigToRegistryStep) writeDoguSection() error {

	return nil
}

func (wctrs *writeConfigToRegistryStep) writeUserBackendSection() error {
	//log.Info("Setup user backend configuration")
	//
	//reg, err := registry.NewFromContext(cesappCtx)
	//if err != nil {
	//	return errors.Wrap(err, "failed to create registry connection")
	//}
	//
	//err = setupCasUserBackend(reg, conf, cesappCtx)
	//if err != nil {
	//	return errors.Wrap(err, "could not setup user backend for cas")
	//}
	//
	//err = setupLdapMapperUserBackend(reg, conf, cesappCtx)
	//if err != nil {
	//	return errors.Wrap(err, "could not setup user backend for ldap-mapper")
	//}
	return nil
}

//func setupCasUserBackend(reg registry.Registry, conf SetupConfiguration, cesappCtx *core.Context) error {
//	doguConfig := reg.DoguConfig("cas")
//
//	log.Debug("set directory service type in cas ldap-configuration")
//	err := doguConfig.Set("ldap/ds_type", conf.UserBackend.DsType)
//	if err != nil {
//		return errors.Wrap(err, "could not set directory service type")
//	}
//
//	log.Debug("set server in cas ldap-configuration")
//	err = doguConfig.Set("ldap/server", conf.UserBackend.Server)
//	if err != nil {
//		return errors.Wrap(err, "could not set server")
//	}
//
//	log.Debug("set ID attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_id", conf.UserBackend.AttributeID)
//	if err != nil {
//		return errors.Wrap(err, "could not set ID attribute name configuration")
//	}
//
//	log.Debug("set GivenName attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_given_name", conf.UserBackend.AttributeGivenName)
//	if err != nil {
//		return errors.Wrap(err, "could not set GivenName attribute name configuration")
//	}
//
//	log.Debug("set Surname attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_surname", conf.UserBackend.AttributeSurname)
//	if err != nil {
//		return errors.Wrap(err, "could not set Surname attribute name configuration")
//	}
//
//	log.Debug("set Fullname attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_fullname", conf.UserBackend.AttributeFullname)
//	if err != nil {
//		return errors.Wrap(err, "could not set Fullname attribute name configuration")
//	}
//
//	log.Debug("set Mail attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_mail", conf.UserBackend.AttributeMail)
//	if err != nil {
//		return errors.Wrap(err, "could not set Mail attribute name configuration")
//	}
//
//	log.Debug("set Group attribute name in cas ldap-configuration")
//	err = doguConfig.Set("ldap/attribute_group", conf.UserBackend.AttributeGroup)
//	if err != nil {
//		return errors.Wrap(err, "could not set Group attribute name configuration")
//	}
//
//	log.Debug("set group base dn cas ldap-configuration")
//	err = doguConfig.Set("ldap/group_base_dn", conf.UserBackend.GroupBaseDN)
//	if err != nil {
//		return errors.Wrap(err, "could not set group base dn configuration")
//	}
//
//	log.Debug("set group search filter in cas ldap-configuration")
//	err = doguConfig.Set("ldap/group_search_filter", conf.UserBackend.GroupSearchFilter)
//	if err != nil {
//		return errors.Wrap(err, "could not set group search filter configuration")
//	}
//
//	log.Debug("set group name attribute in cas ldap-configuration")
//	err = doguConfig.Set("ldap/group_attribute_name", conf.UserBackend.GroupAttributeName)
//	if err != nil {
//		return errors.Wrap(err, "could not group name attribute name configuration")
//	}
//
//	log.Debug("set BaseDN in cas ldap-configuration")
//	err = doguConfig.Set("ldap/base_dn", conf.UserBackend.BaseDN)
//	if err != nil {
//		return errors.Wrap(err, "could not set BaseDN")
//	}
//
//	log.Debug("set search filter in cas ldap-configuration")
//	err = doguConfig.Set("ldap/search_filter", conf.UserBackend.SearchFilter)
//	if err != nil {
//		return errors.Wrap(err, "could not set search Filter")
//	}
//
//	log.Debug("set ConnectionDN in cas ldap-configuration")
//	err = doguConfig.Set("ldap/connection_dn", conf.UserBackend.ConnectionDN)
//	if err != nil {
//		return errors.Wrap(err, "could not set ConnectionDN")
//	}
//
//	encryption := conf.UserBackend.Encryption
//	if encryption == "" {
//		log.Debug("skip setting empty encryption in cas ldap-configuration")
//	} else {
//		log.Debug("set encryption in cas ldap-configuration")
//		err = doguConfig.Set("ldap/encryption", encryption)
//		if err != nil {
//			return errors.Wrap(err, "could not set encryption")
//		}
//	}
//
//	log.Debug("set encrypted password in cas ldap-configuration")
//	casPublicKeyString, err := doguConfig.Get("public.pem")
//	if err == nil {
//		keyProvider, err := keys.NewKeyProviderFromContext(cesappCtx)
//		if err != nil {
//			return errors.Wrap(err, "could not create key provider")
//		}
//
//		casPublicKey, err := keyProvider.ReadPublicKeyFromString(casPublicKeyString)
//		if err != nil {
//			return errors.Wrap(err, "could not get public key from public.pem")
//		}
//
//		passwordEnc, err := casPublicKey.Encrypt(conf.UserBackend.Password)
//		if err != nil {
//			return errors.Wrap(err, "could not encrypt password")
//		}
//
//		err = doguConfig.Set("ldap/password", passwordEnc)
//		if err != nil {
//			return errors.Wrap(err, "could not set password")
//		}
//	} else {
//		// maybe cas is not installed, continue execution
//		log.Warning("error while trying to get public.pem from cas: " + err.Error())
//	}
//
//	log.Debug("set host in cas ldap-configuration")
//	err = doguConfig.Set("ldap/host", conf.UserBackend.Host)
//	if err != nil {
//		return errors.Wrap(err, "could not set host")
//	}
//
//	log.Debug("set port in cas ldap-configuration")
//	err = doguConfig.Set("ldap/port", conf.UserBackend.Port)
//	if err != nil {
//		return errors.Wrap(err, "could not set port")
//	}
//	return nil
//}
//
//func setupLdapMapperUserBackend(reg registry.Registry, conf SetupConfiguration, cesappCtx *core.Context) error {
//
//	enabled, err := reg.DoguRegistry().IsEnabled("ldap-mapper")
//	if err != nil {
//		return errors.Wrap(err, "failed to check if ldap-mapper is enabled")
//	}
//
//	if !enabled {
//		log.Info("Skip setupLdapMapperUserBackend since ldap-mapper isn't installed")
//		return nil
//	}
//
//	doguConfig := reg.DoguConfig("ldap-mapper")
//
//	log.Debug("set directory service type in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("backend/type", conf.UserBackend.DsType)
//	if err != nil {
//		return errors.Wrap(err, "could not set directory service type")
//	}
//
//	log.Debug("set host in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("backend/host", conf.UserBackend.Host)
//	if err != nil {
//		return errors.Wrap(err, "could not set host")
//	}
//
//	log.Debug("set port in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("backend/port", conf.UserBackend.Port)
//	if err != nil {
//		return errors.Wrap(err, "could not set port")
//	}
//
//	// When using embedded ldap, we do not have to use the following config
//	if conf.UserBackend.DsType == UserBackendDSTypeEmbedded {
//		return nil
//	}
//
//	log.Debug("set BaseDN in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/base_dn", conf.UserBackend.BaseDN)
//	if err != nil {
//		return errors.Wrap(err, "could not set BaseDN")
//	}
//
//	log.Debug("set search filter in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/search_filter", conf.UserBackend.SearchFilter)
//	if err != nil {
//		return errors.Wrap(err, "could not set search Filter")
//	}
//
//	log.Debug("set ID attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/id", conf.UserBackend.AttributeID)
//	if err != nil {
//		return errors.Wrap(err, "could not set ID attribute name configuration")
//	}
//
//	log.Debug("set GivenName attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/given_name", conf.UserBackend.AttributeGivenName)
//	if err != nil {
//		return errors.Wrap(err, "could not set GivenName attribute name configuration")
//	}
//
//	log.Debug("set Surname attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/surname", conf.UserBackend.AttributeSurname)
//	if err != nil {
//		return errors.Wrap(err, "could not set Surname attribute name configuration")
//	}
//
//	log.Debug("set Fullname attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/full_name", conf.UserBackend.AttributeFullname)
//	if err != nil {
//		return errors.Wrap(err, "could not set Fullname attribute name configuration")
//	}
//
//	log.Debug("set Mail attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/mail", conf.UserBackend.AttributeMail)
//	if err != nil {
//		return errors.Wrap(err, "could not set Mail attribute name configuration")
//	}
//
//	log.Debug("set Group attribute name in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/user/group", conf.UserBackend.AttributeGroup)
//	if err != nil {
//		return errors.Wrap(err, "could not set Group attribute name configuration")
//	}
//
//	log.Debug("set group base dn ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/group/base_dn", conf.UserBackend.GroupBaseDN)
//	if err != nil {
//		return errors.Wrap(err, "could not set group base dn configuration")
//	}
//
//	log.Debug("set group search filter in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/group/search_filter", conf.UserBackend.GroupSearchFilter)
//	if err != nil {
//		return errors.Wrap(err, "could not set group search filter configuration")
//	}
//
//	log.Debug("set group name attribute in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/group/name", conf.UserBackend.GroupAttributeName)
//	if err != nil {
//		return errors.Wrap(err, "could not group name attribute name configuration")
//	}
//
//	log.Debug("set group description attribute in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/group/description", conf.UserBackend.GroupAttributeDescription)
//	if err != nil {
//		return errors.Wrap(err, "could not group description attribute description configuration")
//	}
//
//	log.Debug("set group member attribute in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("mapping/group/member", conf.UserBackend.GroupAttributeMember)
//	if err != nil {
//		return errors.Wrap(err, "could not group member attribute member configuration")
//	}
//
//	log.Debug("set encryption in ldap-mapper ldap-configuration")
//	err = doguConfig.Set("backend/encryption", conf.UserBackend.Encryption)
//	if err != nil {
//		return errors.Wrap(err, "could not set encryption")
//	}
//
//	log.Debug("set server in ldap ldap-configuration")
//	err = doguConfig.Set("backend/server", conf.UserBackend.Server)
//	if err != nil {
//		return errors.Wrap(err, "could not set server")
//	}
//
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
//
//	return nil
//}

func (wctrs *writeConfigToRegistryStep) writeRegistryConfigSection() error {

	return nil
}

//// writes registryConfig and registryConfigEncrypted into the registry
//func writeRegistryConfig(cesappCtx *core.Context, conf SetupConfiguration) error {
//	writer, err := newSetupConfigWriter(cesappCtx, false)
//	if err != nil {
//		return errors.Wrapf(err, "failed to create setup config writer")
//	}
//
//	err = writer.writeConfigToRegistry(conf.RegistryConfig)
//	if err != nil {
//		return errors.Wrap(err, "failed to write config in registry")
//	}
//
//	writer, err = newSetupConfigWriter(cesappCtx, true)
//	if err != nil {
//		return errors.Wrapf(err, "failed to create setup config writer")
//	}
//
//	err = writer.writeConfigToRegistry(conf.RegistryConfigEncrypted)
//	if err != nil {
//		return errors.Wrap(err, "failed to write encrypted config in registry")
//	}
//
//	return nil
//}
//
//func newSetupConfigWriter(context *core.Context, encrypted bool) (*ConfigWriter, error) {
//	reg, err := registry.NewFromContext(context)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to create registry connection")
//	}
//	return &ConfigWriter{context, reg, encrypted}, nil
//}
//
//func (writer *ConfigWriter) writeConfigToRegistry(registryConfig CustomKeyValue) error {
//	for fieldName, fieldMap := range registryConfig {
//		err := writer.writeEntriesForConfig(fieldMap, fieldName)
//		if err != nil {
//			return errors.Wrap(err, "failed to write entries in registry config")
//		}
//	}
//	return nil
//}
//
//func (writer *ConfigWriter) writeEntriesForConfig(entries map[string]interface{}, config string) error {
//	var configCtx registry.ConfigurationContext
//
//	log.Infof("write in %s configuration", config)
//	if config == "_global" {
//		if writer.encrypted {
//			return fmt.Errorf("encrypted entries not possible for _global configuration")
//		}
//		configCtx = writer.registry.GlobalConfig()
//	} else {
//		configCtx = writer.registry.DoguConfig(config)
//	}
//
//	contextWriter, err := writer.GetContextConfigWriter(configCtx)
//	if err != nil {
//		return errors.Wrapf(err, "could not get registry config")
//	}
//
//	for fieldName, fieldEntry := range entries {
//		err := contextWriter.handleEntry(fieldName, fieldEntry)
//		if err != nil {
//			return errors.Wrapf(err, "failed to write %s in registry config", fieldName)
//		}
//	}
//
//	return nil
//}
