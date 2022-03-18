package setup

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_fetchYaml(t *testing.T) {
	t.Run("should return error on HTTP not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("HTTP 404 - that's an error"))
		}))
		defer server.Close()

		// when
		_, err := fetchYaml(server.URL, &http.Client{})

		// then
		require.Error(t, err)
		assert.Regexp(t, regexp.MustCompile("could not find YAML file '.+'"), err.Error())
	})

	t.Run("should return file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(multiFileYaml()))
		}))
		defer server.Close()

		// when
		actual, err := fetchYaml(server.URL, &http.Client{})

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Equal(t, []byte(multiFileYaml()), actual)
	})
}

func Test_splitYamlFileSections(t *testing.T) {
	t.Run("should return two sections (with leading delimiter)", func(t *testing.T) {
		const simpleMultiLineYaml = `---
test:
---
anotherTest:
`
		input := []byte(simpleMultiLineYaml)

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, sections[0], "test:\n")
		assert.Equal(t, sections[1], "anotherTest:\n")
	})
	t.Run("should return two sections (without leading delimiter)", func(t *testing.T) {
		const simpleMultiLineYaml = `test:
---
anotherTest:
`
		input := []byte(simpleMultiLineYaml)

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, sections[0], "test:\n")
		assert.Equal(t, sections[1], "anotherTest:\n")
	})
	t.Run("should return sections for complex YAML", func(t *testing.T) {
		input := []byte(multiFileYaml())

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, sections[0], `# A comment for the service
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
`)
		assert.Equal(t, sections[1], `# a comment for the deployment
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
`)
	})
}

func Test_parseYamlResource(t *testing.T) {
	t.Run("should parse valid operator YAML", func(t *testing.T) {
		input := `apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
  labels:
    control-plane: controller-manager
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
  replicas: 1`

		// when
		parseYamlResource(input)

		// then
		//require.NoError(t, err)
		//assert.Equal(t, "", actual)
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
