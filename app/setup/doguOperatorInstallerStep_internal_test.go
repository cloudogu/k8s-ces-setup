package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewDoguOperatorInstallerStep(t *testing.T) {
	t.Parallel()

	// given

	// when
	actual := newDoguOperatorInstallerStep(nil, "https://", "1.2.3")

	// then
	assert.NotNil(t, actual)
	require.Implements(t, (*ExecutorStep)(nil), actual)
}

func TestDoguOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	creator := newDoguOperatorInstallerStep(nil, "https://", "1.2.3")

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Install dogu operator version 1.2.3", description)
}

func TestDoguOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("Setup step fails for because...", func(t *testing.T) {
		// given

		creator := newDoguOperatorInstallerStep(nil, "http://", "1.2.3")

		// when
		err := creator.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Equal(t, "cannot create namespace 1.2.3 with clientset: namespaces \"1.2.3\" already exists", err.Error())
	})

	t.Run("Setup step runs without any problems", func(t *testing.T) {
		// given
		clientSetMock := testclient.NewSimpleClientset()
		creator := newDoguOperatorInstallerStep(nil, "http://", "1.2.3")

		// when
		err := creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		retrievedNamespace, err := clientSetMock.CoreV1().Namespaces().Get(context.Background(), "1.2.3", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "1.2.3", retrievedNamespace.GetName())
	})
}
