package validation

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"net/mail"
)

type adminValidator struct {
}

// NewAdminValidator creates a new validator for the admin section of the setup configuration
func NewAdminValidator() *adminValidator {
	return &adminValidator{}
}

// ValidateAdmin validates all properties of the admin section from a setup json
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (av *adminValidator) ValidateAdmin(admin context.User, dsType string) error {
	if admin.AdminGroup == "" {
		return GetPropertyNotSetError("admin group")
	}
	if dsType == dsTypeExternal {
		return nil
	}
	address := admin.Mail
	if address == "" {
		return GetPropertyNotSetError("admin mail")
	}
	_, err := mail.ParseAddress(address)
	if err != nil {
		return fmt.Errorf("invalid admin mail")
	}
	if admin.Username == "" {
		return GetPropertyNotSetError("admin username")
	}
	if admin.Password == "" {
		return GetPropertyNotSetError("admin password")
	}

	return nil
}
