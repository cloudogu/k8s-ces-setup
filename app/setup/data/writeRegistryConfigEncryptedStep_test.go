package data_test

import (
	gocontext "context"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewWriteRegistryConfigEncryptedStep(t *testing.T) {
	// given
	setupConfig := &context.SetupConfiguration{}
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
		myStep := data.NewWriteRegistryConfigEncryptedStep(nil, nil, "test")

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write registry config encrypted data to the registry", description)
	})
}

func Test_writeRegistryConfigEncryptedStep_PerformSetupStep(t *testing.T) {
	t.Run("success embedded", func(t *testing.T) {
		// given
		admin := context.User{Password: "adminPw"}
		embeddedUserBackend := context.UserBackend{DsType: "embedded"}
		setupConfig := &context.SetupConfiguration{UserBackend: embeddedUserBackend, Admin: admin}
		fakeClient := fake.NewSimpleClientset()
		emptyMap := map[string]map[string]string{}
		writerMock := &mocks.MapWriter{}
		writerMock.On("WriteConfigToStringDataMap", mock.Anything).Return(emptyMap, nil)
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")
		step.Writer = writerMock

		// when
		err := step.PerformSetupStep()
		require.NoError(t, err)

		// then
		secret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "cas-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-mapper-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		require.NotNil(t, secret)
		assert.NotNil(t, secret.StringData)
		assert.Equal(t, 1, len(secret.StringData))
		assert.Equal(t, "adminPw", secret.StringData["admin_password"])
	})

	t.Run("should override embedded admin pw and append user defined config", func(t *testing.T) {
		// given
		admin := context.User{Password: "adminPw"}
		embeddedUserBackend := context.UserBackend{DsType: "embedded"}
		setupConfig := &context.SetupConfiguration{UserBackend: embeddedUserBackend, Admin: admin}
		fakeClient := fake.NewSimpleClientset()
		registryConfigEncrypted := map[string]map[string]string{}
		registryConfigEncrypted["ldap"] = map[string]string{"admin_password": "overrideThis", "fromUser": "user"}
		writerMock := &mocks.MapWriter{}
		writerMock.On("WriteConfigToStringDataMap", mock.Anything).Return(registryConfigEncrypted, nil)
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")
		step.Writer = writerMock

		// when
		err := step.PerformSetupStep()
		require.NoError(t, err)

		// then
		secret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "cas-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-mapper-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		require.NotNil(t, secret)
		assert.NotNil(t, secret.StringData)
		assert.Equal(t, 2, len(secret.StringData))
		assert.Equal(t, "adminPw", secret.StringData["admin_password"])
		assert.Equal(t, "user", secret.StringData["fromUser"])
	})

	t.Run("success external", func(t *testing.T) {
		// given
		admin := context.User{Password: "adminPw"}
		embeddedUserBackend := context.UserBackend{DsType: "external", ConnectionDN: "connection", Password: "ldapPw"}
		dogus := context.Dogus{Install: []string{"ldap-mapper"}}
		setupConfig := &context.SetupConfiguration{UserBackend: embeddedUserBackend, Admin: admin, Dogus: dogus}
		fakeClient := fake.NewSimpleClientset()
		registryConfigEncrypted := map[string]map[string]string{}
		registryConfigEncrypted["ldap-mapper"] = map[string]string{"backend.password": "overrideThis", "backend.connection_dn": "overrideThis", "fromUser": "user"}
		registryConfigEncrypted["cas"] = map[string]string{"password": "overrideThis", "fromUser": "user"}
		writerMock := &mocks.MapWriter{}
		writerMock.On("WriteConfigToStringDataMap", mock.Anything).Return(registryConfigEncrypted, nil)
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")
		step.Writer = writerMock

		// when
		err := step.PerformSetupStep()
		require.NoError(t, err)

		// then
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		casSecret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "cas-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		ldapMapperSecret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-mapper-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, casSecret)
		assert.NotNil(t, casSecret.StringData)
		assert.Equal(t, 2, len(casSecret.StringData))
		assert.Equal(t, "ldapPw", casSecret.StringData["password"])
		assert.Equal(t, "user", casSecret.StringData["fromUser"])
		require.NotNil(t, ldapMapperSecret)
		assert.NotNil(t, ldapMapperSecret.StringData)
		assert.Equal(t, 3, len(ldapMapperSecret.StringData))
		assert.Equal(t, "ldapPw", ldapMapperSecret.StringData["backend.password"])
		assert.Equal(t, "connection", ldapMapperSecret.StringData["backend.connection_dn"])
		assert.Equal(t, "user", ldapMapperSecret.StringData["fromUser"])
	})

	t.Run("should override external backend pw, connection_dn and append user defined config", func(t *testing.T) {
		// given
		admin := context.User{Password: "adminPw"}
		embeddedUserBackend := context.UserBackend{DsType: "external", ConnectionDN: "connection", Password: "ldapPw"}
		dogus := context.Dogus{Install: []string{"ldap-mapper"}}
		setupConfig := &context.SetupConfiguration{UserBackend: embeddedUserBackend, Admin: admin, Dogus: dogus}
		fakeClient := fake.NewSimpleClientset()
		emptyMap := map[string]map[string]string{}
		writerMock := &mocks.MapWriter{}
		writerMock.On("WriteConfigToStringDataMap", mock.Anything).Return(emptyMap, nil)
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, fakeClient, "test")
		step.Writer = writerMock

		// when
		err := step.PerformSetupStep()
		require.NoError(t, err)

		// then
		_, err = fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-secrets", metav1.GetOptions{})
		require.True(t, errors.IsNotFound(err))
		casSecret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "cas-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		ldapMapperSecret, err := fakeClient.CoreV1().Secrets("test").Get(gocontext.Background(), "ldap-mapper-secrets", metav1.GetOptions{})
		require.NoError(t, err)
		require.NotNil(t, casSecret)
		assert.NotNil(t, casSecret.StringData)
		assert.Equal(t, 1, len(casSecret.StringData))
		assert.Equal(t, "ldapPw", casSecret.StringData["password"])
		require.NotNil(t, ldapMapperSecret)
		assert.NotNil(t, ldapMapperSecret.StringData)
		assert.Equal(t, 2, len(ldapMapperSecret.StringData))
		assert.Equal(t, "ldapPw", ldapMapperSecret.StringData["backend.password"])
		assert.Equal(t, "connection", ldapMapperSecret.StringData["backend.connection_dn"])
	})

	t.Run("fail to write config to map", func(t *testing.T) {
		// given
		setupConfig := &context.SetupConfiguration{}
		writerMock := &mocks.MapWriter{}
		writerMock.On("WriteConfigToStringDataMap", mock.Anything).Return(nil, assert.AnError)
		step := data.NewWriteRegistryConfigEncryptedStep(setupConfig, nil, "test")
		step.Writer = writerMock

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write registry config encrypted")
	})
}

func Test_mapConfigurationWriter_WriteConfigToMap(t *testing.T) {
	// given
	ldapMapperRegistryConfig := map[string]interface{}{
		"user": map[string]string{
			"base_dn": "userBase",
		},
		"group": map[string]string{
			"base_dn": "groupBase",
			"test":    "test",
		},
	}
	ldapRegistryConfig := map[string]interface{}{
		"admin_password": "password",
	}

	registryConfig := context.CustomKeyValue{}
	registryConfig["ldap-mapper"] = ldapMapperRegistryConfig
	registryConfig["ldap"] = ldapRegistryConfig
	mapWriter := data.NewStringDataConfigurationWriter()

	// when
	result, err := mapWriter.WriteConfigToStringDataMap(registryConfig)

	// then
	require.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 3, len(result["ldap-mapper"]))
	assert.Equal(t, "userBase", result["ldap-mapper"]["user.base_dn"])
	assert.Equal(t, "groupBase", result["ldap-mapper"]["group.base_dn"])
	assert.Equal(t, "test", result["ldap-mapper"]["group.test"])
	assert.Equal(t, 1, len(result["ldap"]))
	assert.Equal(t, "password", result["ldap"]["admin_password"])
}
