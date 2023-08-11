package data

import (
	appctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func Test_createLoadBalancerStep_PerformSetupStep(t *testing.T) {
	t.Run("creates load-balancer if none exists", func(t *testing.T) {
		// given
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: cesLoadbalancerName,
			},
		}
		fakeClient := fake.NewSimpleClientset(serviceResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep()

		// then
		require.ErrorContains(t, err, `services "ces-loadbalancer" already exists`)
	})
	t.Run("deletes existing load-balancer and creates a new one", func(t *testing.T) {
		// given
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: cesLoadbalancerName,
			},
		}
		fakeClient := fake.NewSimpleClientset(serviceResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep()

		// then
		require.ErrorContains(t, err, `services "ces-loadbalancer" already exists`)
	})
}
