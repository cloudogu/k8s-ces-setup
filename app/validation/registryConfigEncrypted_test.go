package validation_test

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRegistryConfigEncryptedValidator(t *testing.T) {
	// when
	validator := validation.NewRegistryConfigEncryptedValidator()

	// then
	require.NotNil(t, validator)
}

func Test_registryConfigEncryptedValidator_ValidateRegistryConfigEncrypted(t *testing.T) {
	registryConfig := context.CustomKeyValue{}
	ldapConfig := map[string]interface{}{
		"key": "value",
	}
	registryConfig["ldap"] = ldapConfig

	t.Run("success", func(t *testing.T) {
		// given
		validator := validation.NewRegistryConfigEncryptedValidator()
		dogus := context.Dogus{Install: []string{"official/ldap", "testing/cas"}}
		config := &context.SetupConfiguration{RegistryConfigEncrypted: registryConfig, Dogus: dogus}

		// when
		err := validator.ValidateRegistryConfigEncrypted(config)

		// then
		require.NoError(t, err)
	})

	t.Run("key is not in dogu install list", func(t *testing.T) {
		// given
		validator := validation.NewRegistryConfigEncryptedValidator()
		dogus := context.Dogus{Install: []string{"testing/cas"}}

		config := &context.SetupConfiguration{RegistryConfigEncrypted: registryConfig, Dogus: dogus}

		// when
		err := validator.ValidateRegistryConfigEncrypted(config)

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "key ldap does not exist in dogu install list")
	})
}
