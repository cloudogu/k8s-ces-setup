package component

import (
	"context"
	"strings"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	v1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

var etcdClientSetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: ctx.Config{
		TargetNamespace:     testTargetNamespaceName,
		EtcdClientImageRepo: "registryurl/registryname/repo:tag",
	},
}

func TestNewEtcdClientInstallerStep(t *testing.T) {

	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()

	// when
	creator := NewEtcdClientInstallerStep(clientSetMock, etcdClientSetupCtx)

	// then
	assert.NotNil(t, creator)
}

func TestEtcdClientInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()
	creator := NewEtcdClientInstallerStep(clientSetMock, etcdClientSetupCtx)

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Install etcd client from registryurl/registryname/repo:tag", description)
}

func TestEtcdClientInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("Setup step runs without any problems", func(t *testing.T) {
		// given
		clientSetMock := testclient.NewSimpleClientset()
		creator := NewEtcdClientInstallerStep(clientSetMock, etcdClientSetupCtx)

		// when
		err := creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		actual, err := clientSetMock.CoreV1().Pods(testTargetNamespaceName).Get(context.Background(), "etcd-client", metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, "etcd-client", actual.GetName())
		require.Len(t, actual.Spec.Containers, 1)
		assert.Equal(t, "registryurl/registryname/repo:tag", actual.Spec.Containers[0].Image)
		assert.Equal(t, "sleep infinity", strings.Join(actual.Spec.Containers[0].Command, " "))
		assert.Contains(t, actual.Spec.Containers[0].Env, v1.EnvVar{Name: "ETCDCTL_API", Value: "2"})
		assert.Contains(t, actual.Spec.Containers[0].Env, v1.EnvVar{Name: "ETCDCTL_ENDPOINTS", Value: "http://etcd.myfavouritenamespace-1.svc.cluster.local:4001"})
	})
}
