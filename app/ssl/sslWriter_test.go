package ssl_test

import (
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/ssl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSSLWriter(t *testing.T) {

}

func Test_sslWriter_WriteCertificate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", "self-signed").Return(nil)
		globalConfig.On("Set", "certificate/server.crt", "cert").Return(nil)
		globalConfig.On("Set", "certificate/server.key", "key").Return(nil)
		writer := ssl.NewSSLWriter(globalConfig)

		// when
		err := writer.WriteCertificate("self-signed", "cert", "key")

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write type", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", "self-signed").Return(assert.AnError)
		writer := ssl.NewSSLWriter(globalConfig)

		// when
		err := writer.WriteCertificate("self-signed", "cert", "key")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate type")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write certificate", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", "self-signed").Return(nil)
		globalConfig.On("Set", "certificate/server.crt", "cert").Return(assert.AnError)
		writer := ssl.NewSSLWriter(globalConfig)

		// when
		err := writer.WriteCertificate("self-signed", "cert", "key")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})

	t.Run("failed to write certificate key", func(t *testing.T) {
		// given
		globalConfig := &mocks.ConfigurationContext{}
		globalConfig.On("Set", "certificate/type", "self-signed").Return(nil)
		globalConfig.On("Set", "certificate/server.crt", "cert").Return(nil)
		globalConfig.On("Set", "certificate/server.key", "key").Return(assert.AnError)
		writer := ssl.NewSSLWriter(globalConfig)

		// when
		err := writer.WriteCertificate("self-signed", "cert", "key")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set certificate key")
		mock.AssertExpectationsForObjects(t, globalConfig)
	})
}
