package setup_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/k8s-ces-setup/app/setup"

	"github.com/cloudogu/k8s-ces-setup/app/setup/mocks"

	remoteMocks "github.com/cloudogu/cesapp-lib/remote/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/require"
)

func getSetupCtx() context.SetupContext {
	return context.SetupContext{
		AppConfig: context.Config{
			TargetNamespace: "mynamespace",
		},
		StartupConfiguration: context.SetupConfiguration{},
	}
}

func TestNewValidatorStep(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}

		// when
		step := setup.NewValidatorStep(registryMock, &ctx)

		// then
		require.NotNil(t, step)
	})
}

func Test_setupValidatorStep_GetStepDescription(t *testing.T) {
	t.Run("get correct description", func(t *testing.T) {
		// given
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}
		step := setup.NewValidatorStep(registryMock, &ctx)

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
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}
		step := setup.NewValidatorStep(registryMock, &ctx)
		step.Validator = validatorMock

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, validatorMock)
	})
}
