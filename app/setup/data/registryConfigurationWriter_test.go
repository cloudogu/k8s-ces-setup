package data_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewGenericConfigurationWriter(t *testing.T) {
	t.Run("create new generic configuration Writer", func(t *testing.T) {
		// given
		registryMock := &mocks.Registry{}

		// when
		writer := data.NewRegistryConfigurationWriter(registryMock)

		// then
		require.NotNil(t, writer)
		mock.AssertExpectationsForObjects(t, registryMock)
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
		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)

		registryMock := &mocks.Registry{}
		registryMock.On("GlobalConfig").Return(globalRegistryMock)

		writer := data.NewRegistryConfigurationWriter(registryMock)

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, registryMock, globalRegistryMock)
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
		globalRegistryMock := &mocks.ConfigurationContext{}
		globalRegistryMock.On("Set", "test/t1", "myTestt1").Return(nil)
		globalRegistryMock.On("Set", "test/t2", "myTestt2").Return(nil)
		globalRegistryMock.On("Set", "test3", "myTestKey3").Return(nil)

		casRegistryMock := &mocks.ConfigurationContext{}
		casRegistryMock.On("Set", "ldap/attribute_fullname", "myAttributeFullName").Return(nil)

		ldapMapperRegistryMock := &mocks.ConfigurationContext{}
		ldapMapperRegistryMock.On("Set", "backend/type", "external").Return(nil)
		ldapMapperRegistryMock.On("Set", "mapping/group/base_dn", "myGroupBaseDN").Return(nil)
		ldapMapperRegistryMock.On("Set", "mapping/group/member", "myGroupAttributeMember").Return(nil)
		ldapMapperRegistryMock.On("Set", "mapping/user/base_dn", "myBaseDN").Return(nil)
		ldapMapperRegistryMock.On("Set", "mapping/user/full_name", "myAttributeFullName").Return(nil)
		ldapMapperRegistryMock.On("Set", "mapping/user/surname", "myAttributeSurname").Return(nil)

		registryMock := &mocks.Registry{}
		registryMock.On("GlobalConfig").Return(globalRegistryMock)
		registryMock.On("DoguConfig", "cas").Return(casRegistryMock)
		registryMock.On("DoguConfig", "ldap-mapper").Return(ldapMapperRegistryMock)

		writer := data.NewRegistryConfigurationWriter(registryMock)

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, registryMock, globalRegistryMock, casRegistryMock, ldapMapperRegistryMock)
	})
}
