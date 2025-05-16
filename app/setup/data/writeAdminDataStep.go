package data

import (
	"context"
	"fmt"
	"strconv"

	appcontext "github.com/cloudogu/k8s-ces-setup/v4/app/context"

	"github.com/cloudogu/k8s-ces-setup/v4/app/validation"
)

type writeAdminDataStep struct {
	Writer        RegistryWriter
	Configuration *appcontext.SetupJsonConfiguration
}

// NewWriteAdminDataStep create a new setup step which writes the admin data into the registry.
func NewWriteAdminDataStep(writer RegistryWriter, configuration *appcontext.SetupJsonConfiguration) *writeAdminDataStep {
	return &writeAdminDataStep{Writer: writer, Configuration: configuration}
}

// GetStepDescription return the human-readable description of the step.
func (wacs *writeAdminDataStep) GetStepDescription() string {
	return "Write admin data to the registry"
}

// PerformSetupStep writes the configured admin data into the registry.
func (wacs *writeAdminDataStep) PerformSetupStep(context.Context) error {
	registryConfig := appcontext.CustomKeyValue{
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
