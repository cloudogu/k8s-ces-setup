package config_test

import (
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := config.ReadConfig("data/testConfig.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "ecosystem", c.Namespace)
	assert.Equal(t, "0.0.0", c.DoguOperatorVersion)
	assert.Equal(t, "0.0.0", c.EtcdServerVersion)
	assert.Equal(t, logrus.DebugLevel, c.LogLevel)
}

func TestReadConfig_doesNotExist(t *testing.T) {
	_, err := config.ReadConfig("data/doesnotexist.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find configuration")
}

func TestReadConfig_notYaml(t *testing.T) {
	_, err := config.ReadConfig("data/invalidConfig.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal errors")
}
