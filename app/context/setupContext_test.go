package context

import (
	_ "embed"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/testSetupJson.json
var setupJSONBytes []byte

//go:embed testdata/testConfig.yaml
var configBytes []byte

func TestNewSetupContextBuilder(t *testing.T) {
	t.Run("should return new context", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableTargetNamespace, "myTestNamespace")

		// when
		actual := NewSetupContextBuilder("1.2.3")

		// then
		assert.Equal(t, "k8s/dev-resources/k8s-ces-setup.yaml", actual.DevSetupConfigPath)
		assert.Equal(t, "k8s/dev-resources/setup.json", actual.DevStartupConfigPath)
	})
}

func Test_GetEnvVar(t *testing.T) {
	_ = os.Unsetenv(EnvironmentVariableTargetNamespace)

	t.Run("successfully query env var namespace", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableTargetNamespace, "myTestNamespace")

		// when
		ns, err := GetEnvVar(EnvironmentVariableTargetNamespace)

		// then
		require.NoError(t, err)

		assert.Equal(t, "myTestNamespace", ns)
	})

	t.Run("failed to query env var namespace", func(t *testing.T) {
		// when
		_, err := GetEnvVar(EnvironmentVariableTargetNamespace)

		// then
		require.Error(t, err)
	})
}

func TestSetupContextBuilder_NewSetupContext(t *testing.T) {
	t.Setenv(EnvironmentVariableStage, StageDevelopment)
	t.Setenv(EnvironmentVariableTargetNamespace, "myTestNamespace")
	t.Run("success read dev resources", func(t *testing.T) {
		// given
		builder := NewSetupContextBuilder("1.2.3")
		builder.DevStartupConfigPath = "testdata/testSetupJson.json"
		builder.DevSetupConfigPath = "testdata/testConfig.yaml"
		builder.DevDoguRegistrySecretPath = "testdata/testRegistrySecret.yaml"
		fakeClient := fake.NewSimpleClientset()

		// when
		actual, err := builder.NewSetupContext(fakeClient)

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", actual.AppVersion)
		assert.Equal(t, "myTestNamespace", actual.AppConfig.TargetNamespace)
		assert.Equal(t, "https://dop.yaml", actual.AppConfig.DoguOperatorURL)
		assert.Equal(t, "https://sd.yaml", actual.AppConfig.ServiceDiscoveryURL)
		assert.Equal(t, "https://etcds.yaml", actual.AppConfig.EtcdServerResourceURL)
		assert.Equal(t, "https://etcdc.yaml", actual.AppConfig.EtcdClientImageRepo)
		assert.Equal(t, "pkcs1v15", actual.AppConfig.KeyProvider)
		assert.Equal(t, "user", actual.doguRegistrySecret.Username)
		assert.Equal(t, "pw", actual.doguRegistrySecret.Password)
		assert.Equal(t, "endpoint", actual.doguRegistrySecret.Endpoint)
	})

	_ = os.Unsetenv(EnvironmentVariableStage)

	t.Run("success read cluster resources", func(t *testing.T) {
		// given
		builder := NewSetupContextBuilder("1.2.3")
		startupData := map[string]string{"setup.json": string(setupJSONBytes)}
		startupConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-json",
			Namespace: "myTestNamespace",
		}, Data: startupData}

		configData := map[string]string{"k8s-ces-setup.yaml": string(configBytes)}
		configConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-config",
			Namespace: "myTestNamespace",
		}, Data: configData}

		registrySecretData := map[string]string{"endpoint": "endpoint", "username": "username", "password": "password"}
		registrySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-dogu-operator-dogu-registry",
			Namespace: "myTestNamespace"},
			StringData: registrySecretData,
		}

		fakeClient := fake.NewSimpleClientset(startupConfigmap, configConfigmap, registrySecret)

		// when
		actual, err := builder.NewSetupContext(fakeClient)

		// then
		require.NoError(t, err)
		require.NotNil(t, actual)
	})

	t.Run("config not found", func(t *testing.T) {
		// given
		builder := NewSetupContextBuilder("1.2.3")
		startupData := map[string]string{"setup.json": string(setupJSONBytes)}
		startupConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-json",
			Namespace: "myTestNamespace",
		}, Data: startupData}
		registrySecretData := map[string]string{"endpoint": "endpoint", "username": "username", "password": "password"}
		registrySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-dogu-operator-dogu-registry",
			Namespace: "myTestNamespace"},
			StringData: registrySecretData,
		}

		fakeClient := fake.NewSimpleClientset(startupConfigmap, registrySecret)

		// when
		_, err := builder.NewSetupContext(fakeClient)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get setup configuration from cluster")
	})

	t.Run("setup json not found", func(t *testing.T) {
		// given
		builder := NewSetupContextBuilder("1.2.3")
		configData := map[string]string{"k8s-ces-setup.yaml": string(configBytes)}
		configConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-config",
			Namespace: "myTestNamespace",
		}, Data: configData}
		registrySecretData := map[string]string{"endpoint": "endpoint", "username": "username", "password": "password"}
		registrySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-dogu-operator-dogu-registry",
			Namespace: "myTestNamespace"},
			StringData: registrySecretData,
		}

		fakeClient := fake.NewSimpleClientset(configConfigmap, registrySecret)

		// when
		_, err := builder.NewSetupContext(fakeClient)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get setup.json configmap")
	})

	t.Run("dogu registry secret not found", func(t *testing.T) {
		// given
		builder := NewSetupContextBuilder("1.2.3")
		configData := map[string]string{"k8s-ces-setup.yaml": string(configBytes)}
		configConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-config",
			Namespace: "myTestNamespace",
		}, Data: configData}
		startupData := map[string]string{"setup.json": string(setupJSONBytes)}
		startupConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-ces-setup-json",
			Namespace: "myTestNamespace",
		}, Data: startupData}

		fakeClient := fake.NewSimpleClientset(configConfigmap, startupConfigmap)

		// when
		_, err := builder.NewSetupContext(fakeClient)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dogu registry secret k8s-dogu-operator-dogu-registry not found")
	})
}
