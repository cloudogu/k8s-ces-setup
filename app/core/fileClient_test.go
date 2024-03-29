package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestNewFileClient(t *testing.T) {
	actual := NewFileClient("1.2.3")

	require.NotNil(t, actual)
	assert.NotNil(t, actual.httpClient)
	assert.Equal(t, "1.2.3", actual.version)
}

func Test_fileClient_Get(t *testing.T) {
	t.Run("should return error on HTTP not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("HTTP 404 - that's an error"))
		}))
		defer server.Close()

		sut := defaultHttpClient{httpClient: &http.Client{}}

		// when
		_, err := sut.Get(server.URL, "", "")

		// then
		require.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("response for YAML file '.+' returned with non-200 reply"), err.Error())
	})

	t.Run("should return file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(multiFileYaml()))
		}))
		defer server.Close()

		sut := defaultHttpClient{httpClient: &http.Client{}}

		// when
		actual, err := sut.Get(server.URL, "", "")

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Equal(t, []byte(multiFileYaml()), actual)
	})

	t.Run("should return file with basic auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			require.True(t, ok)
			assert.Equal(t, "username", username)
			assert.Equal(t, "password", password)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(multiFileYaml()))
		}))
		defer server.Close()

		sut := defaultHttpClient{httpClient: &http.Client{}}

		// when
		actual, err := sut.Get(server.URL, "username", "password")

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Equal(t, []byte(multiFileYaml()), actual)
	})
}

func multiFileYaml() string {
	return `---
# A comment for the service
apiVersion: v1
kind: Service
metadata:
  name: your-app
  app.kubernetes.io/name: your-app
  labels:
    app: your-app
spec:
  type: NodePort
  ports:
---
# a comment for the deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: your-app
  name: your-app
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: your-app
  template:
    metadata:
      labels:
        app: your-app
        app.kubernetes.io/name: your-app
    spec:
`
}
