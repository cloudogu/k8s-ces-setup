package context_test

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		// when
		c, err := context.ReadConfig("testdata/testConfig.yaml")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "ecosystem", c.TargetNamespace)
		assert.Equal(t, "https://dop.yaml", c.DoguOperatorURL)
		assert.Equal(t, "https://sd.yaml", c.ServiceDiscoveryURL)
		assert.Equal(t, "https://etcds.yaml", c.EtcdServerResourceURL)
		assert.Equal(t, "https://etcdc.yaml", c.EtcdClientImageRepo)
		assert.Equal(t, logrus.DebugLevel, c.LogLevel)
	})

	t.Run("fail on non existen config", func(t *testing.T) {
		// when
		_, err := context.ReadConfig("testdata/doesnotexist.yaml")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find configuration")
	})

	t.Run("fail on invalid file content", func(t *testing.T) {
		// when
		_, err := context.ReadConfig("testdata/invalidConfig.yaml")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal errors")
	})
}
