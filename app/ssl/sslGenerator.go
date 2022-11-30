package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

type sslGenerator struct {
}

// NewSSLGenerator creates a new sslGenerator instance to generate a self-signed cert and key
func NewSSLGenerator() *sslGenerator {
	return &sslGenerator{}
}

// GenerateSelfSignedCert generates a self-signed certificate for the ces and returns the certificate chain and the
// private key as string
func (sg *sslGenerator) GenerateSelfSignedCert(fqdn string, domain string, certExpireDays int) (string, string, error) {
	// create x509 certificate and key (ca)
	ca, caPrivateKey, err := sg.createCertTemplateWithKey(certExpireDays, domain)
	if err != nil {
		return "", "", fmt.Errorf("failed to create ca cert: %w", err)
	}
	sg.appendCaTemplate(ca)
	caCertBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return "", "", err
	}

	// create x509 certificate and key
	cert, certPrivKey, err := sg.createCertTemplateWithKey(certExpireDays, domain)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cert: %w", err)
	}
	sg.appendCertTemplate(cert, fqdn)

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivateKey)
	if err != nil {
		return "", "", err
	}

	formattedEncodeErrorMsg := "failed to encode certificate: %w"
	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf(formattedEncodeErrorMsg, err)
	}

	caCertPEM := new(bytes.Buffer)
	err = pem.Encode(caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf(formattedEncodeErrorMsg, err)
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return "", "", fmt.Errorf(formattedEncodeErrorMsg, err)
	}

	chain := fmt.Sprintf("%s%s", certPEM.String(), caCertPEM.String())
	return chain, certPrivKeyPEM.String(), nil
}

func (sg *sslGenerator) createCertTemplateWithKey(validDays int, domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"DE"},
			Province:           []string{"Lower Saxony"},
			Locality:           []string{"Brunswick"},
			Organization:       []string{domain},
			OrganizationalUnit: []string{domain},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 0, validDays),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	return certificate, certPrivKey, nil
}

func (sg *sslGenerator) appendCaTemplate(ca *x509.Certificate) {
	ca.Subject.CommonName = "CES Self Signed"
	ca.IsCA = true
	ca.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	ca.BasicConstraintsValid = true
}

func (sg *sslGenerator) appendCertTemplate(cert *x509.Certificate, fqdn string) {
	cert.KeyUsage = x509.KeyUsageDigitalSignature
	cert.SubjectKeyId = []byte{1, 2, 3, 4, 6}
	ip := net.ParseIP(fqdn)
	if ip != nil {
		cert.IPAddresses = append(cert.IPAddresses, ip)
	}
	cert.DNSNames = append(cert.DNSNames, fqdn)
	cert.DNSNames = append(cert.DNSNames, "local.cloudogu.com")
	cert.Subject.CommonName = fqdn
}
