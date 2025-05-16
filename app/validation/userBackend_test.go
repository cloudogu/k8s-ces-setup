package validation

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/v4/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserBackendValidator(t *testing.T) {
	// when
	validator := NewUserBackendValidator()

	// then
	require.NotNil(t, validator)
}

func Test_userBackendValidator_validateEmbeddedBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      context.UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid attributeID backend type embedded", context.UserBackend{DsType: "embedded"}, "invalid attributeID valid options are", assert.Error, assert.Contains},
		{"invalid attributeFullName backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid"}, "invalid attributeFullName valid options are", assert.Error, assert.Contains},
		{"invalid attributeMail backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn"}, "invalid attributeMail valid options are", assert.Error, assert.Contains},
		{"invalid attributeGroup backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail"}, "invalid attributeGroup valid options are", assert.Error, assert.Contains},
		{"invalid searchFilter backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf"}, "invalid searchFilter valid options are", assert.Error, assert.Contains},
		{"invalid host backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)"}, "invalid host valid options are", assert.Error, assert.Contains},
		{"invalid port backend type embedded", context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)", Host: "ldap"}, "invalid port valid options are", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &userBackendValidator{}

			// when
			result := validator.validateEmbeddedBackend(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful embedded backend validation", func(t *testing.T) {
		// given
		userBackend := context.UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)", Host: "ldap", Port: "389"}
		validator := &userBackendValidator{}

		// when
		result := validator.validateEmbeddedBackend(userBackend)

		// then
		require.NoError(t, result)
	})
}

func Test_userBackendValidator_validateExternalBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      context.UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid attributeID", context.UserBackend{Server: "activeDirectory"}, "invalid attributeID valid options are", assert.Error, assert.Contains},
		{"invalid attributeFullName", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName"}, "invalid attributeFullName valid options are", assert.Error, assert.Contains},
		{"invalid attributeMail", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn"}, "invalid attributeMail valid options are", assert.Error, assert.Contains},
		{"invalid attributeGroup", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail"}, "invalid attributeGroup valid options are", assert.Error, assert.Contains},
		{"invalid searchFilter", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf"}, "invalid searchFilter valid options are", assert.Error, assert.Contains},
		{"invalid server", context.UserBackend{DsType: "external"}, "invalid server valid options are", assert.Error, assert.Contains},
		{"attributeGivenName not set", context.UserBackend{DsType: "external", Server: "custom", BaseDN: "n", ConnectionDN: "n", Host: "n", Port: "2", Encryption: "ssl", Password: "n"}, "no attributeGivenName set", assert.Error, assert.Contains},
		{"attributeSurName not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "name", BaseDN: "n", ConnectionDN: "n", Host: "n", Port: "2", Encryption: "ssl", Password: "n"}, "no attributeSurName set", assert.Error, assert.Contains},
		{"baseDn not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n"}, "no baseDn set", assert.Error, assert.Contains},
		{"connectionDn not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n"}, "no connectionDn set", assert.Error, assert.Contains},
		{"password not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Host: "n", Port: "2", Encryption: "ssl"}, "no password set", assert.Error, assert.Contains},
		{"host not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n"}, "no host set", assert.Error, assert.Contains},
		{"port not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n"}, "no port set", assert.Error, assert.Contains},
		{"port is not number", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "asd"}, "failed to validate property port: the given value is not a number", assert.Error, assert.Contains},
		{"encryption not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2"}, "invalid encryption valid options are", assert.Error, assert.Contains},
		{"groupBaseDN not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2", Encryption: "ssl"}, "no groupBaseDN set", assert.Error, assert.Contains},
		{"groupSearchFilter not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2", Encryption: "ssl", GroupBaseDN: "n"}, "no groupSearchFilter set", assert.Error, assert.Contains},
		{"groupAttributeName not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n"}, "no groupAttributeName set", assert.Error, assert.Contains},
		{"groupAttributeDescription not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n"}, "no groupAttributeDescription set", assert.Error, assert.Contains},
		{"groupAttributeMember not set", context.UserBackend{DsType: "external", Server: "custom", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "2", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n", GroupAttributeDescription: "n"}, "no groupAttributeMember set", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &userBackendValidator{}

			// when
			result := validator.validateExternalBackend(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful external backend validation", func(t *testing.T) {
		// given
		userBackend := context.UserBackend{DsType: "external", Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "1", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n", GroupAttributeDescription: "n", GroupAttributeMember: "n"}
		validator := &userBackendValidator{}

		// when
		result := validator.validateExternalBackend(userBackend)

		// then
		require.NoError(t, result)
	})
}

func Test_userBackendValidator_validateActiveDirectoryServerBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      context.UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid attributeID", context.UserBackend{Server: "activeDirectory"}, "invalid attributeID valid options are", assert.Error, assert.Contains},
		{"invalid attributeFullName", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName"}, "invalid attributeFullName valid options are", assert.Error, assert.Contains},
		{"invalid attributeMail", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn"}, "invalid attributeMail valid options are", assert.Error, assert.Contains},
		{"invalid attributeGroup", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail"}, "invalid attributeGroup valid options are", assert.Error, assert.Contains},
		{"invalid searchFilter", context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf"}, "invalid searchFilter valid options are", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &userBackendValidator{}

			// when
			result := validator.validateActiveDirectoryServer(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful active directory server validation", func(t *testing.T) {
		// given
		userBackend := context.UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)"}
		validator := &userBackendValidator{}

		// when
		result := validator.validateActiveDirectoryServer(userBackend)

		// then
		require.NoError(t, result)
	})
}

func Test_userBackendValidator_ValidateUserBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      context.UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid datasource type", context.UserBackend{}, "invalid dsType valid options are [embedded external]", assert.Error, assert.Contains},
		{"invalid external config", context.UserBackend{DsType: "external"}, "invalid server valid options are", assert.Error, assert.Contains},
		{"invalid embedded config", context.UserBackend{DsType: "embedded"}, "invalid attributeID valid options are [uid]", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &userBackendValidator{}

			// when
			result := validator.ValidateUserBackend(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful external backend validation", func(t *testing.T) {
		// given
		userBackend := context.UserBackend{DsType: "external", Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "1", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n", GroupAttributeDescription: "n", GroupAttributeMember: "n"}
		validator := &userBackendValidator{}

		// when
		result := validator.ValidateUserBackend(userBackend)

		// then
		require.NoError(t, result)
	})
}
