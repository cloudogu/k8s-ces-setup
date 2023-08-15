package context

import (
	_ "embed"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
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
		builder.DevHelmRepositoryDataPath = "testdata/testHelmRepoData.yaml"
		fakeClient := fake.NewSimpleClientset()

		// when
		actual, err := builder.NewSetupContext(testCtx, fakeClient)

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", actual.AppVersion)
		assert.Equal(t, "myTestNamespace", actual.AppConfig.TargetNamespace)
		assert.Equal(t, "k8s/k8s-component-operator:0.0.2", actual.AppConfig.ComponentOperatorChart)
		assert.Equal(t, "https://etcdc.yaml", actual.AppConfig.EtcdClientImageRepo)
		assert.Equal(t, "pkcs1v15", actual.AppConfig.KeyProvider)
		assert.Equal(t, "user", actual.DoguRegistryConfiguration.Username)
		assert.Equal(t, "pw", actual.DoguRegistryConfiguration.Password)
		assert.Equal(t, "endpoint", actual.DoguRegistryConfiguration.Endpoint)
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

		helmConfigmap := &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name:      "component-operator-helm-repository",
			Namespace: "myTestNamespace",
		}, Data: map[string]string{"endpoint": "http://helm.repo"}}

		fakeClient := fake.NewSimpleClientset(startupConfigmap, configConfigmap, registrySecret, helmConfigmap)

		// when
		actual, err := builder.NewSetupContext(testCtx, fakeClient)

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
		_, err := builder.NewSetupContext(testCtx, fakeClient)

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
		_, err := builder.NewSetupContext(testCtx, fakeClient)

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
		_, err := builder.NewSetupContext(testCtx, fakeClient)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dogu registry secret k8s-dogu-operator-dogu-registry not found")
	})

	t.Run("helm repo config not found", func(t *testing.T) {
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
		registrySecretData := map[string]string{"endpoint": "endpoint", "username": "username", "password": "password"}
		registrySecret := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-dogu-operator-dogu-registry",
			Namespace: "myTestNamespace"},
			StringData: registrySecretData,
		}

		fakeClient := fake.NewSimpleClientset(configConfigmap, startupConfigmap, registrySecret)

		// when
		_, err := builder.NewSetupContext(testCtx, fakeClient)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "helm repository configMap component-operator-helm-repository not found")
	})

	t.Run("cannot not read current namespace", func(t *testing.T) {
		// given
		originalNamespace, _ := GetEnvVar(EnvironmentVariableTargetNamespace)
		defer func() { _ = os.Setenv(EnvironmentVariableTargetNamespace, originalNamespace) }()
		err := os.Unsetenv(EnvironmentVariableTargetNamespace)
		require.NoError(t, err)

		builder := NewSetupContextBuilder("1.2.3")

		// when
		_, err = builder.NewSetupContext(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not read current namespace: POD_NAMESPACE must be set")
	})

	t.Run("cannot not find setup config file", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableStage, StageDevelopment)

		builder := NewSetupContextBuilder("1.2.3")
		builder.DevSetupConfigPath = "invalid"

		// when
		_, err := builder.NewSetupContext(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find configuration at invalid")
	})

	t.Run("cannot unmarshal setup json", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableStage, StageDevelopment)

		builder := NewSetupContextBuilder("1.2.3")
		builder.DevSetupConfigPath = "testdata/testConfig.yaml"
		builder.DevStartupConfigPath = "testdata/invalidSetupJson.json"

		// when
		_, err := builder.NewSetupContext(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal setup configuration testdata/invalidSetupJson.json")
	})

	t.Run("cannot not find dogu registry secret", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableStage, StageDevelopment)

		builder := NewSetupContextBuilder("1.2.3")
		builder.DevSetupConfigPath = "testdata/testConfig.yaml"
		builder.DevStartupConfigPath = "testdata/testSetupJson.json"
		builder.DevDoguRegistrySecretPath = "invalid.yaml"

		// when
		_, err := builder.NewSetupContext(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find registry secret at invalid.yaml")
	})

	t.Run("cannot not find helm repo config", func(t *testing.T) {
		// given
		t.Setenv(EnvironmentVariableStage, StageDevelopment)

		builder := NewSetupContextBuilder("1.2.3")
		builder.DevSetupConfigPath = "testdata/testConfig.yaml"
		builder.DevStartupConfigPath = "testdata/testSetupJson.json"
		builder.DevDoguRegistrySecretPath = "testdata/testRegistrySecret.yaml"
		builder.DevHelmRepositoryDataPath = "invalid.yaml"

		// when
		_, err := builder.NewSetupContext(testCtx, nil)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find configuration at invalid.yaml")
	})
}

func TestGetSetupStateConfigMap(t *testing.T) {
	t.Run("should create config map", func(t *testing.T) {
		// given
		namespace := "test-namespace"
		clientset := fake.NewSimpleClientset()

		// when
		actual, err := GetSetupStateConfigMap(testCtx, clientset, namespace)

		// then
		require.NoError(t, err)
		assert.Equal(t, "k8s-setup-config", actual.Name)
		assert.Equal(t, "test-namespace", actual.Namespace)
		assert.Contains(t, actual.Labels, "app")
		assert.Contains(t, actual.Labels, "app.kubernetes.io/name")
	})
	t.Run("should fail when getting config map", func(t *testing.T) {
		// given
		namespace := "test-namespace"
		clientset := fake.NewSimpleClientset()
		clientset.CoreV1().(*fakecorev1.FakeCoreV1).PrependReactor("get", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, &v1.ConfigMap{}, assert.AnError
		})

		// when
		_, err := GetSetupStateConfigMap(testCtx, clientset, namespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get configmap [k8s-ces-setup-config]")
	})
	t.Run("should fail to create configmap", func(t *testing.T) {
		// given
		namespace := "test-namespace"
		clientset := fake.NewSimpleClientset()
		clientset.CoreV1().(*fakecorev1.FakeCoreV1).PrependReactor("create", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, &v1.ConfigMap{}, assert.AnError
		})

		// when
		_, err := GetSetupStateConfigMap(testCtx, clientset, namespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create configmap [k8s-setup-config]")
	})
}
