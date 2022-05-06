package data_test

import (
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewWriteSSL(t *testing.T) {
	// when
	step := data.NewWriteSSLStep(nil, nil)

	// then
	require.NotNil(t, step)
}

func Test_writeSSLStep_GetStepDescription(t *testing.T) {
	// given
	step := data.NewWriteSSLStep(nil, nil)

	// when
	description := step.GetStepDescription()

	// then
	assert.Contains(t, description, "Write SSL certificate and key")
}

func Test_writeSSLStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	naming := context.Naming{CertificateType: "selfsigned", Certificate: "cert", CertificateKey: "key"}
	config := &context.SetupConfiguration{Naming: naming}
	t.Run("success", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", naming.CertificateType).Return(nil)
		globalConfig.On("Set", "certificate/server.crt", naming.Certificate).Return(nil)
		globalConfig.On("Set", "certificate/server.key", naming.CertificateKey).Return(nil)
		step := data.NewWriteSSLStep(config, globalConfig)

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write type", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", naming.CertificateType).Return(assert.AnError)
		step := data.NewWriteSSLStep(config, globalConfig)

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate type")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write certificate", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", naming.CertificateType).Return(nil)
		globalConfig.On("Set", "certificate/server.crt", naming.Certificate).Return(assert.AnError)
		step := data.NewWriteSSLStep(config, globalConfig)

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write certificate key", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", naming.CertificateType).Return(nil)
		globalConfig.On("Set", "certificate/server.crt", naming.Certificate).Return(nil)
		globalConfig.On("Set", "certificate/server.key", naming.CertificateKey).Return(assert.AnError)
		step := data.NewWriteSSLStep(config, globalConfig)

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate key")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})
}
