package data

import (
	"context"
	appctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

const (
	testNamespace        = "ecosystem"
	testLoadbalancerName = "ces-test-loadbalancer"
)

func Test_fqdnRetrieverStep_PerformSetupStep(t *testing.T) {
	t.Run("should successfully set FQDN to IP when the load-balancer receives an external IP address", func(t *testing.T) {
		// given
		mockedLoadBalancerResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testLoadbalancerName,
				Namespace: testNamespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{DoguLabelName: nginxIngressName},
			},
		}
		fakeClient := fake.NewSimpleClientset(mockedLoadBalancerResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}

		sut := NewFQDNRetrieverStep(config, fakeClient, testNamespace, testLoadbalancerName)

		// simulate asynchronous IP setting by cluster provider
		timer := time.NewTimer(time.Second * 2)
		go func() {
			<-timer.C
			patch := []byte(`{"status":{"loadBalancer":{"ingress":[{"ip": "111.222.111.222"}]}}}`)
			service, err := fakeClient.CoreV1().Services(testNamespace).Patch(
				context.Background(),
				testLoadbalancerName,
				types.MergePatchType,
				patch,
				metav1.PatchOptions{},
			)
			require.NoError(t, err)
			assert.NotNil(t, service)
		}()

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, "111.222.111.222", config.Naming.Fqdn)
	})
}
func TestCreateLoadBalancerStep_PerformSetupStep(t *testing.T) {
	t.Run("failed due to missing nginx-ingress", func(t *testing.T) {
		// given
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}}

		step := NewCreateLoadBalancerStep(config, nil, testNamespace, testLoadbalancerName)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid configuration: FQDN can only be created if nginx-ingress will be installed")
	})
}

func TestNewFQDNRetrieverStep(t *testing.T) {
	// when
	step := NewFQDNRetrieverStep(nil, nil, "", "")

	// then
	require.NotNil(t, step)
}

func Test_fqdnRetrieverStep_GetStepDescription(t *testing.T) {
	// given
	step := NewFQDNRetrieverStep(nil, nil, "", "")

	// when
	description := step.GetStepDescription()

	// then
	assert.Equal(t, "Retrieving a new FQDN from the IP of a loadbalancer service", description)
}
