package component

import (
	"github.com/cloudogu/k8s-apply-lib/apply"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testTargetNamespaceName = "myfavouritenamespace-1"
const doguOperatorURL = "http://url.server.com/dogu/operator.yaml"

var doguOperatorSetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: &ctx.Config{
		TargetNamespace: testTargetNamespaceName,
		DoguOperatorURL: doguOperatorURL,
	},
}

func TestNewDoguOperatorInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual, _ := NewDoguOperatorInstallerStep(doguOperatorSetupCtx, &mockK8sClient{})

	// then
	assert.NotNil(t, actual)
}

func TestDoguOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewDoguOperatorInstallerStep(doguOperatorSetupCtx, &mockK8sClient{})

	// when
	description := installer.GetStepDescription()

	// then
	assert.Equal(t, "Install dogu operator from http://url.server.com/dogu/operator.yaml", description)
}

func TestDoguOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("should perform an installation without resource modification", func(t *testing.T) {
		// given
		var doguOpYamlBytes apply.YamlDocument = []byte("yaml result goes here")
		mockedResourceRegistryClient := newMockResourceRegistryClient(t)
		mockedResourceRegistryClient.EXPECT().GetResourceFileContent(doguOperatorURL).Return(doguOpYamlBytes, nil)
		mockedK8sClient := newMockK8sClient(t)
		mockedK8sClient.EXPECT().ApplyWithOwner(doguOpYamlBytes, testTargetNamespaceName, mock.Anything).Return(nil)

		installer := doguOperatorInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            doguOperatorURL,
			resourceRegistryClient: mockedResourceRegistryClient,
			k8sClient:              mockedK8sClient,
		}

		// when
		err := installer.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
