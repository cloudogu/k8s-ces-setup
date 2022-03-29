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

func Test_newSecretCreator(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()

	// when
	creator := newSecretCreator(clientSetMock, "namespace", "currentNamespace")

	// then
	assert.NotNil(t, creator)
}

func Test_secretCreator_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	clientSetMock := testclient.NewSimpleClientset()
	creator := newSecretCreator(clientSetMock, "myTestNamespace", "currentNamespace")

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Create new Secrets in namespace myTestNamespace", description)
}

func Test_secretCreator_PerformSetupStep(t *testing.T) {
	t.Run("Setup step runs without any problems", func(t *testing.T) {
		// given
		clientSetMock := testclient.NewSimpleClientset()
		fmt.Println(dccSecret)
		dccSecret, _ = clientSetMock.CoreV1().Secrets("currentNamespace").Create(context.TODO(), dccSecret, metav1.CreateOptions{})
		_, err := clientSetMock.CoreV1().Secrets("currentNamespace").Create(context.TODO(), dockerSecret, metav1.CreateOptions{})
		require.NoError(t, err)

		creator := newSecretCreator(clientSetMock, "myTestNamespace", "currentNamespace")

		// when
		err = creator.PerformSetupStep()

		// then
		require.NoError(t, err)

		retrievedDccSecret, err := clientSetMock.CoreV1().Secrets(creator.TargetNamespace).Get(context.Background(), "dogu-cloudogu-com", metav1.GetOptions{})
		require.NoError(t, err)
		retrievedDockerSecret, err := clientSetMock.CoreV1().Secrets(creator.TargetNamespace).Get(context.Background(), "registry-cloudogu-com", metav1.GetOptions{})
		require.NoError(t, err)

		assert.Equal(t, "dogu-cloudogu-com", retrievedDccSecret.GetName())
		assert.Equal(t, dccSecret.Data, retrievedDccSecret.Data)
		assert.Equal(t, dccSecret.StringData, retrievedDccSecret.StringData)
		assert.Equal(t, "registry-cloudogu-com", retrievedDockerSecret.GetName())
		assert.Equal(t, dockerSecret.Data, retrievedDockerSecret.Data)
		assert.Equal(t, dockerSecret.StringData, retrievedDockerSecret.StringData)
	})
}
