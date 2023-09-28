package context

import (
	"context"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

var testCtx = context.Background()

func TestReadConfig(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		// when
		c, err := ReadConfigFromFile("testdata/testConfig.yaml")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "ecosystem", c.TargetNamespace)
		assert.Equal(t, "k8s/k8s-component-operator-crd:0.0.2", c.ComponentOperatorCrdChart)
		assert.Equal(t, "k8s/k8s-component-operator:0.0.2", c.ComponentOperatorChart)
		assert.Len(t, c.Components, 3)
		assert.Equal(t, "1.2.3", c.Components["k8s/k8s-etcd"])
		assert.Equal(t, "0.0.1", c.Components["k8s/k8s-dogu-operator"])
		assert.Equal(t, "latest", c.Components["k8s/k8s-service-discovery"])
		assert.Equal(t, "https://etcdc.yaml", c.EtcdClientImageRepo)
		assert.Equal(t, logrus.DebugLevel, *c.LogLevel)
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
		myFileMap := map[string]string{"k8s-ces-setup.yaml": "component_operator_crd_chart: https://crd.chart\ncomponent_operator_chart: https://url.com"}
		mockedConfig := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      SetupConfigConfigmap,
				Namespace: testNamespace,
			},
			Data: myFileMap,
		}
		client := fake.NewSimpleClientset(mockedConfig)

		// when
		actual, err := ReadConfigFromCluster(testCtx, client, testNamespace)

		// then
		require.NoError(t, err)
		extected := &Config{ComponentOperatorCrdChart: "https://crd.chart", ComponentOperatorChart: "https://url.com"}
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
		_, err := ReadConfigFromCluster(testCtx, client, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal configuration from configmap")
	})
}
