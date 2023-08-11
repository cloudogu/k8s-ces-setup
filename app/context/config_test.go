package context

import (
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		// when
		c, err := ReadConfigFromFile("testdata/testConfig.yaml")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "ecosystem", c.TargetNamespace)
		assert.Equal(t, "https://dop.yaml", c.DoguOperatorURL)
		assert.Equal(t, "https://sd.yaml", c.ServiceDiscoveryURL)
		assert.Equal(t, "https://etcds.yaml", c.EtcdServerResourceURL)
		assert.Equal(t, "https://etcdc.yaml", c.EtcdClientImageRepo)
		assert.Equal(t, logrus.DebugLevel, c.LogLevel)
	})

	t.Run("fail on non existent config", func(t *testing.T) {
		// when
		_, err := ReadConfigFromFile("testdata/doesnotexist.yaml")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find configuration")
	})

	t.Run("fail on invalid file content", func(t *testing.T) {
		// when
		_, err := ReadConfigFromFile("testdata/invalidConfig.yaml")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal errors")
	})
}

func TestReadConfigFromCluster(t *testing.T) {
	const testNamespace = "test-namespace"
	t.Run("should return marshalled config", func(t *testing.T) {
		// given
		myFileMap := map[string]string{"k8s-ces-setup.yaml": `dogu_operator_url: https://url.com`}
		mockedConfig := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SetupConfigConfigmap,
				Namespace: testNamespace,
			},
			Data: myFileMap,
		}
		client := fake.NewSimpleClientset(mockedConfig)

		// when
		actual, err := ReadConfigFromCluster(client, testNamespace)

		// then
		require.NoError(t, err)
		extected := &Config{DoguOperatorURL: "https://url.com"}
		assert.Equal(t, extected, actual)
	})
	t.Run("should fail during marshalling config", func(t *testing.T) {
		// given
		myFileMap := map[string]string{"k8s-ces-setup.yaml": `"dogu_operator_url: https://url.com`}
		mockedConfig := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SetupConfigConfigmap,
				Namespace: testNamespace,
			},
			Data: myFileMap,
		}
		client := fake.NewSimpleClientset(mockedConfig)

		// when
		_, err := ReadConfigFromCluster(client, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal configuration from configmap")
	})
}
