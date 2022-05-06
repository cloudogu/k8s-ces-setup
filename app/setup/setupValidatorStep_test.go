package setup_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/k8s-ces-setup/app/setup"

	"github.com/cloudogu/k8s-ces-setup/app/setup/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func createTestMocks(withSecret bool) (*fake.Clientset, *context.SetupContext) {
	setupCtx := &context.SetupContext{
		AppConfig: context.Config{
			TargetNamespace: "mynamespace",
		},
		StartupConfiguration: context.SetupConfiguration{},
	}

	if withSecret {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      context.SecretDoguRegistry,
				Namespace: "mynamespace",
			},
			StringData: map[string]string{"username": "myuser", "password": "mypass", "endpoint": "myendpoint"},
		}

		return fake.NewSimpleClientset(secret), setupCtx
	} else {
		return fake.NewSimpleClientset(), setupCtx
	}
}

func TestNewValidatorStep(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		// given
		clientMock, ctx := createTestMocks(true)

		// when
		step, err := setup.NewValidatorStep(clientMock, ctx)

		// then
		require.NoError(t, err)
		require.NotNil(t, step)
	})

	t.Run("failed to get secret", func(t *testing.T) {
		// given
		clientMock, setupCtx := createTestMocks(false)

		// when
		_, err := setup.NewValidatorStep(clientMock, setupCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get secret")
	})

}

func Test_setupValidatorStep_GetStepDescription(t *testing.T) {
	t.Run("get correct description", func(t *testing.T) {
		// given
		clientMock, ctx := createTestMocks(true)
		step, err := setup.NewValidatorStep(clientMock, ctx)
		require.NoError(t, err)

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, "Validating the setup configuration", description)
	})
}

func Test_setupValidatorStep_PerformSetupStep(t *testing.T) {
	t.Run("sucessful performing step", func(t *testing.T) {
		// given
		validatorMock := &mocks.ConfigurationValidator{}
		validatorMock.On("ValidateConfiguration", mock.Anything).Return(nil)
		client, setupContext := createTestMocks(true)
		step, err := setup.NewValidatorStep(client, setupContext)
		require.NoError(t, err)
		step.Validator = validatorMock

		// when
		err = step.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
