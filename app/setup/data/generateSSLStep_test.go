package data_test

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		config := &context.SetupConfiguration{Naming: context.Naming{CertificateType: "selfsigned", Fqdn: "192.168.56.2", Domain: "myces"}}
		step := data.NewGenerateSSLStep(config)
		generatorMock := &mocks.SSLGenerator{}
		generatorMock.On("GenerateSelfSignedCert", "192.168.56.2", "myces", data.CertExpireDays).Return("cert", "key", nil)
		step.SslGenerator = generatorMock

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		assert.Equal(t, "cert", config.Naming.Certificate)
		assert.NotEmpty(t, "key", config.Naming.CertificateKey)
		mock.AssertExpectationsForObjects(t, generatorMock)
	})

	t.Run("failed to generate certificate", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{CertificateType: "selfsigned", Fqdn: "192.168.56.2", Domain: "myces"}}
		step := data.NewGenerateSSLStep(config)
		generatorMock := &mocks.SSLGenerator{}
		generatorMock.On("GenerateSelfSignedCert", "192.168.56.2", "myces", data.CertExpireDays).Return("cert", "key", assert.AnError)
		step.SslGenerator = generatorMock

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		mock.AssertExpectationsForObjects(t, generatorMock)
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
