package data

import (
	"context"
	"fmt"
	appctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var testCtx = context.Background()

func Test_createLoadBalancerStep_PerformSetupStep(t *testing.T) {
	t.Run("creates load-balancer if none exists", func(t *testing.T) {
		// given
		fakeClient := fake.NewClientset()
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		actual, err := fakeClient.CoreV1().Services(testNamespace).Get(testCtx, "ces-loadbalancer", metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotNil(t, actual)
		expected := map[string]string{"app": "ces"}
		assert.Equal(t, expected, actual.ObjectMeta.Labels)
	})
	t.Run("deletes misconfigured but still existing service and creates a new one", func(t *testing.T) {
		// given
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{

				Name:   cesLoadbalancerName,
				Labels: map[string]string{"delete": "me"},
			},
			Spec: corev1.ServiceSpec{Type: corev1.ServiceTypeClusterIP},
		}
		fakeClient := fake.NewClientset(serviceResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		actual, err := fakeClient.CoreV1().Services(testNamespace).Get(testCtx, "ces-loadbalancer", metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotNil(t, actual)
		expected := map[string]string{"app": "ces"}
		assert.Equal(t, expected, actual.ObjectMeta.Labels)
		assert.Equal(t, corev1.ServiceTypeLoadBalancer, actual.Spec.Type)
	})

	t.Run("reuse existing loadbalancer with initial port mapping to keep current ip", func(t *testing.T) {
		// given
		expectedLabels := map[string]string{"app": "ces", "actual-ip": "1.2.3.4"}
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{

				Name:      cesLoadbalancerName,
				Labels:    expectedLabels,
				Namespace: testNamespace,
			},
			Spec: corev1.ServiceSpec{Type: corev1.ServiceTypeClusterIP},
		}
		fakeClient := fake.NewClientset(serviceResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}
		sut := NewCreateLoadBalancerStep(config, fakeClient, testNamespace)

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		actual, err := fakeClient.CoreV1().Services(testNamespace).Get(testCtx, "ces-loadbalancer", metav1.GetOptions{})
		require.NoError(t, err)
		assert.NotNil(t, actual)
		require.Equal(t, expectedLabels, actual.Labels)
		assert.Equal(t, corev1.ServiceTypeLoadBalancer, actual.Spec.Type)
		expectedIPSingleStackPolicy := corev1.IPFamilyPolicySingleStack
		assert.Equal(t, &expectedIPSingleStackPolicy, actual.Spec.IPFamilyPolicy)
		assert.Equal(t, []corev1.IPFamily{corev1.IPv4Protocol}, actual.Spec.IPFamilies)
		assert.Equal(t, map[string]string{DoguLabelName: nginxIngressName}, actual.Spec.Selector)

		expectedServicePorts := []corev1.ServicePort{
			{
				Name:       fmt.Sprintf("%s-%d", nginxIngressName, 80),
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(80),
				TargetPort: intstr.FromInt32(80),
			},
			{
				Name:       fmt.Sprintf("%s-%d", nginxIngressName, 443),
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(443),
				TargetPort: intstr.FromInt32(443),
			},
		}
		assert.Equal(t, expectedServicePorts, actual.Spec.Ports)
	})
}
