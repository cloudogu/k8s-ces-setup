package context_test

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := context.ReadConfig("testdata/testConfig.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "ecosystem", c.TargetNamespace)
	assert.Equal(t, "https://dop.yaml", c.DoguOperatorURL)
	assert.Equal(t, "https://etcds.yaml", c.EtcdServerResourceURL)
	assert.Equal(t, "https://etcdc.yaml", c.EtcdClientImageRepo)
	assert.Equal(t, logrus.DebugLevel, c.LogLevel)
}

func TestReadConfig_doesNotExist(t *testing.T) {
	_, err := context.ReadConfig("testdata/doesnotexist.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find configuration")
}

func TestReadConfig_notYaml(t *testing.T) {
	_, err := context.ReadConfig("testdata/invalidConfig.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal errors")
}
