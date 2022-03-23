package setup

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testNamespaceName = "myfavouritenamespace-1"

func Test_replaceNamespace(t *testing.T) {
	t.Run("should replace namespace within simple namespaced resource", func(t *testing.T) {
		input := simpleTestRoleBinding()
		// when
		actual := replaceNamespacedResources(input, testNamespaceName)

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
`, testNamespaceName, testNamespaceName)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("should replace namespace in complex resources", func(t *testing.T) {
		input := twoNamespacedResources()

		// when
		actual := replaceNamespacedResources(input, testNamespaceName)

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
			testNamespaceName, testNamespaceName)

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
