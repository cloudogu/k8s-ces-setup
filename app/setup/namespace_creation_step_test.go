package setup_test

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNewNamespaceCreator(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()

	// when
	creator := setup.NewNamespaceCreator(clientSetMock, "namespace")

	// then
	assert.NotNil(t, creator)
}

func TestNamespaceCreator_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()
	creator := setup.NewNamespaceCreator(clientSetMock, "myTestNamespace")

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Create new namespace myTestNamespace", description)
}

func TestNamespaceCreator_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("Setup step fails for creating already existing namespace", func(t *testing.T) {
		// given
		namespace := v1.Namespace{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: "myTestNamespace"},
			Spec:       v1.NamespaceSpec{},
			Status:     v1.NamespaceStatus{},
		}
		clientSetMock := testclient.NewSimpleClientset(&namespace)
		creator := setup.NewNamespaceCreator(clientSetMock, "myTestNamespace")

		// when
		err := creator.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Equal(t, "cannot create namespace myTestNamespace with clientset; namespaces \"myTestNamespace\" already exists", err.Error())
	})

	t.Run("Setup step runs without any problems", func(t *testing.T) {
		// given
		clientSetMock := testclient.NewSimpleClientset()
		creator := setup.NewNamespaceCreator(clientSetMock, "myTestNamespace")

		// when
		err := creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		retrievedNamespace, err := clientSetMock.CoreV1().Namespaces().Get(context.Background(), "myTestNamespace", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "myTestNamespace", retrievedNamespace.GetName())
	})
}
