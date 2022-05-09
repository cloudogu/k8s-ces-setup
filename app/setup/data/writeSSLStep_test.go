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
		writerMock := &mocks.SSLWriter{}
		writerMock.On("WriteCertificate", naming.CertificateType, naming.Certificate, naming.CertificateKey).Return(nil)
		step := &data.WriteSSLStep{}
		step.SSLWriter = writerMock
		step.Config = config

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, writerMock)
	})

	t.Run("failed to write certificate", func(t *testing.T) {
		// given
		writerMock := &mocks.SSLWriter{}
		writerMock.On("WriteCertificate", naming.CertificateType, naming.Certificate, naming.CertificateKey).Return(assert.AnError)
		step := &data.WriteSSLStep{}
		step.SSLWriter = writerMock
		step.Config = config

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		mock.AssertExpectationsForObjects(t, writerMock)
	})
}
