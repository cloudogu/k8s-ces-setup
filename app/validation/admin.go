package validation

import (
	"fmt"
	"net/mail"

	"github.com/cloudogu/k8s-ces-setup/v4/app/context"
)

type adminValidator struct{}

// NewAdminValidator creates a new validator for the admin section of the setup configuration
func NewAdminValidator() *adminValidator {
	return &adminValidator{}
}

// ValidateAdmin validates all properties of the admin section from a setup json
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (av *adminValidator) ValidateAdmin(admin context.User, dsType string) error {
	if admin.AdminGroup == "" {
		return getPropertyNotSetError("admin group")
	}

	if dsType == DsTypeExternal {
		return nil
	}

	if admin.Mail == "" {
		return getPropertyNotSetError("admin mail")
	}
	_, err := mail.ParseAddress(admin.Mail)
	if err != nil {
		return fmt.Errorf("invalid admin mail")
	}
	if admin.Username == "" {
		return getPropertyNotSetError("admin username")
	}
	if admin.Password == "" {
		return getPropertyNotSetError("admin password")
	}

	return nil
}
