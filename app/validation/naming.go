package validation

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/mail"
	"strings"

	"github.com/cloudogu/k8s-ces-setup/v4/app/context"
)

type namingValidator struct{}

// NewNamingValidator creates a new validator for the naming section of the setup configuration
func NewNamingValidator() *namingValidator {
	return &namingValidator{}
}

// ValidateNaming validates all properties of the naming section from a setup json
// see: https://docs.cloudogu.com/docs/system-components/ces-setup/operations/setup-json_de/
func (nv *namingValidator) ValidateNaming(naming context.Naming) error {
	if naming.Fqdn == "" {
		return getPropertyNotSetError("fqdn")
	}

	if naming.Domain == "" {
		return getPropertyNotSetError("domain")
	}

	certificateType := naming.CertificateType
	if certificateType != "selfsigned" && certificateType != "external" {
		return getInvalidOptionError("certificateType", "selfsigned", "external")
	}

	if certificateType == "external" {
		err := validateCertificates(naming)
		if err != nil {
			return err
		}
	}

	if naming.RelayHost == "" {
		return getPropertyNotSetError("relayHost")
	}
	address := naming.MailAddress

	if address != "" {
		_, err := mail.ParseAddress(address)
		if err != nil {
			return fmt.Errorf("failed to validate mail address: %w", err)
		}
	}

	if naming.UseInternalIp {
		internalIP := naming.InternalIp
		ip := net.ParseIP(internalIP)
		if ip == nil {
			return fmt.Errorf("failed to parse internal ip: %s", internalIP)
		}
	}

	return nil
}

func validateCertificates(naming context.Naming) error {
	cert := naming.Certificate
	if cert == "" {
		return getPropertyNotSetError("certificate")
	}

	certs := SplitPemCertificates(cert)
	for i, cert := range certs {
		block, _ := pem.Decode([]byte(cert))
		if block == nil {
			return fmt.Errorf("failed to decode %d-th certificate in [certificate] property", i)
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse %d-th certificate in [certificate] property: %w", i, err)
		}
	}

	key := naming.CertificateKey
	if key == "" {
		return getPropertyNotSetError("certificate key")
	}

	keyBlock, _ := pem.Decode([]byte(key))
	if keyBlock == nil {
		return fmt.Errorf("failed to parse certificate key")
	}
	return nil
}

// SplitPemCertificates splits a certificate chain in pem format and returns all certificates of the chain as []string
func SplitPemCertificates(chain string) []string {
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

func getInvalidOptionError(property string, validOptions ...string) error {
	return fmt.Errorf("invalid %s valid options are %s", property, validOptions)
}

func getPropertyNotSetError(property string) error {
	return fmt.Errorf("no %s set", property)
}
