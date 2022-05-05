package validation

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"net"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type namingValidator struct {
}

// NewNamingValidator creates a new validator for the naming section of the setup configuration
func NewNamingValidator() *namingValidator {
	return &namingValidator{}
}

// ValidateNaming validates all properties of the naming section from a setup json
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (nv *namingValidator) ValidateNaming(naming context.Naming) error {
	ip := net.ParseIP(naming.Fqdn)
	domain := CheckDomain(naming.Fqdn)
	if ip == nil && domain != nil {
		return fmt.Errorf("failed to parse fqdn: %s", naming.Fqdn)
	}
	err := CheckDomain(naming.Domain)
	if err != nil {
		return fmt.Errorf("failed to validate domain: %w", err)
	}
	certificateType := naming.CertificateType
	if certificateType != "selfsigned" && certificateType != "external" {
		return GetInvalidOptionError("certificateType", "selfsigned", "external")
	}
	if certificateType == "external" {
		certificate := naming.Certificate
		certs := nv.splitPemCertificates(certificate)
		for _, cert := range certs {
			block, _ := pem.Decode([]byte(cert))
			if block == nil {
				return fmt.Errorf("failed to parse certificate PEM")
			}
			_, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %w", err)
			}
		}
		keyBlock, _ := pem.Decode([]byte(naming.CertificateKey))
		if keyBlock == nil {
			return fmt.Errorf("failed to parse private key PEM")
		}
	}
	err = CheckDomain(naming.RelayHost)
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

func (nv *namingValidator) splitPemCertificates(chain string) []string {
	sep := "-----BEGIN CERTIFICATE-----\n"
	result := []string{}
	split := strings.Split(chain, sep)
	for _, s := range split {
		if s == "" {
			continue
		}
		result = append(result, fmt.Sprintf("%s%s", sep, s))
	}
	return result
}

func GetInvalidOptionError(property string, validOptions ...string) error {
	return fmt.Errorf("invalid %s valid options are %s", property, validOptions)
}

func GetPropertyNotSetError(property string) error {
	return fmt.Errorf("no %s set", property)
}

// CheckDomain returns an error if the domain name is not valid
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
// https://gist.github.com/chmike/d4126a3247a6d9a70922fc0e8b4f4013
func CheckDomain(name string) error {
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
