package setup

import (
	"testing"

	"github.com/stretchr/testify/mock"

	remoteMocks "github.com/cloudogu/cesapp-lib/remote/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/require"
)

func getSetupCtx() context.SetupContext {
	return context.SetupContext{
		AppConfig: &context.Config{
			TargetNamespace: "mynamespace",
		},
		SetupJsonConfiguration: &context.SetupJsonConfiguration{},
	}
}

func TestNewValidatorStep(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}

		// when
		step := NewValidatorStep(registryMock, &ctx)

		// then
		require.NotNil(t, step)
	})
}

func Test_setupValidatorStep_GetStepDescription(t *testing.T) {
	t.Run("get correct description", func(t *testing.T) {
		// given
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}
		step := NewValidatorStep(registryMock, &ctx)

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, "Validating the setup configuration", description)
	})
}

func Test_setupValidatorStep_PerformSetupStep(t *testing.T) {
	t.Run("sucessful performing step", func(t *testing.T) {
		// given
		validatorMock := NewMockConfigurationValidator(t)
		validatorMock.EXPECT().ValidateConfiguration(mock.Anything).Return(nil)
		ctx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}
		step := NewValidatorStep(registryMock, &ctx)
		step.setupJsonValidator = validatorMock

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
