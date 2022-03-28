package setup

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	testclient "k8s.io/client-go/kubernetes/fake"
	"testing"
)

//go:embed testdata/dccSecret.yaml
var dccSecretBytes []byte
var dccSecret = &v1.Secret{}

//go:embed testdata/dockerSecret.yaml
var dockerSecretBytes []byte
var dockerSecret = &v1.Secret{}

func init() {
	err := yaml.Unmarshal(dccSecretBytes, dccSecret)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(dockerSecretBytes, dockerSecret)
	if err != nil {
		panic(err)
	}
}

func Test_getEnvVar(t *testing.T) {
	t.Run("successfully query env var namespace", func(t *testing.T) {
		// given
		t.Setenv("NAMESPACE", "myTestNamespace")

		// when
		ns, err := getEnvVar("NAMESPACE")

		// then
		require.NoError(t, err)

		assert.Equal(t, "myTestNamespace", ns)
	})

	t.Run("failed to query env var namespace", func(t *testing.T) {
		// when
		_, err := getEnvVar("NAMESPACE")

		// then
		require.Error(t, err)
	})
}

func Test_newSecretCreator(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()

	// when
	creator := newSecretCreator(clientSetMock, "namespace")

	// then
	assert.NotNil(t, creator)
}

func Test_secretCreator_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()
	creator := newSecretCreator(clientSetMock, "myTestNamespace")

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Create new Secrets in namespace myTestNamespace", description)
}

func Test_secretCreator_PerformSetupStep(t *testing.T) {
	t.Run("Setup step runs without any problems", func(t *testing.T) {
		t.Setenv("NAMESPACE", "actualNamespace")
		// given
		clientSetMock := testclient.NewSimpleClientset()
		fmt.Println(dccSecret)
		dccSecret, _ = clientSetMock.CoreV1().Secrets("actualNamespace").Create(context.TODO(), dccSecret, metav1.CreateOptions{})
		_, err := clientSetMock.CoreV1().Secrets("actualNamespace").Create(context.TODO(), dockerSecret, metav1.CreateOptions{})
		require.NoError(t, err)

		creator := newSecretCreator(clientSetMock, "myTestNamespace")

		// when
		err = creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		retrievedDccSecret, err := clientSetMock.CoreV1().Secrets(creator.Namespace).Get(context.Background(), "dogu-cloudogu-com", metav1.GetOptions{})
		require.NoError(t, err)
		retrievedDockerSecret, err := clientSetMock.CoreV1().Secrets(creator.Namespace).Get(context.Background(), "docker-image-pull", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "dogu-cloudogu-com", retrievedDccSecret.GetName())
		assert.Equal(t, dccSecret.Data, retrievedDccSecret.Data)
		assert.Equal(t, dccSecret.StringData, retrievedDccSecret.StringData)
		assert.Equal(t, "docker-image-pull", retrievedDockerSecret.GetName())
		assert.Equal(t, dockerSecret.Data, retrievedDockerSecret.Data)
		assert.Equal(t, dockerSecret.StringData, retrievedDockerSecret.StringData)
	})
}
