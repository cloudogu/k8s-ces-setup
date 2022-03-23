package setup

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testNamespaceName = "myfavouritenamespace-1"

func Test_replaceNamespace(t *testing.T) {
	t.Run("should replace namespace", func(t *testing.T) {
		input := simpleTestRoleBinding()
		// when
		actual := replaceNamespace(input, testNamespaceName)

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
