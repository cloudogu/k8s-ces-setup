package context

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/mail"
	"unicode/utf8"
)

const (
	dsTypeEmbedded = "embedded"
	dsTypeExternal = "external"
)

type validator struct {
	configuration SetupConfiguration
}

// NewStartupConfigurationValidator creates a new setup json validator
func NewStartupConfigurationValidator(configuration SetupConfiguration) *validator {
	return &validator{configuration: configuration}
}

// ValidateConfiguration checks the section naming, user backend and user from the setup.json configuration
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (v *validator) ValidateConfiguration() error {
	//TODO Dogu Deps
	naming := v.configuration.Naming
	err := v.validateNaming(naming)
	if err != nil && naming.Completed {
		return fmt.Errorf("failed to validate naming section: %w", err)
	}

	userBackend := v.configuration.UserBackend
	err = v.validateUserBackend(userBackend)
	if err != nil && userBackend.Completed {
		return fmt.Errorf("failed to validate user userBackend section: %w", err)
	}

	admin := v.configuration.Admin
	err = v.validateAdminUser(admin, userBackend.DsType)
	if err != nil && admin.Completed {
		return fmt.Errorf("failed to validate admin user section: %w", err)
	}

	return nil
}

func (v *validator) validateUserBackend(backend UserBackend) error {
	dsType := backend.DsType
	if dsType != dsTypeEmbedded && dsType != dsTypeExternal {
		return fmt.Errorf("invalid user backend type %s valid options are embedded or external", dsType)
	}

	var result error
	server := backend.Server
	if server == "activeDirectory" {
		result = v.validateActiveDirectoryServer(backend)
	}
	if dsType == dsTypeExternal {
		result = v.validateExternalBackend(backend)
	}
	if dsType == dsTypeEmbedded {
		result = v.validateEmbeddedBackend(backend)
	}

	return result
}

func (v *validator) validateActiveDirectoryServer(backend UserBackend) error {
	id := backend.AttributeID
	fullName := backend.AttributeFullname
	mail := backend.AttributeMail
	group := backend.AttributeGroup
	searchFilter := backend.SearchFilter

	if id != "sAMAccountName" {
		return v.getInvalidOptionError("attributeID", "sAMAccountName")
	}
	if fullName != "cn" {
		return v.getInvalidOptionError("attributeFullName", "cn")
	}
	if mail != "mail" {
		return v.getInvalidOptionError("attributeMail", "mail")
	}
	if group != "memberOf" {
		return v.getInvalidOptionError("attributeGroup", "memberOf")
	}
	if searchFilter != "(objectClass=person)" {
		return v.getInvalidOptionError("searchFilter", "(objectClass=person)")
	}

	return nil
}

func (v *validator) validateExternalBackend(backend UserBackend) error {
	host := backend.Host
	port := backend.Port
	server := backend.Server

	if server != "activeDirectory" && server != "custom" {
		return v.getInvalidOptionError("server", "activeDirectory", "custom")
	}
	if backend.AttributeGivenName == "" {
		return v.getPropertyNotSetError("attributeGivenName")
	}
	if backend.AttributeSurname == "" {
		return v.getPropertyNotSetError("attributeSurName")
	}
	if backend.BaseDN == "" {
		return v.getPropertyNotSetError("baseDn")
	}
	if backend.ConnectionDN == "" {
		return v.getPropertyNotSetError("connectionDn")
	}
	if backend.Password == "" {
		return v.getPropertyNotSetError("password")
	}
	if host == "" {
		return v.getPropertyNotSetError("host")
	}
	if port == "" {
		return v.getPropertyNotSetError("port")
	}
	encryption := backend.Encryption
	if encryption != "none" && encryption != "ssl" && encryption != "sslAny" && encryption != "startTLS" && encryption != "startTLSAny" {
		return v.getInvalidOptionError("encryption", "none", "ssl", "sslAny", "startTLS", "startTLSAny")
	}
	if backend.GroupBaseDN == "" {
		return v.getPropertyNotSetError("groupBaseDN")
	}
	if backend.GroupSearchFilter == "" {
		return v.getPropertyNotSetError("groupSearchFilter")
	}
	if backend.GroupAttributeName == "" {
		return v.getPropertyNotSetError("groupAttributeName")
	}
	if backend.GroupAttributeDescription == "" {
		return v.getPropertyNotSetError("groupAttributeDescription")
	}
	if backend.GroupAttributeMember == "" {
		return v.getPropertyNotSetError("groupAttributeMember")
	}

	return nil
}

func (v *validator) validateEmbeddedBackend(backend UserBackend) error {
	id := backend.AttributeID
	fullName := backend.AttributeFullname
	mail := backend.AttributeMail
	group := backend.AttributeGroup
	searchFilter := backend.SearchFilter
	host := backend.Host
	port := backend.Port

	if id != "uid" {
		return v.getInvalidOptionError("attributeID", "uid")
	}
	if fullName != "cn" {
		return v.getInvalidOptionError("attributeFullName", "cn")
	}
	if mail != "mail" {
		return v.getInvalidOptionError("attributeMail", "mail")
	}
	if group != "memberOf" {
		return v.getInvalidOptionError("attributeGroup", "memberOf")
	}
	if searchFilter != "(objectClass=person)" {
		return v.getInvalidOptionError("searchFilter", "(objectClass=person)")
	}
	if host != "ldap" {
		return v.getInvalidOptionError("host", "ldap")
	}
	if port != "389" {
		return v.getInvalidOptionError("port", "389")
	}

	return nil
}

func (v *validator) validateNaming(naming Naming) error {
	ip := net.ParseIP(naming.Fqdn)
	domain := checkDomain(naming.Fqdn)
	if ip == nil && domain != nil {
		return fmt.Errorf("failed to parse fqdn: %s", naming.Fqdn)
	}
	err := checkDomain(naming.Domain)
	if err != nil {
		return fmt.Errorf("failed to validate domain: %w", err)
	}
	certificateType := naming.CertificateType
	if certificateType != "selfsigned" && certificateType != "external" {
		return v.getInvalidOptionError("certificateType", "selfsigned", "external")
	}
	if certificateType == "external" {
		block, _ := pem.Decode([]byte(naming.Certificate))
		if block == nil {
			return fmt.Errorf("failed to parse certificate PEM")
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %w", err)
		}
		keyBlock, _ := pem.Decode([]byte(naming.CertificateKey))
		if keyBlock == nil {
			return fmt.Errorf("failed to parse private key PEM")
		}
	}
	err = checkDomain(naming.RelayHost)
	if err != nil {
		return fmt.Errorf("failed to validate mail relay host: %w", err)
	}
	address := naming.MailAddress
	if address != "" {
		_, err = mail.ParseAddress(address)
		if err != nil {
			return fmt.Errorf("failed to validate mail address: %w", err)
		}
	}
	if naming.UseInternalIp {
		internalIP := naming.InternalIp
		ip = net.ParseIP(internalIP)
		if ip == nil {
			return fmt.Errorf("failed to parse internal ip: %s", internalIP)
		}
	}

	return nil
}

func (v *validator) validateAdminUser(admin User, dsType string) error {
	if admin.AdminGroup == "" {
		return v.getPropertyNotSetError("admin group")
	}
	if dsType == dsTypeExternal {
		return nil
	}
	address := admin.Mail
	if address == "" {
		return v.getPropertyNotSetError("admin mail")
	}

	_, err := mail.ParseAddress(address)
	if err != nil {
		return fmt.Errorf("invalid admin mail")
	}

	if admin.Username == "" {
		return v.getPropertyNotSetError("admin username")
	}
	if admin.Password == "" {
		return v.getPropertyNotSetError("admin password")
	}

	return nil
}

func (v *validator) getInvalidOptionError(property string, validOptions ...string) error {
	return fmt.Errorf("invalid %s valid options are %s", property, validOptions)
}

func (v *validator) getPropertyNotSetError(property string) error {
	return fmt.Errorf("no %s set", property)
}

// https://gist.github.com/chmike/d4126a3247a6d9a70922fc0e8b4f4013
// checkDomain returns an error if the domain name is not valid
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
func checkDomain(name string) error {
	if name == "" {
		return fmt.Errorf("domain is empty")
	}

	switch {
	case len(name) == 0:
		return nil // an empty domain name will result in a cookie without a domain restriction
	case len(name) > 255:
		return fmt.Errorf("cookie domain: name length is %d, can't exceed 255", len(name))
	}
	var l int
	for i := 0; i < len(name); i++ {
		b := name[i]
		if b == '.' {
			// check domain labels validity
			switch {
			case i == l:
				return fmt.Errorf("cookie domain: invalid character '%c' at offset %d: label can't begin with a period", b, i)
			case i-l > 63:
				return fmt.Errorf("cookie domain: byte length of label '%s' is %d, can't exceed 63", name[l:i], i-l)
			case name[l] == '-':
				return fmt.Errorf("cookie domain: label '%s' at offset %d begins with a hyphen", name[l:i], l)
			case name[i-1] == '-':
				return fmt.Errorf("cookie domain: label '%s' at offset %d ends with a hyphen", name[l:i], l)
			}
			l = i + 1
			continue
		}
		// test label character validity, note: tests are ordered by decreasing validity frequency
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			// show the printable unicode character starting at byte offset i
			c, _ := utf8.DecodeRuneInString(name[i:])
			if c == utf8.RuneError {
				return fmt.Errorf("cookie domain: invalid rune at offset %d", i)
			}
			return fmt.Errorf("cookie domain: invalid character '%c' at offset %d", c, i)
		}
	}
	// check top level domain validity
	switch {
	case l == len(name):
		return fmt.Errorf("cookie domain: missing top level domain, domain can't end with a period")
	case len(name)-l > 63:
		return fmt.Errorf("cookie domain: byte length of top level domain '%s' is %d, can't exceed 63", name[l:], len(name)-l)
	case name[l] == '-':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d begins with a hyphen", name[l:], l)
	case name[len(name)-1] == '-':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d ends with a hyphen", name[l:], l)
	case name[l] >= '0' && name[l] <= '9':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d begins with a digit", name[l:], l)
	}
	return nil
}
