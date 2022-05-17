package data

import (
	"fmt"
	"strconv"

	"github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/k8s-ces-setup/app/validation"
)

type writeAdminDataStep struct {
	Writer        RegistryWriter
	Configuration *context.SetupConfiguration
}

// NewWriteAdminDataStep create a new setup step which writes the admin data into the registry.
func NewWriteAdminDataStep(writer RegistryWriter, configuration *context.SetupConfiguration) *writeAdminDataStep {
	return &writeAdminDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wacs *writeAdminDataStep) GetStepDescription() string {
	return "Write admin data to the registry"
}

// PerformSetupStep writes the configured admin data into the registry.
func (wacs *writeAdminDataStep) PerformSetupStep() error {
	registryConfig := context.CustomKeyValue{
		"_global": map[string]interface{}{
			"admin_group": wacs.Configuration.Admin.AdminGroup,
		},
	}

	if wacs.Configuration.UserBackend.DsType == validation.DsTypeEmbedded {
		registryConfig["ldap"] = map[string]interface{}{
			"admin_username": wacs.Configuration.Admin.Username,
			"admin_mail":     wacs.Configuration.Admin.Mail,
			"admin_member":   strconv.FormatBool(wacs.Configuration.Admin.AdminMember),
		}
	}

	err := wacs.Writer.WriteConfigToRegistry(registryConfig)
	if err != nil {
		return fmt.Errorf("failed to write admin data to registry: %w", err)
	}

	return nil
}
