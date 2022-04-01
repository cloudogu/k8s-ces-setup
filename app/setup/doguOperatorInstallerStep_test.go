package setup

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testTargetNamespaceName = "myfavouritenamespace-1"

func Test_replaceNamespace(t *testing.T) {
	t.Run("should replace namespace within simple namespaced resource", func(t *testing.T) {
		input := simpleTestRoleBinding()
		// when
		actual := replaceNamespacedResources(input, testTargetNamespaceName)

		// then
		expected := fmt.Sprintf(`apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8s-dogu-operator-leader-election-rolebinding
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8s-dogu-operator-leader-election-role
subjects:
- kind: ServiceAccount
  name: k8s-dogu-operator-controller-manager
  namespace: %s
`, testTargetNamespaceName, testTargetNamespaceName)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("should replace namespace in complex resources", func(t *testing.T) {
		input := twoNamespacedResources()

		// when
		actual := replaceNamespacedResources(input, testTargetNamespaceName)

		// then
		expected := fmt.Sprintf(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-dogu-operator-controller-manager
  namespace: %s
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8s-dogu-operator-leader-election-role
  namespace: %s
rules:
- apiGroups:
`,
			testTargetNamespaceName, testTargetNamespaceName)

		assert.Equal(t, expected, string(actual))
	})
}

func Test_removeLegacyNamespaceFromResources(t *testing.T) {
	input := []byte(`apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: ecosystem
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.8.0
  creationTimestamp: null
  name: dogus.k8s.cloudogu.com`)

	// when
	actual := removeLegacyNamespaceFromResources(input)

	// then
	expectedNoResource := `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.8.0
  creationTimestamp: null
  name: dogus.k8s.cloudogu.com`
	assert.Equal(t, expectedNoResource, string(actual))
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
		assert.Equal(t, "test:\n", string(sections[0]))
		assert.Equal(t, "anotherTest:\n", string(sections[1]))
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
		assert.Equal(t, "test:\n", string(sections[0]))
		assert.Equal(t, "anotherTest:\n", string(sections[1]))
	})
	t.Run("should return sections for complex YAML", func(t *testing.T) {
		input := []byte(multiFileYaml())

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, `# A comment for the service
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
`, string(sections[0]))
		assert.Equal(t, `# a comment for the deployment
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
`, string(sections[1]))
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

func simpleTestRoleBinding() []byte {
	return []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8s-dogu-operator-leader-election-rolebinding
  namespace: ecosystem
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8s-dogu-operator-leader-election-role
subjects:
- kind: ServiceAccount
  name: k8s-dogu-operator-controller-manager
  namespace: ecosystem
`)
}

func twoNamespacedResources() []byte {
	return []byte(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-dogu-operator-controller-manager
  namespace: ecosystem
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8s-dogu-operator-leader-election-role
  namespace: ecosystem
rules:
- apiGroups:
`)
}
