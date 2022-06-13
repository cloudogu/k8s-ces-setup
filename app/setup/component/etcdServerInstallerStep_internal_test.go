package component

import (
	"github.com/cloudogu/k8s-apply-lib/apply"
	"github.com/stretchr/testify/mock"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const etcdServerResourceURL = "https://url.server.com/etcd/resource.yaml"

var etcdServerSetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: &ctx.Config{
		TargetNamespace:       testTargetNamespaceName,
		EtcdServerResourceURL: etcdServerResourceURL,
	},
}

func TestNewEtcdServerInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual, _ := NewEtcdServerInstallerStep(etcdServerSetupCtx, &mockK8sClient{})

	// then
	assert.NotNil(t, actual)
}

func TestEtcdServerInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewEtcdServerInstallerStep(etcdServerSetupCtx, &mockK8sClient{})

	// when
	description := installer.GetStepDescription()

	// then
	assert.Equal(t, "Install etcd server from https://url.server.com/etcd/resource.yaml", description)
}

func TestEtcdServerInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("should perform an installation without resource modification", func(t *testing.T) {
		// given
		var yamlBytes apply.YamlDocument = []byte("yaml result goes here")

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", etcdServerResourceURL).Return([]byte(yamlBytes), nil)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("ApplyWithOwner", yamlBytes, testTargetNamespaceName, mock.Anything).Return(nil)

		installer := etcdServerInstallerStep{
			namespace:   testTargetNamespaceName,
			resourceURL: etcdServerResourceURL,
			fileClient:  mockedFileClient,
			k8sClient:   mockedK8sClient,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mockedFileClient.AssertExpectations(t)
		mockedK8sClient.AssertExpectations(t)
	})
}
