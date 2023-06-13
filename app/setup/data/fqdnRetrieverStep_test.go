package data

import (
	gocontext "context"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func Test_fqdnRetrieverStep_PerformSetupStep(t *testing.T) {
	namespace := "ecosystem"
	cesLoadbalancerName := "ces-loadbalancer"
	t.Run("successfully set FQDN", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{Fqdn: ""}, Dogus: context.Dogus{Install: []string{"nginx-ingress"}}}
		fakeClient := fake.NewSimpleClientset()

		step := NewFQDNRetrieverStep(config, fakeClient, namespace)

		// when
		timer := time.NewTimer(time.Second * 2)
		go func() {
			<-timer.C
			patch := []byte(`{"status":{"loadBalancer":{"ingress":[{"ip": "555.444.333.222"}]}}}`)
			service, err := fakeClient.CoreV1().Services(namespace).Patch(
				gocontext.Background(),
				cesLoadbalancerName,
				types.MergePatchType,
				patch,
				metav1.PatchOptions{},
			)
			require.NoError(t, err)
			assert.NotNil(t, service)
		}()
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		assert.Equal(t, "555.444.333.222", config.Naming.Fqdn)
	})
	t.Run("failed due to missing nginx-ingress", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{Fqdn: ""}}
		fakeClient := fake.NewSimpleClientset()

		step := NewFQDNRetrieverStep(config, fakeClient, namespace)

		// when
		err := step.PerformSetupStep()

		// then
		require.ErrorContains(t, err, "invalid configuration. FQDN can only be created with nginx-ingress installed")
	})
	t.Run("failure when service with name already exists", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{Fqdn: ""}, Dogus: context.Dogus{Install: []string{"nginx-ingress"}}}
		fakeClient := fake.NewSimpleClientset()
		serviceResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: cesLoadbalancerName,
			},
		}
		_, err := fakeClient.CoreV1().Services(namespace).Create(gocontext.Background(), serviceResource, metav1.CreateOptions{})
		if err != nil {
			assert.Fail(t, "Test-service could not be created")
		}
		step := NewFQDNRetrieverStep(config, fakeClient, namespace)

		// when
		err = step.PerformSetupStep()

		// then
		require.ErrorContains(t, err, `services "ces-loadbalancer" already exists`)
	})
}

func TestNewFQDNRetrieverStep(t *testing.T) {
	// when
	step := NewFQDNRetrieverStep(nil, nil, "")

	// then
	require.NotNil(t, step)
}

func Test_fqdnRetrieverStep_GetStepDescription(t *testing.T) {
	// given
	step := NewFQDNRetrieverStep(nil, nil, "")

	// when
	description := step.GetStepDescription()

	// then
	assert.Equal(t, "Retrieving a new FQDN from the IP of a loadbalancer service", description)
}
