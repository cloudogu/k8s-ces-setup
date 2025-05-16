package data

import (
	"context"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const ecosystemNamespace = "ecosystem"

func TestNewWriteEcosystemCertificateDataStep(t *testing.T) {
	t.Run("successfully create new naming data step", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		secretClient := NewMockSecretClient(t)

		// when
		myStep := NewWriteEcosystemCertificateDataStep(secretClient, testConfig)

		// then
		assert.NotNil(t, myStep)
	})
}

func Test_writeEcosystemCertificateDataStep_GetStepDescription(t *testing.T) {
	t.Run("successfully get naming data step description", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		secretClient := NewMockSecretClient(t)

		// when
		myStep := NewWriteEcosystemCertificateDataStep(secretClient, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write ecosystem certificate data to a secret", description)
	})
}

func Test_writeEcosystemCertificateDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	t.Run("fail to create secret for certificate", func(t *testing.T) {
		certificateSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: certificateSecretName,
			},
			Data: map[string][]byte{
				v1.TLSCertKey:       []byte(""),
				v1.TLSPrivateKeyKey: []byte(""),
			},
		}
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		secretClient := NewMockSecretClient(t)
		secretClient.EXPECT().Create(testCtx, certificateSecret, metav1.CreateOptions{}).Return(nil, assert.AnError)

		myStep := NewWriteEcosystemCertificateDataStep(secretClient, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("succeed to create secret for certificate", func(t *testing.T) {
		certificateSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: certificateSecretName,
			},
			Data: map[string][]byte{
				v1.TLSCertKey:       []byte(""),
				v1.TLSPrivateKeyKey: []byte(""),
			},
		}
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		secretClient := NewMockSecretClient(t)
		secretClient.EXPECT().Create(testCtx, certificateSecret, metav1.CreateOptions{}).Return(nil, nil)

		myStep := NewWriteEcosystemCertificateDataStep(secretClient, testConfig)

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
