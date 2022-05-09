package data

import (
	"fmt"
	"strconv"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/k8s-ces-setup/app/validation"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type writeAdminConfigStep struct {
	Registry      registry.Registry
	Configuration *context.SetupConfiguration
}

// NewWriteAdminConfigStep create a new setup step which writes the admin configuration into the registry.
func NewWriteAdminConfigStep(registry registry.Registry, configuration *context.SetupConfiguration) *writeAdminConfigStep {
	return &writeAdminConfigStep{Registry: registry, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wctrs *writeAdminConfigStep) GetStepDescription() string {
	return "Write admin configuration to the registry"
}

func (wctrs *writeAdminConfigStep) PerformSetupStep() error {
	logrus.Info("Write admin configuration into registry")

	err := wctrs.Registry.GlobalConfig().Set("admin_group", wctrs.Configuration.Admin.AdminGroup)
	if err != nil {
		return errors.Wrap(err, "could not set admin group")
	}

	if wctrs.Configuration.UserBackend.DsType == validation.DsTypeEmbedded {

		ldapConfig := map[string]string{
			"admin_username": wctrs.Configuration.Admin.Username,
			"admin_mail":     wctrs.Configuration.Admin.Mail,
			"admin_member":   strconv.FormatBool(wctrs.Configuration.Admin.AdminMember),
		}

		err = writeMapToContext(wctrs.Registry.DoguConfig("ldap"), "ldap/config", ldapConfig)
		if err != nil {
			return err
		}

		// TODO: Currently it is not possible to save encrypted values for the dogus. Should be addressed ASAP
		//logrus.Debug("set encrypted admin password in ldap dogu configuration")
		//ldapPublicKeyString, err := ldapContext.Get("public.pem")
		//if err != nil {
		//	return errors.Wrap(err, "could not get public.pem from ldap")
		//}
		//
		//keyProvider, err := keys.NewKeyProviderFromContext(cesappCtx)
		//if err != nil {
		//	return errors.Wrap(err, "could not create key provider")
		//}
		//
		//ldapPublicKey, err := keyProvider.ReadPublicKeyFromString(ldapPublicKeyString)
		//if err != nil {
		//	return errors.Wrap(err, "could not get public key from public.pem")
		//}
		//
		//passwordEnc, err := ldapPublicKey.Encrypt(wctrs.Configuration.Admin.Password)
		//if err != nil {
		//	return errors.Wrap(err, "could not encrypt password")
		//}
		//
		//err = ldapContext.Set("admin_password", passwordEnc)
		//if err != nil {
		//	return errors.Wrap(err, "could not set admin password")
		//}
	}

	return nil
}

func writeMapToContext(context registry.ConfigurationContext, scope string, configMap map[string]string) error {
	for key, value := range configMap {
		err := context.Set(key, value)
		if err != nil {
			return fmt.Errorf("could not write [%s] into etcd registry: %w", key, err)
		}
	}

	return nil
}
