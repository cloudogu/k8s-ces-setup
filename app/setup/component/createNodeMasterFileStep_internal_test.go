package component

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

var testCtx = context.Background()

func TestNewNodeMasterCreationStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		config := &rest.Config{}

		// when
		step, err := NewNodeMasterCreationStep(config, "namespace")

		// then
		require.NoError(t, err)
		assert.NotNil(t, step)
		assert.NotNil(t, step.ClientSet)
		assert.Equal(t, "namespace", step.TargetNamespace)
	})
}

func Test_nodeMasterCreationStep_GetStepDescription(t *testing.T) {
	t.Run("get description", func(t *testing.T) {
		// given
		nodeMasterCreationStep := nodeMasterCreationStep{}

		// when
		description := nodeMasterCreationStep.GetStepDescription()

		// then
		assert.Equal(t, "Setup node master file", description)
	})
}

func Test_nodeMasterCreationStep_PerformSetupStep(t *testing.T) {
	t.Run("config map is begin created", func(t *testing.T) {
		// given
		nodeMasterCreationStep := nodeMasterCreationStep{}
		nodeMasterCreationStep.TargetNamespace = "my-namespace"
		nodeMasterCreationStep.ClientSet = fake.NewSimpleClientset()

		// when
		err := nodeMasterCreationStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)

		cm, err := nodeMasterCreationStep.ClientSet.CoreV1().ConfigMaps("my-namespace").Get(testCtx, "node-master-file", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "etcd.my-namespace.svc.cluster.local", cm.Data["node_master"])
	})

	t.Run("config map is begin updated", func(t *testing.T) {
		// given
		existingCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node-master-file",
				Namespace: "my-namespace",
			},
			Data: map[string]string{
				"node_master": "epic.etcd.at.here.local",
			},
		}
		nodeMasterCreationStep := nodeMasterCreationStep{}
		nodeMasterCreationStep.TargetNamespace = "my-namespace"
		nodeMasterCreationStep.ClientSet = fake.NewSimpleClientset(existingCM)

		// when
		err := nodeMasterCreationStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)

		cm, err := nodeMasterCreationStep.ClientSet.CoreV1().ConfigMaps("my-namespace").Get(testCtx, "node-master-file", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "etcd.my-namespace.svc.cluster.local", cm.Data["node_master"])
	})
}
