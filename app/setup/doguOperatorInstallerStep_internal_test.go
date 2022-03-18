package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewDoguOperatorInstallerStep(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()

	// when
	actual := newDoguOperatorInstallerStep(clientSetMock, "https://", "1.2.3")

	// then
	assert.NotNil(t, actual)
	require.Implements(t, (*ExecutorStep)(nil), actual)
}

func TestDoguOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()
	creator := newDoguOperatorInstallerStep(clientSetMock, "https://", "1.2.3")

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Install dogu operator version 1.2.3", description)
}

func TestDoguOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("Setup step fails for because...", func(t *testing.T) {
		// given
		namespace := v1.Namespace{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: "myNamespace"},
			Spec:       v1.NamespaceSpec{},
			Status:     v1.NamespaceStatus{},
		}

		clientSetMock := testclient.NewSimpleClientset(&namespace)
		creator := newDoguOperatorInstallerStep(clientSetMock, "http://", "1.2.3")

		// when
		err := creator.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Equal(t, "cannot create namespace 1.2.3 with clientset: namespaces \"1.2.3\" already exists", err.Error())
	})

	t.Run("Setup step runs without any problems", func(t *testing.T) {
		// given
		clientSetMock := testclient.NewSimpleClientset()
		creator := newDoguOperatorInstallerStep(clientSetMock, "http://", "1.2.3")

		// when
		err := creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		retrievedNamespace, err := clientSetMock.CoreV1().Namespaces().Get(context.Background(), "1.2.3", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "1.2.3", retrievedNamespace.GetName())
	})
}
