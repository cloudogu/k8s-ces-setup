package setup

import (
	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var doguOperatorSetupCtx = ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: ctx.Config{
		TargetNamespace: testTargetNamespaceName,
		DoguOperatorURL: "http://url.server.com/dogu/operator.yaml",
	},
}

func TestNewDoguOperatorInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual := newDoguOperatorInstallerStep(nil, doguOperatorSetupCtx)

	// then
	assert.NotNil(t, actual)
	require.Implements(t, (*ExecutorStep)(nil), actual)
}

func TestDoguOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	creator := newDoguOperatorInstallerStep(nil, doguOperatorSetupCtx)

	// when
	description := creator.GetStepDescription()

	// then
	assert.Equal(t, "Install dogu operator from http://url.server.com/dogu/operator.yaml", description)
}

func TestDoguOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

}
