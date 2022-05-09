package data_test

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/validation"

	mocks2 "github.com/cloudogu/cesapp-lib/registry/mocks"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/cloudogu/k8s-ces-setup/app/setup/data/mocks"
)

func TestNewWriteLdapDataStep(t *testing.T) {
	t.Parallel()

	t.Run("successfully create new dogu data step", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistry := &mocks2.Registry{}
		testConfig := &context.SetupConfiguration{}

		// when
		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// then
		assert.NotNil(t, myStep)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeLdapDataStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("successfully get dogu data step description", func(t *testing.T) {
		// given
		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistry := &mocks2.Registry{}
		testConfig := &context.SetupConfiguration{}
		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		description := myStep.GetStepDescription()

		// then
		assert.Equal(t, "Write ldap data to the registry", description)
		mock.AssertExpectationsForObjects(t, mockRegistryWriter)
	})
}

func Test_writeLdapDataStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("failed to check if ldap-mapper is enabled", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(false, assert.AnError)

		mockRegistryWriter := &mocks.RegistryWriter{}
		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})

	t.Run("failed to write to the registry", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{}
		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(false, nil)

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", mock.Anything).Return(assert.AnError)

		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.ErrorIs(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})

	ldapConfiguration := context.UserBackend{
		DsType:                    "external",
		Server:                    "myServer",
		AttributeID:               "myAttributeID",
		AttributeGivenName:        "myAttributeGivenName",
		AttributeSurname:          "myAttributeSurname",
		AttributeFullname:         "myAttributeFullName",
		AttributeMail:             "myAttributeMail",
		AttributeGroup:            "myAttributeGroup",
		BaseDN:                    "myBaseDN",
		SearchFilter:              "mySearchFilter",
		ConnectionDN:              "myConnectionDN",
		Password:                  "myPassword",
		Host:                      "myHost",
		Port:                      "myPort",
		LoginID:                   "myLoginID",
		LoginPassword:             "myLoginPassword",
		Encryption:                "myEncryption",
		GroupBaseDN:               "myGroupBaseDN",
		GroupSearchFilter:         "myGroupSearchFilter",
		GroupAttributeName:        "myGroupAttributeName",
		GroupAttributeDescription: "myGroupAttributeDescription",
		GroupAttributeMember:      "myGroupAttributeMember",
	}

	t.Run("successfully write all dogu data to the registry with: embedded ldap, no encryption and no ldap-mapper enabled", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{UserBackend: ldapConfiguration}
		testConfig.UserBackend.DsType = validation.DsTypeEmbedded
		testConfig.UserBackend.Encryption = ""

		registryConfig := context.CustomKeyValue{
			"cas": map[string]interface{}{
				"ldap": map[string]string{
					"attribute_fullname":   "myAttributeFullName",
					"attribute_given_name": "myAttributeGivenName",
					"attribute_group":      "myAttributeGroup",
					"attribute_id":         "myAttributeID",
					"attribute_mail":       "myAttributeMail",
					"base_dn":              "myBaseDN",
					"connection_dn":        "myConnectionDN",
					"ds_type":              "embedded",
					"group_attribute_name": "myGroupAttributeName",
					"group_base_dn":        "myGroupBaseDN",
					"group_search_filter":  "myGroupSearchFilter",
					"host":                 "myHost",
					"port":                 "myPort",
					"search_filter":        "mySearchFilter",
					"server":               "myServer"}},
		}

		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(false, nil)

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})

	t.Run("successfully write all dogu data to the registry with: embedded ldap, with encryption and no ldap-mapper enabled", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{UserBackend: ldapConfiguration}
		testConfig.UserBackend.DsType = validation.DsTypeEmbedded
		testConfig.UserBackend.Encryption = "myEncryption"

		registryConfig := context.CustomKeyValue{
			"cas": map[string]interface{}{
				"ldap": map[string]string{
					"encryption":           "myEncryption",
					"attribute_fullname":   "myAttributeFullName",
					"attribute_given_name": "myAttributeGivenName",
					"attribute_group":      "myAttributeGroup",
					"attribute_id":         "myAttributeID",
					"attribute_mail":       "myAttributeMail",
					"base_dn":              "myBaseDN",
					"connection_dn":        "myConnectionDN",
					"ds_type":              "embedded",
					"group_attribute_name": "myGroupAttributeName",
					"group_base_dn":        "myGroupBaseDN",
					"group_search_filter":  "myGroupSearchFilter",
					"host":                 "myHost",
					"port":                 "myPort",
					"search_filter":        "mySearchFilter",
					"server":               "myServer"}},
		}

		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(false, nil)

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})

	t.Run("successfully write all dogu data to the registry with: embedded ldap, with encryption and ldap-mapper enabled", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{UserBackend: ldapConfiguration}
		testConfig.UserBackend.DsType = validation.DsTypeExternal
		testConfig.UserBackend.Encryption = "myEncryption"

		registryConfig := context.CustomKeyValue{
			"cas": map[string]interface{}{
				"ldap": map[string]string{
					"attribute_fullname":   "myAttributeFullName",
					"attribute_given_name": "myAttributeGivenName",
					"attribute_group":      "myAttributeGroup",
					"attribute_id":         "myAttributeID",
					"attribute_mail":       "myAttributeMail",
					"base_dn":              "myBaseDN",
					"connection_dn":        "myConnectionDN",
					"ds_type":              "external",
					"encryption":           "myEncryption",
					"group_attribute_name": "myGroupAttributeName",
					"group_base_dn":        "myGroupBaseDN",
					"group_search_filter":  "myGroupSearchFilter",
					"host":                 "myHost",
					"port":                 "myPort",
					"search_filter":        "mySearchFilter",
					"server":               "myServer",
				},
			},
			"ldap-mapper": map[string]interface{}{
				"backend": map[string]string{
					"host": "myHost",
					"port": "myPort",
					"type": "external",
				},
				"mapping": map[string]interface{}{
					"group": map[string]string{
						"base_dn":       "myGroupBaseDN",
						"description":   "myGroupAttributeDescription",
						"encryption":    "myEncryption",
						"member":        "myGroupAttributeMember",
						"name":          "myGroupAttributeName",
						"search_filter": "myGroupSearchFilter",
						"server":        "myServer",
					},
					"user": map[string]string{
						"base_dn":       "myBaseDN",
						"full_name":     "myAttributeFullName",
						"given_name":    "myAttributeGivenName",
						"group":         "myAttributeGroup",
						"id":            "myAttributeID",
						"mail":          "myAttributeMail",
						"search_filter": "mySearchFilter",
						"surname":       "myAttributeSurname",
					},
				},
			},
		}

		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(true, nil)

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})

	t.Run("successfully write all dogu data to the registry with: external ldap, with encryption and ldap-mapper enabled", func(t *testing.T) {
		// given
		testConfig := &context.SetupConfiguration{UserBackend: ldapConfiguration}
		testConfig.UserBackend.DsType = validation.DsTypeEmbedded
		testConfig.UserBackend.Encryption = "myEncryption"

		registryConfig := context.CustomKeyValue{"cas": map[string]interface{}{
			"ldap": map[string]string{
				"attribute_fullname":   "myAttributeFullName",
				"attribute_given_name": "myAttributeGivenName",
				"attribute_group":      "myAttributeGroup",
				"attribute_id":         "myAttributeID",
				"attribute_mail":       "myAttributeMail",
				"base_dn":              "myBaseDN",
				"connection_dn":        "myConnectionDN",
				"ds_type":              "embedded",
				"encryption":           "myEncryption",
				"group_attribute_name": "myGroupAttributeName",
				"group_base_dn":        "myGroupBaseDN",
				"group_search_filter":  "myGroupSearchFilter",
				"host":                 "myHost",
				"port":                 "myPort",
				"search_filter":        "mySearchFilter",
				"server":               "myServer"},
		},
			"ldap-mapper": map[string]interface{}{
				"backend": map[string]string{
					"host": "myHost",
					"port": "myPort",
					"type": "embedded",
				},
			},
		}

		mockRegistry := &mocks2.Registry{}
		mockDoguRegistry := &mocks2.DoguRegistry{}
		mockRegistry.On("DoguRegistry").Return(mockDoguRegistry)
		mockDoguRegistry.On("IsEnabled", "ldap-mapper").Return(true, nil)

		mockRegistryWriter := &mocks.RegistryWriter{}
		mockRegistryWriter.On("WriteConfigToRegistry", registryConfig).Return(nil)

		myStep := data.NewWriteLdapDataStep(mockRegistry, mockRegistryWriter, testConfig)

		// when
		err := myStep.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry, mockDoguRegistry, mockRegistryWriter)
	})
}
