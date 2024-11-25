package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)

		// when
		step := NewValidatorStep(remoteDoguRepo, &ctx)

		// then
		require.NotNil(t, step)
	})
}

func Test_setupValidatorStep_GetStepDescription(t *testing.T) {
	t.Run("get correct description", func(t *testing.T) {
		// given
		ctx := getSetupCtx()
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		step := NewValidatorStep(remoteDoguRepo, &ctx)

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, "Validating the setup configuration", description)
	})
}

func Test_setupValidatorStep_PerformSetupStep(t *testing.T) {
	t.Run("sucessful performing step", func(t *testing.T) {
		// given
		validatorMock := newMockSetupJsonConfigurationValidator(t)
		validatorMock.EXPECT().Validate(mock.Anything, mock.Anything).Return(nil)
		appCtx := getSetupCtx()
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		step := NewValidatorStep(remoteDoguRepo, &appCtx)
		step.setupJsonValidator = validatorMock

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
