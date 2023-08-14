package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	remoteMocks "github.com/cloudogu/cesapp-lib/remote/mocks"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
)

var testCtx = context.Background()

func getSetupCtx() appcontext.SetupContext {
	return appcontext.SetupContext{
		AppConfig: &appcontext.Config{
			TargetNamespace: "mynamespace",
		},
		SetupJsonConfiguration: &appcontext.SetupJsonConfiguration{},
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
		appCtx := getSetupCtx()
		registryMock := &remoteMocks.Registry{}
		step := NewValidatorStep(registryMock, &appCtx)
		step.setupJsonValidator = validatorMock

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
