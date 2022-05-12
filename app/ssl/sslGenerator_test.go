package ssl_test

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/cloudogu/k8s-ces-setup/app/ssl"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSSLGenerator(t *testing.T) {
	// when
	generator := ssl.NewSSLGenerator()

	// then
	require.NotNil(t, generator)
}

func Test_sslGenerator_GenerateSelfSignedCert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		generator := ssl.NewSSLGenerator()

		// when
		cert, key, err := generator.GenerateSelfSignedCert("fqdn", "myces", 365)

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, cert)
		assert.NotEmpty(t, key)

		certs := validation.SplitPemCertificates(cert)
		assert.Equal(t, 2, len(certs))

		err = validateCert(certs[0])
		require.NoError(t, err)
		err = validateCert(certs[1])
		require.NoError(t, err)
		_, err = validatePEM(key)
		require.NoError(t, err)
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
