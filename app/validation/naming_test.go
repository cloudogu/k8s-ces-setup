package validation

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/v4/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNamingValidator(t *testing.T) {
	// when
	validator := NewNamingValidator()

	// then
	require.NotNil(t, validator)
}

func Test_namingValidator_ValidateNaming(t *testing.T) {
	cert := "-----BEGIN CERTIFICATE-----\nMIIFTzCCBDegAwIBAgIFFlF0N0AwDQYJKoZIhvcNAQELBQAwgYAxCzAJBgNVBAYT\nAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEV\nMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAW\nBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAeFw0yMjA1MDUwOTQyMjBaFw00NjEyMjUw\nOTQyMjBaMH0xCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQ\nBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQL\nDAwxOTIuMTY4LjU2LjIxFTATBgNVBAMMDDE5Mi4xNjguNTYuMjCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBAJ6dGTm0S/K0R0eRwj8KWLFnAhG4uY7jqK4t\npy42kczMGLnMCmO4qPZNmOk6zb/hqwBuU6GzDViNNMS5H3PTJ4er7cWomwGmRT93\ngcM26JrXhBqRI+BcdUM4ldswSZrViNNn2jf+X7LetsjoiUjDwkG1Ye28RaP9wCCh\nz+6Aht6PgMAavkBJR488fohcdmTV4Sv01Wv6iNjhoW1jJr/QoBq7GRIXwv3TUMLh\nLqoxgJ9946oRCRexO+oARlETPIonmUTtSzWiYdhiAoVydNiupXqmCF8EfQYGa3cy\noSTdPwj3M79ntxWZ1FzKaZ9ddR4W7nxBWsZqW5eYZ1UWZtevT7S+W1mxUDnDsJeP\n8gTh7rcyrDHlKHNmOvJlMo3qOBYBRJdGEYABYpz3ToiZsioyre1ORcCZhehs4yn9\nzClBFkqv2uHjY8Ucc+CBvxK6FayXyjXKDPtkBeCp+UAm4VLH8seNlUKwRCylNN4y\n6PP4CY83yUFiGlKG7f5z9zKPJLMYmPWyXVejTaOiIpZ5YkMqoK57p+bCZaAMc+6M\newWzhBJDy6kJFkiP27zWImFa3CtzDPwb/IyGkFBh5KWiQlTyFecd94+Vv2k/G1UD\nq2BVsOo+rt0U+CN1MD5vskOQ8hCa0Xk7aQcXJRR5KryzHRApl+FC7GUnOtKeQRar\nEMy4Pt8lAgMBAAGjgdEwgc4wCQYDVR0TBAIwADAsBglghkgBhvhCAQ0EHxYdT3Bl\nblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYDVR0OBBYEFOVGdYUJGPBPgrku\nWZMraMhm7plYMB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MBMGA1Ud\nJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwIF4DAxBgNVHREEKjAoggwxOTIuMTY4\nLjU2LjKCEmxvY2FsLmNsb3Vkb2d1LmNvbYcEwKg4AjANBgkqhkiG9w0BAQsFAAOC\nAQEAPBJFxh4n0YayAQkxvAcmGACn09ugczCRWlPCylgORxcD7mJdqr61/LMie0iY\n7OnFMggl2xlx+8yrfwtTEzeBNzraOYJKkFeBnZ3yxC63oRminOdgClUDA16D7Guk\nDJ94gJy6ueIA+MXbWkEg5w+suGUCovbJDATnjiAP+xQ3tK4GtyACibP0tHNFzUTe\n3GNTqkSJnV9rjjN7NFfEe+nSQFLghz/nP9k/vyECFyjemG8k5Vd1XNogs13uXpSG\nX4Q78vRC+s2QIm3ZIokh3Uu4bKK4Rl9aRMynt8iJ7ZxlK0+/pJpI8e6yKDdNpe38\nvK6monD5jYOdcYWmqUwh/wgseQ==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4jCCAsqgAwIBAgIUY1X2vb/NYhqWg+hBq5hZncWVHa8wDQYJKoZIhvcNAQEN\nBQAwgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNV\nBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwx\nOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAgFw0yMjA1MDUw\nOTQyMTlaGA8yMDkwMDUyMzA5NDIxOVowgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQI\nDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTky\nLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBT\nZWxmIFNpZ25lZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKFmyOfn\nXnpesxiqApTUSBbO5fg/GhcRCFI2n/kNsezGHv+w1j47kuP9wE6kRYiGGmjD36lU\ndX4abq2UppeWiGycceT5oXdKfLlQP7J2jNPTiPstGXfEk6mGnzzyDz8VXd8EsWfc\nMRPcJyC9l0MXRPuagqnKIipIOEWeqsnuM7IQS62SmfTlBt8MVehlMLoo3L61wH3E\nyLSicZvwCvkBUWowa0K3sStoUyCm8TIOIjPyGaOTmbjWLqkrSoKbhuGvbXXAfJX3\nyup5lsDCl9jAznXGTGJ5ZuAmWUHlbgkO324/9YGhTZUHTErkmTnZ7bHwhFLACHAt\nj3J251HIfxa6eEcCAwEAAaNQME4wHQYDVR0OBBYEFCI01F+3w7czJ4NLLjetRkxR\n+ed9MB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MAwGA1UdEwQFMAMB\nAf8wDQYJKoZIhvcNAQENBQADggEBABjFa2Mja4KQMXWBtGTEMfhJmahU63k5lyNO\nObH/JoepywTTekjzTqNj1qt3vR1+ITNsFHZ36VFcP1BQQBn0v9SUTMGFwmF40hYG\nCreFvO7HBZdsQhCtOfv0tq9gA3NTght+vhl0rQSWPSf3I87xywQFti4OM4kkPKCg\nbbCVsw512o3PLzOMWPolg3LmXH2sJwGe/i9fAFO8twEvcqynC0z2BLnlrucfpDvD\nKM5olC9qszMj6MT3vqMKC112isadYUqn860G4EwUpjj7PH2kQneori9K8BKX+Qx3\nI+keZoQ4jX56rm5W9+IqiXUGz1xpADIbIB6KIQmMMec03z1aX24=\n-----END CERTIFICATE-----\n"
	key := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
	invalidCert := "-----BEGIN CERTIFICATE-----\nMIIFTzCCBDekRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEV\nMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAW\nBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAeFw0yMjA1MDUwOTQyMjBaFw00NjEyMjUw\nOTQyMjBaMH0xCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQ\nBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQL\nDAwxOTIuMTY4LjU2LjIxFTATBgNVBAMMDDE5Mi4xNjguNTYuMjCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBAJ6dGTm0S/K0R0eRwj8KWLFnAhG4uY7jqK4t\npy42kczMGLnMCmO4qPZNmOk6zb/hqwBuU6GzDViNNMS5H3PTJ4er7cWomwGmRT93\ngcM26JrXhBqRI+BcdUM4ldswSZrViNNn2jf+X7LetsjoiUjDwkG1Ye28RaP9wCCh\nz+6Aht6PgMAavkBJR488fohcdmTV4Sv01Wv6iNjhoW1jJr/QoBq7GRIXwv3TUMLh\nLqoxgJ9946oRCRexO+oARlETPIonmUTtSzWiYdhiAoVydNiupXqmCF8EfQYGa3cy\noSTdPwj3M79ntxWZ1FzKaZ9ddR4W7nxBWsZqW5eYZ1UWZtevT7S+W1mxUDnDsJeP\n8gTh7rcyrDHlKHNmOvJlMo3qOBYBRJdGEYABYpz3ToiZsioyre1ORcCZhehs4yn9\nzClBFkqv2uHjY8Ucc+CBvxK6FayXyjXKDPtkBeCp+UAm4VLH8seNlUKwRCylNN4y\n6PP4CY83yUFiGlKG7f5z9zKPJLMYmPWyXVejTaOiIpZ5YkMqoK57p+bCZaAMc+6M\newWzhBJDy6kJFkiP27zWImFa3CtzDPwb/IyGkFBh5KWiQlTyFecd94+Vv2k/G1UD\nq2BVsOo+rt0U+CN1MD5vskOQ8hCa0Xk7aQcXJRR5KryzHRApl+FC7GUnOtKeQRar\nEMy4Pt8lAgMBAAGjgdEwgc4wCQYDVR0TBAIwADAsBglghkgBhvhCAQ0EHxYdT3Bl\nblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYDVR0OBBYEFOVGdYUJGPBPgrku\nWZMraMhm7plYMB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MBMGA1Ud\nJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwIF4DAxBgNVHREEKjAoggwxOTIuMTY4\nLjU2LjKCEmxvY2FsLmNsb3Vkb2d1LmNvbYcEwKg4AjANBgkqhkiG9w0BAQsFAAOC\nAQEAPBJFxh4n0YayAQkxvAcmGACn09ugczCRWlPCylgORxcD7mJdqr61/LMie0iY\n7OnFMggl2xlx+8yrfwtTEzeBNzraOYJKkFeBnZ3yxC63oRminOdgClUDA16D7Guk\nDJ94gJy6ueIA+MXbWkEg5w+suGUCovbJDATnjiAP+xQ3tK4GtyACibP0tHNFzUTe\n3GNTqkSJnV9rjjN7NFfEe+nSQFLghz/nP9k/vyECFyjemG8k5Vd1XNogs13uXpSG\nX4Q78vRC+s2QIm3ZIokh3Uu4bKK4Rl9aRMynt8iJ7ZxlK0+/pJpI8e6yKDdNpe38\nvK6monD5jYOdcYWmqUwh/wgseQ==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4jCCAsqgAwIBAgIUY1X2vb/NYhqWg+hBq5hZncWVHa8wDQYJKoZIhvcNAQEN\nBQAwgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBTYXhvbnkxEjAQBgNV\nBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTkyLjE2OC41Ni4yMRUwEwYDVQQLDAwx\nOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBTZWxmIFNpZ25lZDAgFw0yMjA1MDUw\nOTQyMTlaGA8yMDkwMDUyMzA5NDIxOVowgYAxCzAJBgNVBAYTAkRFMRUwEwYDVQQI\nDAxMb3dlciBTYXhvbnkxEjAQBgNVBAcMCUJydW5zd2ljazEVMBMGA1UECgwMMTky\nLjE2OC41Ni4yMRUwEwYDVQQLDAwxOTIuMTY4LjU2LjIxGDAWBgNVBAMMD0NFUyBT\nZWxmIFNpZ25lZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKFmyOfn\nXnpesxiqApTUSBbO5fg/GhcRCFI2n/kNsezGHv+w1j47kuP9wE6kRYiGGmjD36lU\ndX4abq2UppeWiGycceT5oXdKfLlQP7J2jNPTiPstGXfEk6mGnzzyDz8VXd8EsWfc\nMRPcJyC9l0MXRPuagqnKIipIOEWeqsnuM7IQS62SmfTlBt8MVehlMLoo3L61wH3E\nyLSicZvwCvkBUWowa0K3sStoUyCm8TIOIjPyGaOTmbjWLqkrSoKbhuGvbXXAfJX3\nyup5lsDCl9jAznXGTGJ5ZuAmWUHlbgkO324/9YGhTZUHTErkmTnZ7bHwhFLACHAt\nj3J251HIfxa6eEcCAwEAAaNQME4wHQYDVR0OBBYEFCI01F+3w7czJ4NLLjetRkxR\n+ed9MB8GA1UdIwQYMBaAFCI01F+3w7czJ4NLLjetRkxR+ed9MAwGA1UdEwQFMAMB\nAf8wDQYJKoZIhvcNAQENBQADggEBABjFa2Mja4KQMXWBtGTEMfhJmahU63k5lyNO\nObH/JoepywTTekjzTqNj1qt3vR1+ITNsFHZ36VFcP1BQQBn0v9SUTMGFwmF40hYG\nCreFvO7HBZdsQhCtOfv0tq9gA3NTght+vhl0rQSWPSf3I87xywQFti4OM4kkPKCg\nbbCVsw512o3PLzOMWPolg3LmXH2sJwGe/i9fAFO8twEvcqynC0z2BLnlrucfpDvD\nKM5olC9qszMj6MT3vqMKC112isadYUqn860G4EwUpjj7PH2kQneori9K8BKX+Qx3\nI+keZoQ4jX56rm5W9+IqiXUGz1xpADIbIB6KIQmMMec03z1aX24=\n----END CERTIFICATE-----\n"
	tests := []struct {
		name             string
		naming           context.Naming
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"invalid fqdn", context.Naming{}, "no fqdn set", assert.Error, assert.Contains},
		{"invalid domain", context.Naming{Fqdn: "192.168.56.2"}, "no domain set", assert.Error, assert.Contains},
		{"invalid certificateType", context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com"}, "invalid certificateType valid options are", assert.Error, assert.Contains},
		{"invalid relayhost", context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned"}, "no relayHost set", assert.Error, assert.Contains},
		{"invalid mail address", context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b@a"}, "failed to validate mail address", assert.Error, assert.Contains},
		{"invalid internal ip", context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de", UseInternalIp: true, InternalIp: "1234.123"}, "failed to parse internal ip", assert.Error, assert.Contains},
		{"invalid external cert", context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "external", Certificate: invalidCert}, "failed to decode 0-th certificate in [certificate] property", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &namingValidator{}

			// when
			result := validator.ValidateNaming(tt.naming)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful naming validation with ip", func(t *testing.T) {
		// given
		naming := context.Naming{Fqdn: "192.168.56.2", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de"}
		validator := &namingValidator{}

		// when
		result := validator.ValidateNaming(naming)

		// then
		require.NoError(t, result)
	})

	t.Run("successful naming validation with dns", func(t *testing.T) {
		// given
		naming := context.Naming{Fqdn: "cloudogu.com", Domain: "cloudogu.com", CertificateType: "selfsigned", RelayHost: "relay", MailAddress: "a@b.de"}
		validator := &namingValidator{}

		// when
		result := validator.ValidateNaming(naming)

		// then
		require.NoError(t, result)
	})

	t.Run("successful naming validation with external certificate", func(t *testing.T) {
		// given
		naming := context.Naming{Fqdn: "cloudogu.com", Domain: "cloudogu.com", CertificateType: "external", RelayHost: "relay", MailAddress: "a@b.de", Certificate: cert, CertificateKey: key}
		validator := &namingValidator{}

		// when
		result := validator.ValidateNaming(naming)

		// then
		require.NoError(t, result)
	})
}
