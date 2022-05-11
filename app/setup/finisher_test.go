package setup_test

import (
	context2 "context"
	"testing"

	v13 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"

	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/setup"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewFinisher(t *testing.T) {
	t.Run("finish setup with no resources", func(t *testing.T) {
		// given
		fakeclient := fake.NewSimpleClientset()
		finisher := setup.NewFinisher(fakeclient, "mytestnamespace")

		// when
		err := finisher.FinishSetup()

		// then
		require.NoError(t, err)

		cm, err := context.GetSetupConfigMap(fakeclient, "mytestnamespace")
		require.NoError(t, err)

		assert.Equal(t, "installed", cm.Data[context.SetupStateKey])
	})

	t.Run("finish setup with resources", func(t *testing.T) {
		// given
		n := "mytestnamespace"
		crb1 := &v1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup-cluster-resources", Namespace: n}}
		crb2 := &v1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup-cluster-non-resources", Namespace: n}}
		cr1 := &v1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup-cluster-resources", Namespace: n}}
		cr2 := &v1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup-cluster-non-resources", Namespace: n}}
		rb := &v1.RoleBinding{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup", Namespace: n}}
		r := &v1.Role{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup", Namespace: n}}
		sa := &v12.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup", Namespace: n}}
		s := &v12.Service{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup", Namespace: n}}
		d := &v13.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "k8s-ces-setup", Namespace: n}}

		fakeclient := fake.NewSimpleClientset(crb1, crb2, cr1, cr2, rb, r, sa, s, d)
		finisher := setup.NewFinisher(fakeclient, n)

		// when
		err := finisher.FinishSetup()

		// then
		require.NoError(t, err)

		cm, err := context.GetSetupConfigMap(fakeclient, n)
		require.NoError(t, err)
		assert.Equal(t, "installed", cm.Data[context.SetupStateKey])

		_, err = fakeclient.RbacV1().ClusterRoleBindings().Get(context2.Background(), "k8s-ces-setup-cluster-resources", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.RbacV1().ClusterRoleBindings().Get(context2.Background(), "k8s-ces-setup-cluster-non-resources", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.RbacV1().ClusterRoles().Get(context2.Background(), "k8s-ces-setup-cluster-resources", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.RbacV1().ClusterRoles().Get(context2.Background(), "k8s-ces-setup-cluster-non-resources", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.RbacV1().RoleBindings(n).Get(context2.Background(), "k8s-ces-setup", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.RbacV1().Roles(n).Get(context2.Background(), "k8s-ces-setup", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.CoreV1().ServiceAccounts(n).Get(context2.Background(), "k8s-ces-setup", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.CoreV1().Services(n).Get(context2.Background(), "k8s-ces-setup", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
		_, err = fakeclient.AppsV1().Deployments(n).Get(context2.Background(), "k8s-ces-setup", metav1.GetOptions{})
		assert.True(t, errors.IsNotFound(err))
	})
}

func TestFinisher_FinishSetup(t *testing.T) {
	t.Run("fail to ", func(t *testing.T) {
		// given

		// when

		// then
	})
}
