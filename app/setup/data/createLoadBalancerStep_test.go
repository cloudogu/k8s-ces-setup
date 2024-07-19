package data

import (
	"context"
	appctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var testCtx = context.Background()

func Test_createLoadBalancerStep_PerformSetupStep(t *testing.T) {
	t.Run("creates load-balancer if none exists", func(t *testing.T) {
		// given
		fakeClient := fake.NewSimpleClientset()
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace, testLoadbalancerName)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		actual, err := fakeClient.CoreV1().Services(testNamespace).Get(testCtx, testLoadbalancerName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotNil(t, actual)
		expected := map[string]string{"app": "ces"}
		assert.Equal(t, expected, actual.ObjectMeta.Labels)
	})
	t.Run("deletes misconfigured but still existing service and creates a new one", func(t *testing.T) {
		// given
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{

				Name:   testLoadbalancerName,
				Labels: map[string]string{"delete": "me"},
			},
			Spec: corev1.ServiceSpec{Type: corev1.ServiceTypeClusterIP},
		}
		fakeClient := fake.NewSimpleClientset(serviceResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace, testLoadbalancerName)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		actual, err := fakeClient.CoreV1().Services(testNamespace).Get(testCtx, testLoadbalancerName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotNil(t, actual)
		expected := map[string]string{"app": "ces"}
		assert.Equal(t, expected, actual.ObjectMeta.Labels)
		assert.Equal(t, corev1.ServiceTypeLoadBalancer, actual.Spec.Type)
	})
}
