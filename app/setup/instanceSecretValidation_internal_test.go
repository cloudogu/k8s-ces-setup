package setup

import (
	"testing"

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
	creator := newInstanceSecretValidator(clientSetMock, "namespace")

	// then
	assert.NotNil(t, creator)
}

func TestNamespaceCreator_validate(t *testing.T) {
	t.Parallel()

	t.Run("fails for missing dogu secret", func(t *testing.T) {
		// given
		namespace := v1.Namespace{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: testTargetNamespaceName},
			Spec:       v1.NamespaceSpec{},
			Status:     v1.NamespaceStatus{},
		}
		clientSetMock := testclient.NewSimpleClientset(&namespace)
		creator := newInstanceSecretValidator(clientSetMock, testTargetNamespaceName)

		// when
		err := creator.validate()

		// then
		require.Error(t, err)
		assert.Equal(t, "instance secret validation error: cannot read secret from target namespace myfavouritenamespace-1: secrets \"k8s-dogu-operator-dogu-registry\" not found", err.Error())
	})

	t.Run("fails for missing image secret", func(t *testing.T) {
		// given
		doguSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretNameDoguRegistry,
				Namespace: testTargetNamespaceName,
			},
		}
		clientSetMock := testclient.NewSimpleClientset(doguSecret)
		creator := newInstanceSecretValidator(clientSetMock, testTargetNamespaceName)

		// when
		err := creator.validate()

		// then
		require.Error(t, err)
		assert.Equal(t, "instance secret validation error: cannot read secret from target namespace myfavouritenamespace-1: secrets \"k8s-dogu-operator-docker-registry\" not found", err.Error())
	})

	t.Run("runs without any problems", func(t *testing.T) {
		// given
		doguSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretNameDoguRegistry,
				Namespace: testTargetNamespaceName,
			},
		}
		imageSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretNameImageRegistry,
				Namespace: testTargetNamespaceName,
			},
			Type: v1.SecretTypeDockercfg,
		}
		clientSetMock := testclient.NewSimpleClientset(doguSecret, imageSecret)
		sut := newInstanceSecretValidator(clientSetMock, testTargetNamespaceName)

		// when
		err := sut.validate()

		// then
		require.NoError(t, err)
	})
}
