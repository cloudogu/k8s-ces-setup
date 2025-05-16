package data_test

import (
	"context"
	"testing"

	appcontext "github.com/cloudogu/k8s-ces-setup/v2/app/context"
	"github.com/cloudogu/k8s-ces-setup/v2/app/setup/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewWriteRegistryConfigEncryptedStep(t *testing.T) {
	// given
	setupConfig := &appcontext.SetupJsonConfiguration{}
	fakeClient := fake.NewSimpleClientset()

	// when
	step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")

	// then
	require.NotNil(t, step)
}

func Test_writeRegistryConfigEncryptedStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get registry config data step description", func(t *testing.T) {
		// given
		myStep := data.NewWriteRegistryConfigEncryptedStep(nil, fake.NewSimpleClientset(), "test")

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write registry config encrypted data to the registry", description)
	})
}

func Test_writeRegistryConfigEncryptedStep_PerformSetupStep(t *testing.T) {
	var testCtx = context.Background()
	t.Run("success embedded", func(t *testing.T) {
		// given
		admin := appcontext.User{Password: "adminPw"}
		embeddedUserBackend := appcontext.UserBackend{DsType: "embedded"}
		setupConfig := &appcontext.SetupJsonConfiguration{UserBackend: embeddedUserBackend, Admin: admin}
		fakeClient := fake.NewSimpleClientset()
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")

		// when
		err := step.PerformSetupStep(testCtx)
		require.NoError(t, err)

		// then
		secret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-config", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "cas-config", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-mapper-config", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, secret)
		assert.NotNil(t, secret.StringData)
		assert.Equal(t, 1, len(secret.StringData))
		assert.Equal(t, "admin_password: adminPw\n", secret.StringData["config.yaml"])
	})

	t.Run("should override embedded admin pw and append user defined config", func(t *testing.T) {
		// given
		admin := appcontext.User{Password: "adminPw"}
		embeddedUserBackend := appcontext.UserBackend{DsType: "embedded"}
		setupConfig := &appcontext.SetupJsonConfiguration{
			UserBackend: embeddedUserBackend,
			Admin:       admin,
			RegistryConfigEncrypted: map[string]map[string]any{
				"ldap": {"admin_password": "overrideThis", "fromUser": "user"},
			},
		}
		fakeClient := fake.NewSimpleClientset()
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")

		// when
		err := step.PerformSetupStep(testCtx)
		require.NoError(t, err)

		// then
		secret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-config", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "cas-config", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-mapper-config", metav1.GetOptions{})
		require.NoError(t, err)

		require.NotNil(t, secret)
		assert.NotNil(t, secret.StringData)
		assert.Equal(t, 1, len(secret.StringData))
		assert.Equal(t, "admin_password: adminPw\nfromUser: user\n", secret.StringData["config.yaml"])
	})

	t.Run("success external", func(t *testing.T) {
		// given
		admin := appcontext.User{Password: "adminPw"}
		embeddedUserBackend := appcontext.UserBackend{DsType: "external", ConnectionDN: "connection", Password: "ldapPw"}
		dogus := appcontext.Dogus{Install: []string{"ldap-mapper"}}
		setupConfig := &appcontext.SetupJsonConfiguration{
			UserBackend: embeddedUserBackend,
			Admin:       admin,
			Dogus:       dogus,
			RegistryConfigEncrypted: map[string]map[string]any{
				"ldap-mapper": {"backend.password": "overrideThis", "backend.connection_dn": "overrideThis", "fromUser": "user"},
				"cas":         {"password": "overrideThis", "fromUser": "user"},
			},
		}
		fakeClient := fake.NewSimpleClientset()
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")

		// when
		err := step.PerformSetupStep(testCtx)
		require.NoError(t, err)

		// then
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-config", metav1.GetOptions{})
		require.NoError(t, err)
		casSecret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "cas-config", metav1.GetOptions{})
		require.NoError(t, err)
		ldapMapperSecret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-mapper-config", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, casSecret)
		assert.NotNil(t, casSecret.StringData)
		assert.Equal(t, 1, len(casSecret.StringData))
		assert.Equal(t, "fromUser: user\npassword: ldapPw\n", casSecret.StringData["config.yaml"])
		require.NotNil(t, ldapMapperSecret)
		assert.NotNil(t, ldapMapperSecret.StringData)
		assert.Equal(t, 1, len(ldapMapperSecret.StringData))
		assert.Equal(t, "backend:\n    connection_dn: connection\n    password: ldapPw\nbackend.connection_dn: overrideThis\nbackend.password: overrideThis\nfromUser: user\n", ldapMapperSecret.StringData["config.yaml"])
	})

	t.Run("should override external backend pw, connection_dn and append user defined config", func(t *testing.T) {
		// given
		admin := appcontext.User{Password: "adminPw"}
		embeddedUserBackend := appcontext.UserBackend{DsType: "external", ConnectionDN: "connection", Password: "ldapPw"}
		dogus := appcontext.Dogus{Install: []string{"ldap-mapper"}}
		setupConfig := &appcontext.SetupJsonConfiguration{
			UserBackend:             embeddedUserBackend,
			Admin:                   admin,
			Dogus:                   dogus,
			RegistryConfigEncrypted: map[string]map[string]any{},
		}
		fakeClient := fake.NewSimpleClientset()
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")

		// when
		err := step.PerformSetupStep(testCtx)
		require.NoError(t, err)

		// then
		_, err = fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-config", metav1.GetOptions{})
		require.NoError(t, err)
		casSecret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "cas-config", metav1.GetOptions{})
		require.NoError(t, err)
		ldapMapperSecret, err := fakeClient.CoreV1().Secrets("test").Get(testCtx, "ldap-mapper-config", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, casSecret)
		assert.NotNil(t, casSecret.StringData)
		assert.Equal(t, 1, len(casSecret.StringData))
		assert.Equal(t, "password: ldapPw\n", casSecret.StringData["config.yaml"])
		require.NotNil(t, ldapMapperSecret)
		assert.NotNil(t, ldapMapperSecret.StringData)
		assert.Equal(t, 1, len(ldapMapperSecret.StringData))
		assert.Equal(t, "backend:\n    connection_dn: connection\n    password: ldapPw\n", ldapMapperSecret.StringData["config.yaml"])
	})
}
