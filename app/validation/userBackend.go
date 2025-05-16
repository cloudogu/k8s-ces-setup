package validation

import (
	"fmt"
	"strconv"

	"github.com/cloudogu/k8s-ces-setup/v2/app/context"
)

const (
	// DsTypeEmbedded is used to set up the EcoSystem with an embedded user backend (e.g. the ldap dogu).
	DsTypeEmbedded = "embedded"
	// DsTypeExternal is used to set up the EcoSystem with an external user backend.
	DsTypeExternal = "external"
	searchFilter   = "(objectClass=person)"
)

type userBackendValidator struct {
}

// NewUserBackendValidator creates a new validator for the user backend section of the setup configuration
func NewUserBackendValidator() *userBackendValidator {
	return &userBackendValidator{}
}

// ValidateUserBackend validates all properties of the user backend section from a setup json
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (ubv *userBackendValidator) ValidateUserBackend(backend context.UserBackend) error {
	dsType := backend.DsType
	if dsType != DsTypeEmbedded && dsType != DsTypeExternal {
		return getInvalidOptionError("dsType", DsTypeEmbedded, DsTypeExternal)
	}

	var result error
	if dsType == DsTypeExternal {
		result = ubv.validateExternalBackend(backend)
	}
	if dsType == DsTypeEmbedded {
		result = ubv.validateEmbeddedBackend(backend)
	}

	return result
}

func (ubv *userBackendValidator) validateActiveDirectoryServer(backend context.UserBackend) error {
	if backend.AttributeID != "sAMAccountName" {
		return getInvalidOptionError("attributeID", "sAMAccountName")
	}
	if backend.AttributeFullname != "cn" {
		return getInvalidOptionError("attributeFullName", "cn")
	}
	if backend.AttributeMail != "mail" {
		return getInvalidOptionError("attributeMail", "mail")
	}
	if backend.AttributeGroup != "memberOf" {
		return getInvalidOptionError("attributeGroup", "memberOf")
	}
	if backend.SearchFilter != searchFilter {
		return getInvalidOptionError("searchFilter", searchFilter)
	}

	return nil
}

func (ubv *userBackendValidator) validateExternalBackend(backend context.UserBackend) error {
	err := ubv.validateBackendLocation(backend)
	if err != nil {
		return err
	}

	err = ubv.validateBackendAuth(backend)
	if err != nil {
		return err
	}

	err = ubv.validateBackendGroups(backend)
	if err != nil {
		return err
	}

	return nil
}

func (ubv *userBackendValidator) validateBackendLocation(backend context.UserBackend) error {
	if backend.Server != "activeDirectory" && backend.Server != "custom" {
		return getInvalidOptionError("server", "activeDirectory", "custom")
	}
	if backend.Server == "activeDirectory" {
		err := ubv.validateActiveDirectoryServer(backend)
		if err != nil {
			return err
		}
	}
	if backend.BaseDN == "" {
		return getPropertyNotSetError("baseDn")
	}
	if backend.ConnectionDN == "" {
		return getPropertyNotSetError("connectionDn")
	}
	if backend.Host == "" {
		return getPropertyNotSetError("host")
	}
	if backend.Port == "" {
		return getPropertyNotSetError("port")
	}
	if _, err := strconv.Atoi(backend.Port); err != nil {
		return fmt.Errorf("failed to validate property port: the given value is not a number")
	}
	if backend.Encryption != "none" && backend.Encryption != "ssl" && backend.Encryption != "sslAny" && backend.Encryption != "startTLS" && backend.Encryption != "startTLSAny" {
		return getInvalidOptionError("encryption", "none", "ssl", "sslAny", "startTLS", "startTLSAny")
	}

	return nil
}

func (ubv *userBackendValidator) validateBackendAuth(backend context.UserBackend) error {
	if backend.Password == "" {
		return getPropertyNotSetError("password")
	}
	if backend.AttributeGivenName == "" {
		return getPropertyNotSetError("attributeGivenName")
	}
	if backend.AttributeSurname == "" {
		return getPropertyNotSetError("attributeSurName")
	}

	return nil
}

func (ubv *userBackendValidator) validateBackendGroups(backend context.UserBackend) error {
	if backend.GroupBaseDN == "" {
		return getPropertyNotSetError("groupBaseDN")
	}
	if backend.GroupSearchFilter == "" {
		return getPropertyNotSetError("groupSearchFilter")
	}
	if backend.GroupAttributeName == "" {
		return getPropertyNotSetError("groupAttributeName")
	}
	if backend.GroupAttributeDescription == "" {
		return getPropertyNotSetError("groupAttributeDescription")
	}
	if backend.GroupAttributeMember == "" {
		return getPropertyNotSetError("groupAttributeMember")
	}

	return nil
}

func (ubv *userBackendValidator) validateEmbeddedBackend(backend context.UserBackend) error {
	if backend.AttributeID != "uid" {
		return getInvalidOptionError("attributeID", "uid")
	}
	if backend.AttributeFullname != "cn" {
		return getInvalidOptionError("attributeFullName", "cn")
	}
	if backend.AttributeMail != "mail" {
		return getInvalidOptionError("attributeMail", "mail")
	}
	if backend.AttributeGroup != "memberOf" {
		return getInvalidOptionError("attributeGroup", "memberOf")
	}
	if backend.SearchFilter != searchFilter {
		return getInvalidOptionError("searchFilter", searchFilter)
	}
	if backend.Host != "ldap" {
		return getInvalidOptionError("host", "ldap")
	}
	if backend.Port != "389" {
		return getInvalidOptionError("port", "389")
	}

	return nil
}
