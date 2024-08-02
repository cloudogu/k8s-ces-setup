package data_test

import (
	gocontext "context"
	"oras.land/oras-go/pkg/context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
)

func TestNewWriteNamingDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new naming data step", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &appcontext.SetupJsonConfiguration{}
		fakeClient := fake.NewSimpleClientset()

		// when
		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// then
		assert.NotNil(t, myStep)
	})
}

func Test_writeNamingDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get naming data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		testConfig := &appcontext.SetupJsonConfiguration{}
		fakeClient := fake.NewSimpleClientset()

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write naming data to the registry", description)
	})
}

func Test_writeNamingDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	t.Run("fail to write anything in the registry", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{}
		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(mock.Anything).Return(assert.AnError)
		fakeClient := fake.NewSimpleClientset()

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("successfully apply all naming entries", func(t *testing.T) {
		// given
		testConfig := &appcontext.SetupJsonConfiguration{Naming: appcontext.Naming{
			Fqdn:            "myFqdn",
			Domain:          "myDomain",
			MailAddress:     "my@mail.address",
			CertificateType: "self-signed",
			Certificate:     "myCertificate",
			CertificateKey:  "myCertificateKey",
			RelayHost:       "myRelayHost",
			UseInternalIp:   true,
			InternalIp:      "1.2.3.4",
		}}

		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"fqdn":                   "myFqdn",
				"domain":                 "myDomain",
				"certificate/type":       "self-signed",
				"certificate/server.crt": "myCertificate",
				"certificate/server.key": "myCertificateKey",
				"mail_address":           "my@mail.address",
				"k8s/use_internal_ip":    "true",
				"k8s/internal_ip":        "1.2.3.4",
			},
			"postfix": map[string]interface{}{
				"relayhost": "myRelayHost",
			},
		}

		mockRegistryWriter := data.NewMockRegistryWriter(t)
		mockRegistryWriter.EXPECT().WriteConfigToRegistry(registryConfig).Return(nil)

		fakeClient := fake.NewSimpleClientset()

		myStep := data.NewWriteNamingDataStep(mockRegistryWriter, testConfig, fakeClient, "ecosystem")

		// when
		err := myStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)

		tlsSecret, err := fakeClient.CoreV1().Secrets("ecosystem").Get(gocontext.Background(), "ecosystem-certificate", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, []byte("myCertificate"), tlsSecret.Data[v1.TLSCertKey])
		assert.Equal(t, []byte("myCertificateKey"), tlsSecret.Data[v1.TLSPrivateKeyKey])
	})

	// TODO add tests for internal ip when addressed
}
