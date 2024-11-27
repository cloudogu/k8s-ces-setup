package data

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"testing"
)

func Test_disableDefaultSAAutomountStep_PerformSetupStep(t *testing.T) {
	t.Run("successfully disable default service account token automount", func(t *testing.T) {
		// given
		defaultSA := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: testNamespace,
			},
		}
		fakeClient := fake.NewClientset(defaultSA)
		sut := NewDisableDefaultSAAutomountStep(fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		// test default service account token automount deactivation
		sa, err := sut.clientSet.CoreV1().ServiceAccounts(testNamespace).Get(testCtx, "default", metav1.GetOptions{})
		require.NoError(t, err)
		assert.False(t, *sa.AutomountServiceAccountToken)
	})

	t.Run("failed because setup could not get default service account", func(t *testing.T) {
		// given
		fakeClient := fake.NewClientset()
		fakeClient.PrependReactor("get", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, assert.AnError
		})
		sut := NewDisableDefaultSAAutomountStep(fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get default service account")
	})

	t.Run("failed because setup could not update default service account", func(t *testing.T) {
		// given
		fakeClient := fake.NewClientset()
		fakeClient.PrependReactor("get", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, &corev1.ServiceAccount{}, nil
		})
		fakeClient.PrependReactor("update", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, assert.AnError
		})
		sut := NewDisableDefaultSAAutomountStep(fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to deactivate token automount on default service account")
	})
}
