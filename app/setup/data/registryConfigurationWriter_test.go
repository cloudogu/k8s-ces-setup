package data_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/stretchr/testify/require"
)

func TestNewGenericConfigurationWriter(t *testing.T) {
	t.Run("create new generic configuration Writer", func(t *testing.T) {
		// given
		globalReg := &data.MockConfigurationRegistry{}
		doguReg := &data.MockConfigurationRegistry{}

		// when
		writer := data.NewRegistryConfigurationWriter(
			globalReg,
			func(ctx context.Context, name string) (data.ConfigurationRegistry, error) {
				return doguReg, nil
			},
		)

		// then
		require.NotNil(t, writer)
	})
}

func TestGenericConfigurationWriter_WriteConfigToRegistry(t *testing.T) {
	t.Run("failed to write to config", func(t *testing.T) {
		// given
		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"test3": "myTestKey3",
			},
		}

		globalReg := &data.MockConfigurationRegistry{}
		doguReg := &data.MockConfigurationRegistry{}
		globalReg.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

		writer := data.NewRegistryConfigurationWriter(globalReg, func(ctx context.Context, name string) (data.ConfigurationRegistry, error) {
			return doguReg, nil
		})

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("set all keys correctly", func(t *testing.T) {
		// given
		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"test": map[string]string{
					"t1": "myTestt1",
					"t2": "myTestt2",
				},
				"test3": "myTestKey3",
			},
			"cas": map[string]interface{}{
				"ldap": map[string]string{
					"attribute_fullname": "myAttributeFullName",
				},
			},
			"ldap-mapper": map[string]interface{}{
				"backend": map[string]string{
					"type": "external",
				},
				"mapping": map[string]interface{}{
					"group": map[string]string{
						"base_dn": "myGroupBaseDN",
						"member":  "myGroupAttributeMember",
					},
					"user": map[string]string{
						"base_dn":   "myBaseDN",
						"full_name": "myAttributeFullName",
						"surname":   "myAttributeSurname",
					},
				},
			},
		}

		globalReg := &data.MockConfigurationRegistry{}
		globalReg.EXPECT().Set(mock.Anything, "test/t1", "myTestt1").Return(nil)
		globalReg.EXPECT().Set(mock.Anything, "test/t2", "myTestt2").Return(nil)
		globalReg.EXPECT().Set(mock.Anything, "test3", "myTestKey3").Return(nil)

		doguReg := &data.MockConfigurationRegistry{}
		doguReg.EXPECT().Set(mock.Anything, "ldap/attribute_fullname", "myAttributeFullName").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "backend/type", "external").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "mapping/group/base_dn", "myGroupBaseDN").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "mapping/group/member", "myGroupAttributeMember").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "mapping/user/base_dn", "myBaseDN").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "mapping/user/full_name", "myAttributeFullName").Return(nil)
		doguReg.EXPECT().Set(mock.Anything, "mapping/user/surname", "myAttributeSurname").Return(nil)

		writer := data.NewRegistryConfigurationWriter(globalReg, func(ctx context.Context, name string) (data.ConfigurationRegistry, error) {
			return doguReg, nil
		})

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.NoError(t, err)
	})
}
