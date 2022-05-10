package validation

import (
	"fmt"
	"strconv"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

const (
	DsTypeEmbedded = "embedded"
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
		return GetInvalidOptionError("dsType", DsTypeEmbedded, DsTypeExternal)
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
		return GetInvalidOptionError("attributeID", "sAMAccountName")
	}
	if backend.AttributeFullname != "cn" {
		return GetInvalidOptionError("attributeFullName", "cn")
	}
	if backend.AttributeMail != "mail" {
		return GetInvalidOptionError("attributeMail", "mail")
	}
	if backend.AttributeGroup != "memberOf" {
		return GetInvalidOptionError("attributeGroup", "memberOf")
	}
	if backend.SearchFilter != searchFilter {
		return GetInvalidOptionError("searchFilter", searchFilter)
	}

	return nil
}

func (ubv *userBackendValidator) validateExternalBackend(backend context.UserBackend) error {
	if backend.Server != "activeDirectory" && backend.Server != "custom" {
		return GetInvalidOptionError("server", "activeDirectory", "custom")
	}
	if backend.Server == "activeDirectory" {
		err := ubv.validateActiveDirectoryServer(backend)
		if err != nil {
			return err
		}
	}
	if backend.AttributeGivenName == "" {
		return GetPropertyNotSetError("attributeGivenName")
	}
	if backend.AttributeSurname == "" {
		return GetPropertyNotSetError("attributeSurName")
	}
	if backend.BaseDN == "" {
		return GetPropertyNotSetError("baseDn")
	}
	if backend.ConnectionDN == "" {
		return GetPropertyNotSetError("connectionDn")
	}
	if backend.Password == "" {
		return GetPropertyNotSetError("password")
	}
	if backend.Host == "" {
		return GetPropertyNotSetError("host")
	}
	if backend.Port == "" {
		return GetPropertyNotSetError("port")
	}
	if _, err := strconv.Atoi(backend.Port); err != nil {
		return fmt.Errorf("failed to validate property port: the given value is not a number")
	}
	if backend.Encryption != "none" && backend.Encryption != "ssl" && backend.Encryption != "sslAny" && backend.Encryption != "startTLS" && backend.Encryption != "startTLSAny" {
		return GetInvalidOptionError("encryption", "none", "ssl", "sslAny", "startTLS", "startTLSAny")
	}
	if backend.GroupBaseDN == "" {
		return GetPropertyNotSetError("groupBaseDN")
	}
	if backend.GroupSearchFilter == "" {
		return GetPropertyNotSetError("groupSearchFilter")
	}
	if backend.GroupAttributeName == "" {
		return GetPropertyNotSetError("groupAttributeName")
	}
	if backend.GroupAttributeDescription == "" {
		return GetPropertyNotSetError("groupAttributeDescription")
	}
	if backend.GroupAttributeMember == "" {
		return GetPropertyNotSetError("groupAttributeMember")
	}

	return nil
}

func (ubv *userBackendValidator) validateEmbeddedBackend(backend context.UserBackend) error {
	if backend.AttributeID != "uid" {
		return GetInvalidOptionError("attributeID", "uid")
	}
	if backend.AttributeFullname != "cn" {
		return GetInvalidOptionError("attributeFullName", "cn")
	}
	if backend.AttributeMail != "mail" {
		return GetInvalidOptionError("attributeMail", "mail")
	}
	if backend.AttributeGroup != "memberOf" {
		return GetInvalidOptionError("attributeGroup", "memberOf")
	}
	if backend.SearchFilter != searchFilter {
		return GetInvalidOptionError("searchFilter", searchFilter)
	}
	if backend.Host != "ldap" {
		return GetInvalidOptionError("host", "ldap")
	}
	if backend.Port != "389" {
		return GetInvalidOptionError("port", "389")
	}

	return nil
}
