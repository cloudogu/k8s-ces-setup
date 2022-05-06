package data_test

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGenerateSSL(t *testing.T) {
	// when
	step := data.NewGenerateSSLStep(nil)

	// then
	require.NotNil(t, step)
}

func Test_generateSSLStep_GetStepDescription(t *testing.T) {
	// given
	step := data.NewGenerateSSLStep(nil)

	// when
	description := step.GetStepDescription()

	// then
	assert.Equal(t, "Generate SSL certificate and key", description)
}

func Test_generateSSLStep_PerformSetupStep(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{CertificateType: "selfsigned"}}
		step := data.NewGenerateSSLStep(config)

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, config.Naming.Certificate)
		assert.NotEmpty(t, config.Naming.CertificateKey)

		certs := validation.SplitPemCertificates(config.Naming.Certificate)
		assert.Equal(t, 2, len(certs))

		err = validateCert(certs[0])
		require.NoError(t, err)
		err = validateCert(certs[1])
		require.NoError(t, err)
		_, err = validatePEM(config.Naming.CertificateKey)
		require.NoError(t, err)
	})

	t.Run("let external cert unchanged", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{CertificateType: "external", Certificate: "bitte nicht", CertificateKey: "bitte nicht"}}
		step := data.NewGenerateSSLStep(config)

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		assert.Equal(t, "bitte nicht", config.Naming.Certificate)
		assert.Equal(t, "bitte nicht", config.Naming.CertificateKey)
	})
}

func validateCert(cert string) error {
	block, err := validatePEM(cert)
	_, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return assert.AnError
	}

	return nil
}

func validatePEM(pemStr string) (*pem.Block, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, assert.AnError
	}

	return block, nil
}
