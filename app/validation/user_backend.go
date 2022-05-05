package validation

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

const (
	dsTypeEmbedded = "embedded"
	dsTypeExternal = "external"
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
	if dsType != dsTypeEmbedded && dsType != dsTypeExternal {
		return fmt.Errorf("invalid user backend type %s valid options are embedded or external", dsType)
	}

	var result error
	server := backend.Server
	if server == "activeDirectory" {
		result = ubv.validateActiveDirectoryServer(backend)
	}
	if dsType == dsTypeExternal {
		result = ubv.validateExternalBackend(backend)
	}
	if dsType == dsTypeEmbedded {
		result = ubv.validateEmbeddedBackend(backend)
	}

	return result
}

func (ubv *userBackendValidator) validateActiveDirectoryServer(backend context.UserBackend) error {
	id := backend.AttributeID
	fullName := backend.AttributeFullname
	mail := backend.AttributeMail
	group := backend.AttributeGroup
	searchFilter := backend.SearchFilter

	if id != "sAMAccountName" {
		return GetInvalidOptionError("attributeID", "sAMAccountName")
	}
	if fullName != "cn" {
		return GetInvalidOptionError("attributeFullName", "cn")
	}
	if mail != "mail" {
		return GetInvalidOptionError("attributeMail", "mail")
	}
	if group != "memberOf" {
		return GetInvalidOptionError("attributeGroup", "memberOf")
	}
	if searchFilter != "(objectClass=person)" {
		return GetInvalidOptionError("searchFilter", "(objectClass=person)")
	}

	return nil
}

func (ubv *userBackendValidator) validateExternalBackend(backend context.UserBackend) error {
	host := backend.Host
	port := backend.Port
	server := backend.Server

	if server != "activeDirectory" && server != "custom" {
		return GetInvalidOptionError("server", "activeDirectory", "custom")
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
	if host == "" {
		return GetPropertyNotSetError("host")
	}
	if port == "" {
		return GetPropertyNotSetError("port")
	}
	encryption := backend.Encryption
	if encryption != "none" && encryption != "ssl" && encryption != "sslAny" && encryption != "startTLS" && encryption != "startTLSAny" {
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
	id := backend.AttributeID
	fullName := backend.AttributeFullname
	mail := backend.AttributeMail
	group := backend.AttributeGroup
	searchFilter := backend.SearchFilter
	host := backend.Host
	port := backend.Port

	if id != "uid" {
		return GetInvalidOptionError("attributeID", "uid")
	}
	if fullName != "cn" {
		return GetInvalidOptionError("attributeFullName", "cn")
	}
	if mail != "mail" {
		return GetInvalidOptionError("attributeMail", "mail")
	}
	if group != "memberOf" {
		return GetInvalidOptionError("attributeGroup", "memberOf")
	}
	if searchFilter != "(objectClass=person)" {
		return GetInvalidOptionError("searchFilter", "(objectClass=person)")
	}
	if host != "ldap" {
		return GetInvalidOptionError("host", "ldap")
	}
	if port != "389" {
		return GetInvalidOptionError("port", "389")
	}

	return nil
}
