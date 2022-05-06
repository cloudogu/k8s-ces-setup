package data

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"math/big"
	"net"
	"time"
)

const (
	certExpireDays = 365
)

type generateSSLStep struct {
	config *context.SetupConfiguration
}

// NewGenerateSSLStep create a new setup step which on generates ssl certificates
func NewGenerateSSLStep(config *context.SetupConfiguration) *generateSSLStep {
	return &generateSSLStep{config: config}
}

// GetStepDescription return the human-readable description of the step
func (gss *generateSSLStep) GetStepDescription() string {
	return fmt.Sprintf("Generate SSL certificate and key")
}

// PerformSetupStep either generates a certificate if necessary and writes it to the setup configuration
func (gss *generateSSLStep) PerformSetupStep() error {
	naming := gss.config.Naming
	// Generation not needed
	if naming.CertificateType == "external" {
		return nil
	}

	// create x509 certificate and key (ca)
	ca, caPrivateKey, err := gss.createCertTemplateWithKey(certExpireDays)
	if err != nil {
		return fmt.Errorf("failed to create ca cert: %w", err)
	}
	gss.appendCaTemplate(ca)
	caCertBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	// create x509 certificate and key
	cert, certPrivKey, err := gss.createCertTemplateWithKey(certExpireDays)
	if err != nil {
		return fmt.Errorf("failed to create cert: %w", err)
	}
	gss.appendCertTemplate(cert, naming)

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	caCertPEM := new(bytes.Buffer)
	err = pem.Encode(caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	gss.config.Naming.Certificate = fmt.Sprintf("%s%s", caCertPEM.String(), certPEM.String())
	gss.config.Naming.CertificateKey = certPrivKeyPEM.String()

	return nil
}

func (gss *generateSSLStep) createCertTemplateWithKey(validDays int) (*x509.Certificate, *rsa.PrivateKey, error) {
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
			Country:      []string{"DE"},
			Province:     []string{"Lower Saxony"},
			Locality:     []string{"Brunswick"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 0, validDays),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	return certificate, certPrivKey, nil
}

func (gss *generateSSLStep) appendCaTemplate(ca *x509.Certificate) {
	ca.Subject.CommonName = "CES Self Signed"
	ca.IsCA = true
	ca.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	ca.BasicConstraintsValid = true
}

func (gss *generateSSLStep) appendCertTemplate(cert *x509.Certificate, naming context.Naming) {
	cert.KeyUsage = x509.KeyUsageDigitalSignature
	cert.SubjectKeyId = []byte{1, 2, 3, 4, 6}
	fqdn := naming.Fqdn
	ip := net.ParseIP(fqdn)
	if ip != nil {
		cert.IPAddresses = append(cert.IPAddresses, ip)
	}
	cert.DNSNames = append(cert.DNSNames, fqdn)
	cert.DNSNames = append(cert.DNSNames, "local.cloudogu.com")
	cert.Subject.CommonName = fqdn

	if naming.UseInternalIp {
		ip := net.ParseIP(naming.InternalIp)
		if ip != nil {
			cert.IPAddresses = append(cert.IPAddresses, ip)
		}
	}
}
