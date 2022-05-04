package context

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/mail"
	"unicode/utf8"
)

type validator struct {
	configuration SetupConfiguration
}

func NewStartupConfigurationValidator(configuration SetupConfiguration) *validator {
	return &validator{configuration: configuration}
}

func (v *validator) ValidateConfiguration() error {
	err := v.validateNaming(v.configuration.Naming)
	if err != nil {
		return fmt.Errorf("failed to validate naming section: %w", err)
	}

	userBackend := v.configuration.UserBackend
	err = v.validateUserBackend(userBackend)
	if err != nil {
		return fmt.Errorf("failed to validate user userBackend section: %w", err)
	}

	if userBackend.DsType == "embedded" {
		err = v.validateAdminUser(v.configuration.Admin)
		if err != nil {
			return fmt.Errorf("failed to validate admin user section: %w", err)
		}
	}

	return nil
}

func (v *validator) validateUserBackend(backend UserBackend) error {
	dsType := backend.DsType
	if dsType != "embedded" && dsType != "external" {
		return fmt.Errorf("invalid user backend type %s valid options are embedded or external", dsType)
	}

	var result error
	server := backend.Server
	if server == "activeDirectory" {
		result = v.validateActiveDirectoryServer(backend)
	}
	if dsType == "external" {
		result = v.validateExternalBackend(backend)
	}
	if dsType == "embedded" {
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
		return fmt.Errorf("invalid uid %s valid option is sAMAccountName", id)
	}
	if fullName != "cn" {
		return fmt.Errorf("invalid attributeFullName %s valid option is cn", fullName)
	}
	if mail != "mail" {
		return fmt.Errorf("invalid atrtibuteMail %s valid option is mail", mail)
	}
	if group != "memberOf" {
		return fmt.Errorf("invalid attributeGroup %s valid option is memberOf", group)
	}
	if searchFilter != "(objectClass=person)" {
		return fmt.Errorf("invalid searchFilter %s valid option is (objectClass=person)", searchFilter)
	}

	return nil
}

func (v *validator) validateExternalBackend(backend UserBackend) error {
	host := backend.Host
	port := backend.Port
	server := backend.Server

	if server != "activeDirectory" && server != "custom" {
		return fmt.Errorf("invalid user backend server %s valid options are activeDirectory or custom", server)
	}
	if backend.AttributeGivenName == "" {
		return fmt.Errorf("no attributeGivenName set")
	}
	if backend.AttributeSurname == "" {
		return fmt.Errorf("no attributeSurName set")
	}
	if backend.BaseDN == "" {
		return fmt.Errorf("no baseDn set")
	}
	if backend.ConnectionDN == "" {
		return fmt.Errorf("no connectionDn set")
	}
	if backend.Password == "" {
		return fmt.Errorf("no password set")
	}
	if host == "" {
		return fmt.Errorf("no host set")
	}
	if port == "" {
		return fmt.Errorf("no port set")
	}
	encryption := backend.Encryption
	if encryption != "none" && encryption != "ssl" && encryption != "sslAny" && encryption != "startTLS" && encryption != "startTLSAny" {
		return fmt.Errorf("invalid encryption %s valid options are none, ssl, sslAny, startTLS or startTLSAny", encryption)
	}
	if backend.GroupBaseDN == "" {
		return fmt.Errorf("no groupBaseDN set")
	}
	if backend.GroupSearchFilter == "" {
		return fmt.Errorf("no groupSearchFilter set")
	}
	if backend.GroupAttributeName == "" {
		return fmt.Errorf("no groupAttributeName set")
	}
	if backend.GroupSearchFilter == "" {
		return fmt.Errorf("no groupAttributeDescription set")
	}
	if backend.GroupAttributeMember == "" {
		return fmt.Errorf("no groupAttributeMember set")
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
		return fmt.Errorf("invalid uid %s valid option is uid", id)
	}
	if fullName != "cn" {
		return fmt.Errorf("invalid attributeFullName %s valid option is cn", fullName)
	}
	if mail != "mail" {
		return fmt.Errorf("invalid atrtibuteMail %s valid option is mail", fullName)
	}
	if group != "memberOf" {
		return fmt.Errorf("invalid attributeGroup %s valid option is memberOf", group)
	}
	if searchFilter != "(objectClass=person)" {
		return fmt.Errorf("invalid searchFilter %s valid option is (objectClass=person)", searchFilter)
	}
	if host != "ldap" {
		return fmt.Errorf("invalid host %s valid option is ldap", searchFilter)
	}
	if port != "389" {
		return fmt.Errorf("invalid port %s valid option is 389", port)
	}

	return nil
}

func (v *validator) validateNaming(naming Naming) error {
	ip := net.ParseIP(naming.Fqdn)
	if ip == nil {
		return fmt.Errorf("failed to parse fqdn: %s", naming.Fqdn)
	}

	err := checkDomain(naming.Domain)
	if err != nil {
		return fmt.Errorf("failed to validate domain: %w", err)
	}

	certificateType := naming.CertificateType
	if certificateType != "selfsigned" && certificateType != "external" {
		return fmt.Errorf("invalid certification type %s valid options are selfsigned or external", certificateType)
	}

	if certificateType == "external" {
		block, _ := pem.Decode([]byte(naming.Certificate))
		if block == nil {
			return fmt.Errorf("failed to parse certificate PEM")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %w", err)
		}
		_, err = cert.Verify(x509.VerifyOptions{})
		if err != nil {
			return fmt.Errorf("failed to verfiy certificate: %w", err)
		}
		if naming.CertificateKey == "" {
			return fmt.Errorf("no certificate key")
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

func (v *validator) validateAdminUser(admin User) error {
	// TODO
	return nil
}

// https://gist.github.com/chmike/d4126a3247a6d9a70922fc0e8b4f4013
// checkDomain returns an error if the domain name is not valid
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
func checkDomain(name string) error {
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
