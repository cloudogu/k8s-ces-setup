package component

import (
	"fmt"
	"github.com/cloudogu/k8s-apply-lib/apply"
	"testing"

	"github.com/stretchr/testify/mock"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/client-go/rest"
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
	actual, _ := NewServiceDiscoveryInstallerStep(&rest.Config{}, serviceDiscoverySetupCtx)

	// then
	assert.NotNil(t, actual)
}

func TestServiceDiscoveryInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewServiceDiscoveryInstallerStep(&rest.Config{}, serviceDiscoverySetupCtx)

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

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", serviceDiscoveryResourceURL).Return([]byte(yamlBytes), nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", []byte(yamlBytes), testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", []byte(yamlBytes))
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", yamlBytes, testTargetNamespaceName).Return(nil)

		installer := serviceDiscoveryInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            serviceDiscoveryResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockedFileClient, mockedK8sClient, mockedFileModder)
	})

	t.Run("should split yaml file into two parts and apply them each", func(t *testing.T) {
		// given
		yamlDoc1 := `yamlDoc1: 1
	namespace: aNamespaceToBeReplaced`
		yamlDoc2 := `yamlDoc1: 2
	namespace: aNamespaceToBeReplaced`
		var yamlBytes apply.YamlDocument = []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))
		var typedYamlDoc1 apply.YamlDocument = []byte(yamlDoc1 + "\n")
		var typedYamlDoc2 apply.YamlDocument = []byte(yamlDoc2 + "\n")

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", serviceDiscoveryResourceURL).Return([]byte(yamlBytes), nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", []byte(yamlBytes), testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", []byte(yamlBytes))
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", typedYamlDoc1, testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", typedYamlDoc2, testTargetNamespaceName).Return(nil)

		installer := serviceDiscoveryInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            serviceDiscoveryResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockedFileClient, mockedK8sClient, mockedFileModder)
	})
	t.Run("should fail on second apply", func(t *testing.T) {
		// given
		yamlDoc1 := `yamlDoc1: 1
	namespace: aNamespaceToBeReplaced`
		yamlDoc2 := `yamlDoc1: 2
	namespace: aNamespaceToBeReplaced`
		var yamlBytes apply.YamlDocument = []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))
		var typedYamlDoc1 apply.YamlDocument = []byte(yamlDoc1 + "\n")
		var typedYamlDoc2 apply.YamlDocument = []byte(yamlDoc2 + "\n")

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", serviceDiscoveryResourceURL).Return([]byte(yamlBytes), nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", []byte(yamlBytes), testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", []byte(yamlBytes))
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", typedYamlDoc1, testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", typedYamlDoc2, testTargetNamespaceName).Return(assert.AnError)

		installer := serviceDiscoveryInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            serviceDiscoveryResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.Error(t, err)
		mock.AssertExpectationsForObjects(t, mockedFileClient, mockedK8sClient, mockedFileModder)
	})
}
