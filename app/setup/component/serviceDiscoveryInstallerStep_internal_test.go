package component

import (
	"github.com/cloudogu/k8s-apply-lib/apply"
	"testing"

	"github.com/stretchr/testify/mock"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const serviceDiscoveryResourceURL = "https://url.server.com/service-discovery/resource.yaml"

var serviceDiscoverySetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: &ctx.Config{
		TargetNamespace:     testTargetNamespaceName,
		ServiceDiscoveryURL: serviceDiscoveryResourceURL,
	},
}

func TestNewServiceDiscoveryInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual, _ := NewServiceDiscoveryInstallerStep(serviceDiscoverySetupCtx, &mockK8sClient{})

	// then
	assert.NotNil(t, actual)
}

func TestServiceDiscoveryInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewServiceDiscoveryInstallerStep(serviceDiscoverySetupCtx, &mockK8sClient{})

	// when
	description := installer.GetStepDescription()

	// then
	assert.Equal(t, "Install service discovery from https://url.server.com/service-discovery/resource.yaml", description)
}

func TestServiceDiscoveryInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("should perform an installation without resource modification", func(t *testing.T) {
		// given
		var yamlBytes apply.YamlDocument = []byte("yaml result goes here")

		mockedResourceRegistryClient := newMockResourceRegistryClient(t)
		mockedResourceRegistryClient.EXPECT().GetResourceFileContent(serviceDiscoveryResourceURL).Return(yamlBytes, nil)
		mockedK8sClient := newMockK8sClient(t)
		mockedK8sClient.EXPECT().ApplyWithOwner(yamlBytes, testTargetNamespaceName, mock.Anything).Return(nil)

		installer := serviceDiscoveryInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            serviceDiscoveryResourceURL,
			resourceRegistryClient: mockedResourceRegistryClient,
			k8sClient:              mockedK8sClient,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
	})
}
