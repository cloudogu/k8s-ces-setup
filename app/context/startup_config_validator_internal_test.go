package context

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStartupConfigurationValidator(t *testing.T) {
	t.Run("successfull creating validator", func(t *testing.T) {
		// when
		validator := NewStartupConfigurationValidator(SetupConfiguration{})

		// then
		require.NotNil(t, validator)
	})
}

func Test_validator_ValidateConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		config           SetupConfiguration
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid naming section", SetupConfiguration{Naming: Naming{Completed: true}}, "failed to validate naming section", assert.Error, assert.Contains},
		{"invalid user backend section", SetupConfiguration{Naming: Naming{Completed: false}, UserBackend: UserBackend{Completed: true}}, "failed to validate user userBackend section", assert.Error, assert.Contains},
		{"invalid admin user section", SetupConfiguration{Naming: Naming{Completed: false}, UserBackend: UserBackend{Completed: false}, Admin: User{Completed: true}}, "failed to validate admin user section", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			v := &validator{
				configuration: tt.config,
			}

			// when
			result := v.ValidateConfiguration()

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}
}

func Test_validator_validateEmbeddedBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid attributeID backend type embedded", UserBackend{DsType: "embedded"}, "invalid attributeID valid options are", assert.Error, assert.Contains},
		{"invalid attributeFullName backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid"}, "invalid attributeFullName valid options are", assert.Error, assert.Contains},
		{"invalid attributeMail backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn"}, "invalid attributeMail valid options are", assert.Error, assert.Contains},
		{"invalid attributeGroup backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail"}, "invalid attributeGroup valid options are", assert.Error, assert.Contains},
		{"invalid searchFilter backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf"}, "invalid searchFilter valid options are", assert.Error, assert.Contains},
		{"invalid host backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)"}, "invalid host valid options are", assert.Error, assert.Contains},
		{"invalid port backend type embedded", UserBackend{DsType: "embedded", AttributeID: "uid", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)", Host: "ldap"}, "invalid port valid options are", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &validator{configuration: SetupConfiguration{UserBackend: tt.userBackend}}

			// when
			result := validator.validateEmbeddedBackend(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}
}

func Test_validator_validateExternalBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid server", UserBackend{DsType: "external"}, "invalid server valid options are", assert.Error, assert.Contains},
		{"attributeGivenName not set", UserBackend{DsType: "external", Server: "activeDirectory"}, "no attributeGivenName set", assert.Error, assert.Contains},
		{"attributeSurName not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "name"}, "no attributeSurName set", assert.Error, assert.Contains},
		{"baseDn not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n"}, "no baseDn set", assert.Error, assert.Contains},
		{"connectionDn not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n"}, "no connectionDn set", assert.Error, assert.Contains},
		{"password not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n"}, "no password set", assert.Error, assert.Contains},
		{"host not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n"}, "no host set", assert.Error, assert.Contains},
		{"host not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n"}, "no port set", assert.Error, assert.Contains},
		{"encryption not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n"}, "invalid encryption valid options are", assert.Error, assert.Contains},
		{"groupBaseDN not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl"}, "no groupBaseDN set", assert.Error, assert.Contains},
		{"groupSearchFilter not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl", GroupBaseDN: "n"}, "no groupSearchFilter set", assert.Error, assert.Contains},
		{"groupAttributeName not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n"}, "no groupAttributeName set", assert.Error, assert.Contains},
		{"groupAttributeDescription not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n"}, "no groupAttributeDescription set", assert.Error, assert.Contains},
		{"groupAttributeMember not set", UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n", GroupAttributeDescription: "n"}, "no groupAttributeMember set", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &validator{configuration: SetupConfiguration{UserBackend: tt.userBackend}}

			// when
			result := validator.validateExternalBackend(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful external backend validation", func(t *testing.T) {
		// given
		userBackend := UserBackend{DsType: "external", Server: "activeDirectory", AttributeGivenName: "n", AttributeSurname: "n", BaseDN: "n", ConnectionDN: "n", Password: "n", Host: "n", Port: "n", Encryption: "ssl", GroupBaseDN: "n", GroupSearchFilter: "n", GroupAttributeName: "n", GroupAttributeDescription: "n", GroupAttributeMember: "n"}
		validator := &validator{configuration: SetupConfiguration{UserBackend: userBackend}}

		// when
		result := validator.validateExternalBackend(userBackend)

		// then
		require.NoError(t, result)
	})
}

func Test_validator_validateActiveDirectoryServerBackend(t *testing.T) {
	tests := []struct {
		name             string
		userBackend      UserBackend
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid attributeID", UserBackend{Server: "activeDirectory"}, "invalid attributeID valid options are", assert.Error, assert.Contains},
		{"invalid attributeFullName", UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName"}, "invalid attributeFullName valid options are", assert.Error, assert.Contains},
		{"invalid attributeMail", UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn"}, "invalid attributeMail valid options are", assert.Error, assert.Contains},
		{"invalid attributeGroup", UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail"}, "invalid attributeGroup valid options are", assert.Error, assert.Contains},
		{"invalid searchFilter", UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf"}, "invalid searchFilter valid options are", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &validator{configuration: SetupConfiguration{UserBackend: tt.userBackend}}

			// when
			result := validator.validateActiveDirectoryServer(tt.userBackend)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful active directory server validation", func(t *testing.T) {
		// given
		userBackend := UserBackend{Server: "activeDirectory", AttributeID: "sAMAccountName", AttributeFullname: "cn", AttributeMail: "mail", AttributeGroup: "memberOf", SearchFilter: "(objectClass=person)"}
		validator := &validator{configuration: SetupConfiguration{UserBackend: userBackend}}

		// when
		result := validator.validateActiveDirectoryServer(userBackend)

		// then
		require.NoError(t, result)
	})
}

func Test_validator_validateAdminUser(t *testing.T) {
	tests := []struct {
		name             string
		user             User
		dsType           string
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"no admin group set", User{}, "", "no admin group set", assert.Error, assert.Contains},
		{"no admin mail set", User{AdminGroup: "group"}, dsTypeEmbedded, "no admin mail set", assert.Error, assert.Contains},
		{"invalid admin mail set", User{AdminGroup: "group", Mail: "t"}, dsTypeEmbedded, "invalid admin mail", assert.Error, assert.Contains},
		{"no admin username set", User{AdminGroup: "group", Mail: "test@test.de"}, dsTypeEmbedded, "no admin username set", assert.Error, assert.Contains},
		{"no admin password set", User{AdminGroup: "group", Mail: "test@test.de", Username: "name"}, dsTypeEmbedded, "no admin password set", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &validator{}

			// when
			result := validator.validateAdminUser(tt.user, tt.dsType)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful active directory server validation", func(t *testing.T) {
		// given
		adminUser := User{AdminGroup: "group", Mail: "test@test.de", Username: "name", Password: "password"}
		validator := &validator{configuration: SetupConfiguration{Admin: adminUser}}

		// when
		result := validator.validateAdminUser(adminUser, dsTypeEmbedded)

		// then
		require.NoError(t, result)
	})
}

func Test_validator_validateNaming(t *testing.T) {
	cert := "-----BEGIN CERTIFICATE-----\nMIIFTzCCBDegAwIBAgIFFlF0N0AwDQYJKoZIhvcNAQELBQAwgYAxCzAJBgNVBAYT\nAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEV\nMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAW\nBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAeFw0yMjA1MDUwOTQyMjBaFw00NjEyMjUw\nOTQyMjBaMH0xCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQ\nBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQL\nDAwxOTIuMTY4LjU2LjIxFTATBgNVBAMMDDE5Mi4xNjguNTYuMjCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBAJ6dGTm0S/K0R0eRwj8KWLFnAhG4uY7jqK4t\npy42kczMGLnMCmO4qPZNmOk6zb/hqwBuU6GzDViNNMS5H3PTJ4er7cWomwGmRT93\ngcM26JrXhBqRI+BcdUM4ldswSZrViNNn2jf+X7LetsjoiUjDwkG1Ye28RaP9wCCh\nz+6Aht6PgMAavkBJR488fohcdmTV4Sv01Wv6iNjhoW1jJr/QoBq7GRIXwv3TUMLh\nLqoxgJ9946oRCRexO+oARlETPIonmUTtSzWiYdhiAoVydNiupXqmCF8EfQYGa3cy\noSTdPwj3M79ntxWZ1FzKaZ9ddR4W7nxBWsZqW5eYZ1UWZtevT7S+W1mxUDnDsJeP\n8gTh7rcyrDHlKHNmOvJlMo3qOBYBRJdGEYABYpz3ToiZsioyre1ORcCZhehs4yn9\nzClBFkqv2uHjY8Ucc+CBvxK6FayXyjXKDPtkBeCp+UAm4VLH8seNlUKwRCylNN4y\n6PP4CY83yUFiGlKG7f5z9zKPJLMYmPWyXVejTaOiIpZ5YkMqoK57p+bCZaAMc+6M\newWzhBJDy6kJFkiP27zWImFa3CtzDPwb/IyGkFBh5KWiQlTyFecd94+Vv2k/G1UD\nq2BVsOo+rt0U+CN1MD5vskOQ8hCa0Xk7aQcXJRR5KryzHRApl+FC7GUnOtKeQRar\nEMy4Pt8lAgMBAAGjgdEwgc4wCQYDVR0TBAIwADAsBglghkgBhvhCAQ0EHxYdT3Bl\nblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYDVR0OBBYEFOVGdYUJGPBPgrku\nWZMraMhm7plYMB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MBMGA1Ud\nJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwIF4DAxBgNVHREEKjAoggwxOTIuMTY4\nLjU2LjKCEmxvY2FsLmNsb3Vkb2d1LmNvbYcEwKg4AjANBgkqhkiG9w0BAQsFAAOC\nAQEAPBJFxh4n0YayAQkxvAcmGACn09ugczCRWlPCylgORxcD7mJdqr61/LMie0iY\n7OnFMggl2xlx+8yrfwtTEzeBNzraOYJKkFeBnZ3yxC63oRminOdgClUDA16D7Guk\nDJ94gJy6ueIA+MXbWkEg5w+suGUCovbJDATnjiAP+xQ3tK4GtyACibP0tHNFzUTe\n3GNTqkSJnV9rjjN7NFfEe+nSQFLghz/nP9k/vyECFyjemG8k5Vd1XNogs13uXpSG\nX4Q78vRC+s2QIm3ZIokh3Uu4bKK4Rl9aRMynt8iJ7ZxlK0+/pJpI8e6yKDdNpe38\nvK6monD5jYOdcYWmqUwh/wgseQ==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4jCCAsqgAwIBAgIUY1X2vb/NYhqWg+hBq5hZncWVHa8wDQYJKoZIhvcNAQEN\nBQAwgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNV\nBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwx\nOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAgFw0yMjA1MDUw\nOTQyMTlaGA8yMDkwMDUyMzA5NDIxOVowgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQI\nDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTky\nLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBT\nZWxmIFNpZ25lZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKFmyOfn\nXnpesxiqApTUSBbO5fg/GhcRCFI2n/kNsezGHv+w1j47kuP9wE6kRYiGGmjD36lU\ndX4abq2UppeWiGycceT5oXdKfLlQP7J2jNPTiPstGXfEk6mGnzzyDz8VXd8EsWfc\nMRPcJyC9l0MXRPuagqnKIipIOEWeqsnuM7IQS62SmfTlBt8MVehlMLoo3L61wH3E\nyLSicZvwCvkBUWowa0K3sStoUyCm8TIOIjPyGaOTmbjWLqkrSoKbhuGvbXXAfJX3\nyup5lsDCl9jAznXGTGJ5ZuAmWUHlbgkO324/9YGhTZUHTErkmTnZ7bHwhFLACHAt\nj3J251HIfxa6eEcCAwEAAaNQME4wHQYDVR0OBBYEFCI01F+3w7czJ4NLLjetRkxR\n+ed9MB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MAwGA1UdEwQFMAMB\nAf8wDQYJKoZIhvcNAQENBQADggEBABjFa2Mja4KQMXWBtGTEMfhJmahU63k5lyNO\nObH/JoepywTTekjzTqNj1qt3vR1+ITNsFHZ36VFcP1BQQBn0v9SUTMGFwmF40hYG\nCreFvO7HBZdsQhCtOfv0tq9gA3NTght+vhl0rQSWPSf3I87xywQFti4OM4kkPKCg\nbbCVsw512o3PLzOMWPolg3LmXH2sJwGe/i9fAFO8twEvcqynC0z2BLnlrucfpDvD\nKM5olC9qszMj6MT3vqMKC112isadYUqn860G4EwUpjj7PH2kQneori9K8BKX+Qx3\nI+keZoQ4jX56rm5W9+IqiXUGz1xpADIbIB6KIQmMMec03z1aX24=\n-----END CERTIFICATE-----\n"
	invalidCert := "-----BEGIN CERTIFICATE-----\nMIIFTzCCBDekRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEV\nMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAW\nBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAeFw0yMjA1MDUwOTQyMjBaFw00NjEyMjUw\nOTQyMjBaMH0xCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQ\nBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQL\nDAwxOTIuMTY4LjU2LjIxFTATBgNVBAMMDDE5Mi4xNjguNTYuMjCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBAJ6dGTm0S/K0R0eRwj8KWLFnAhG4uY7jqK4t\npy42kczMGLnMCmO4qPZNmOk6zb/hqwBuU6GzDViNNMS5H3PTJ4er7cWomwGmRT93\ngcM26JrXhBqRI+BcdUM4ldswSZrViNNn2jf+X7LetsjoiUjDwkG1Ye28RaP9wCCh\nz+6Aht6PgMAavkBJR488fohcdmTV4Sv01Wv6iNjhoW1jJr/QoBq7GRIXwv3TUMLh\nLqoxgJ9946oRCRexO+oARlETPIonmUTtSzWiYdhiAoVydNiupXqmCF8EfQYGa3cy\noSTdPwj3M79ntxWZ1FzKaZ9ddR4W7nxBWsZqW5eYZ1UWZtevT7S+W1mxUDnDsJeP\n8gTh7rcyrDHlKHNmOvJlMo3qOBYBRJdGEYABYpz3ToiZsioyre1ORcCZhehs4yn9\nzClBFkqv2uHjY8Ucc+CBvxK6FayXyjXKDPtkBeCp+UAm4VLH8seNlUKwRCylNN4y\n6PP4CY83yUFiGlKG7f5z9zKPJLMYmPWyXVejTaOiIpZ5YkMqoK57p+bCZaAMc+6M\newWzhBJDy6kJFkiP27zWImFa3CtzDPwb/IyGkFBh5KWiQlTyFecd94+Vv2k/G1UD\nq2BVsOo+rt0U+CN1MD5vskOQ8hCa0Xk7aQcXJRR5KryzHRApl+FC7GUnOtKeQRar\nEMy4Pt8lAgMBAAGjgdEwgc4wCQYDVR0TBAIwADAsBglghkgBhvhCAQ0EHxYdT3Bl\nblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYDVR0OBBYEFOVGdYUJGPBPgrku\nWZMraMhm7plYMB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MBMGA1Ud\nJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwIF4DAxBgNVHREEKjAoggwxOTIuMTY4\nLjU2LjKCEmxvY2FsLmNsb3Vkb2d1LmNvbYcEwKg4AjANBgkqhkiG9w0BAQsFAAOC\nAQEAPBJFxh4n0YayAQkxvAcmGACn09ugczCRWlPCylgORxcD7mJdqr61/LMie0iY\n7OnFMggl2xlx+8yrfwtTEzeBNzraOYJKkFeBnZ3yxC63oRminOdgClUDA16D7Guk\nDJ94gJy6ueIA+MXbWkEg5w+suGUCovbJDATnjiAP+xQ3tK4GtyACibP0tHNFzUTe\n3GNTqkSJnV9rjjN7NFfEe+nSQFLghz/nP9k/vyECFyjemG8k5Vd1XNogs13uXpSG\nX4Q78vRC+s2QIm3ZIokh3Uu4bKK4Rl9aRMynt8iJ7ZxlK0+/pJpI8e6yKDdNpe38\nvK6monD5jYOdcYWmqUwh/wgseQ==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4jCCAsqgAwIBAgIUY1X2vb/NYhqWg+hBq5hZncWVHa8wDQYJKoZIhvcNAQEN\nBQAwgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNV\nBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwx\nOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAgFw0yMjA1MDUw\nOTQyMTlaGA8yMDkwMDUyMzA5NDIxOVowgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQI\nDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTky\nLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBT\nZWxmIFNpZ25lZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKFmyOfn\nXnpesxiqApTUSBbO5fg/GhcRCFI2n/kNsezGHv+w1j47kuP9wE6kRYiGGmjD36lU\ndX4abq2UppeWiGycceT5oXdKfLlQP7J2jNPTiPstGXfEk6mGnzzyDz8VXd8EsWfc\nMRPcJyC9l0MXRPuagqnKIipIOEWeqsnuM7IQS62SmfTlBt8MVehlMLoo3L61wH3E\nyLSicZvwCvkBUWowa0K3sStoUyCm8TIOIjPyGaOTmbjWLqkrSoKbhuGvbXXAfJX3\nyup5lsDCl9jAznXGTGJ5ZuAmWUHlbgkO324/9YGhTZUHTErkmTnZ7bHwhFLACHAt\nj3J251HIfxa6eEcCAwEAAaNQME4wHQYDVR0OBBYEFCI01F+3w7czJ4NLLjetRkxR\n+ed9MB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MAwGA1UdEwQFMAMB\nAf8wDQYJKoZIhvcNAQENBQADggEBABjFa2Mja4KQMXWBtGTEMfhJmahU63k5lyNO\nObH/JoepywTTekjzTqNj1qt3vR1+ITNsFHZ36VFcP1BQQBn0v9SUTMGFwmF40hYG\nCreFvO7HBZdsQhCtOfv0tq9gA3NTght+vhl0rQSWPSf3I87xywQFti4OM4kkPKCg\nbbCVsw512o3PLzOMWPolg3LmXH2sJwGe/i9fAFO8twEvcqynC0z2BLnlrucfpDvD\nKM5olC9qszMj6MT3vqMKC112isadYUqn860G4EwUpjj7PH2kQneori9K8BKX+Qx3\nI+keZoQ4jX56rm5W9+IqiXUGz1xpADIbIB6KIQmMMec03z1aX24=\n----END CERTIFICATE-----\n"
	tests := []struct {
		name             string
		naming           Naming
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid fqdn", Naming{Fqdn: "1.2"}, "failed to parse fqdn", assert.Error, assert.Contains},
		{"invalid domain", Naming{Fqdn: "192.168.56.2"}, "failed to validate domain", assert.Error, assert.Contains},
		{"invalid certificateType", Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com"}, "invalid certificateType valid options are", assert.Error, assert.Contains},
		{"invalid relayhost", Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "@_"}, "failed to validate mail relay host", assert.Error, assert.Contains},
		{"invalid mail address", Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b@a"}, "failed to validate mail address", assert.Error, assert.Contains},
		{"invalid internal ip", Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de", UseInternalIp: true, InternalIp: "1234.123"}, "failed to parse internal ip", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &validator{}

			// when
			result := validator.validateNaming(tt.naming)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful naming validation with ip", func(t *testing.T) {
		// given
		naming := Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de"}
		validator := &validator{}

		// when
		result := validator.validateNaming(naming)

		// then
		require.NoError(t, result)
	})

	t.Run("successful naming validation with dns", func(t *testing.T) {
		// given
		naming := Naming{Fqdn: "cloudogu.com", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de"}
		validator := &validator{}

		// when
		result := validator.validateNaming(naming)

		// then
		require.NoError(t, result)
	})

	t.Run("successful naming validation with external certificate", func(t *testing.T) {
		// given
		naming := Naming{Fqdn: "cloudogu.com", Domain: "cloudogu.com", CertificateType: "external", RelayHost: "relay", MailAddress: "a@b.de", Certificate: cert, CertificateKey: "key"}
		validator := &validator{}

		// when
		result := validator.validateNaming(naming)

		// then
		require.NoError(t, result)
	})

	t.Run("successful naming validation with external certificate", func(t *testing.T) {
		// given
		naming := Naming{Fqdn: "cloudogu.com", Domain: "cloudogu.com", CertificateType: "external", RelayHost: "relay", MailAddress: "a@b.de", Certificate: invalidCert, CertificateKey: "key"}
		validator := &validator{}

		// when
		result := validator.validateNaming(naming)

		// then
		require.NoError(t, result)
	})
}
