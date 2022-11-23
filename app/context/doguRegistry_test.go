package context

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestReadDoguRegistrySecretFromClusterWithIndexSchema(t *testing.T) {
	// given
	doguRegistrySecret := &DoguRegistrySecret{
		Endpoint:  "testendpoint",
		Username:  "testuser",
		Password:  "testpassword",
		URLSchema: "index",
	}
	dataMap := make(map[string][]byte)
	dataMap["endpoint"] = []byte("testendpoint")
	dataMap["username"] = []byte("testuser")
	dataMap["password"] = []byte("testpassword")
	dataMap["urlschema"] = []byte("index")
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator-dogu-registry", Namespace: "ecosystem"},
		Immutable:  nil,
		Data:       dataMap,
	}
	clientset := fake.NewSimpleClientset(secret)

	// when
	testSecret, err := ReadDoguRegistrySecretFromCluster(clientset, "ecosystem")

	// then
	require.NoError(t, err)
	assert.Equal(t, doguRegistrySecret, testSecret)
}

func TestReadDoguRegistrySecretFromClusterWithDefaultSchema(t *testing.T) {
	// given
	doguRegistrySecret := &DoguRegistrySecret{
		Endpoint:  "testendpoint",
		Username:  "testuser",
		Password:  "testpassword",
		URLSchema: "default",
	}
	dataMap := make(map[string][]byte)
	dataMap["endpoint"] = []byte("testendpoint")
	dataMap["username"] = []byte("testuser")
	dataMap["password"] = []byte("testpassword")
	dataMap["urlschema"] = []byte("default")
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator-dogu-registry", Namespace: "ecosystem"},
		Immutable:  nil,
		Data:       dataMap,
	}
	clientset := fake.NewSimpleClientset(secret)

	// when
	testSecret, err := ReadDoguRegistrySecretFromCluster(clientset, "ecosystem")

	// then
	require.NoError(t, err)
	assert.Equal(t, doguRegistrySecret, testSecret)
}

func TestReadDoguRegistrySecretFromClusterWithEmptySchema(t *testing.T) {
	// given
	doguRegistrySecret := &DoguRegistrySecret{
		Endpoint:  "testendpoint",
		Username:  "testuser",
		Password:  "testpassword",
		URLSchema: "default",
	}
	dataMap := make(map[string][]byte)
	dataMap["endpoint"] = []byte("testendpoint")
	dataMap["username"] = []byte("testuser")
	dataMap["password"] = []byte("testpassword")
	dataMap["urlschema"] = []byte("")
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator-dogu-registry", Namespace: "ecosystem"},
		Immutable:  nil,
		Data:       dataMap,
	}
	clientset := fake.NewSimpleClientset(secret)

	// when
	testSecret, err := ReadDoguRegistrySecretFromCluster(clientset, "ecosystem")

	// then
	require.NoError(t, err)
	assert.Equal(t, doguRegistrySecret, testSecret)
}

func TestReadDoguRegistrySecretFromClusterWithWrongSecretName(t *testing.T) {
	// given
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "NoTaVaLiDnAmE", Namespace: "ecosystem"},
	}
	clientset := fake.NewSimpleClientset(secret)

	// when
	_, err := ReadDoguRegistrySecretFromCluster(clientset, "ecosystem")

	// then
	require.Error(t, err)
}
