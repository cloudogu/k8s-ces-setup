package component

import (
	"fmt"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

const etcdServerResourceURL = "https://url.server.com/etcd/resource.yaml"

var etcdServerSetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: ctx.Config{
		TargetNamespace:       testTargetNamespaceName,
		EtcdServerResourceURL: etcdServerResourceURL,
	},
}

func TestNewEtcdServerInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual, _ := NewEtcdServerInstallerStep(&rest.Config{}, etcdServerSetupCtx)

	// then
	assert.NotNil(t, actual)
}

func TestEtcdServerInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewEtcdServerInstallerStep(&rest.Config{}, etcdServerSetupCtx)

	// when
	description := installer.GetStepDescription()

	// then
	assert.Equal(t, "Install etcd server from https://url.server.com/etcd/resource.yaml", description)
}

func TestEtcdServerInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("should perform an installation without resource modification", func(t *testing.T) {
		// given
		yamlBytes := []byte("yaml result goes here")

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", etcdServerResourceURL).Return(yamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", yamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", yamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", yamlBytes, testTargetNamespaceName).Return(nil)

		installer := etcdServerInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            etcdServerResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mockedFileClient.AssertExpectations(t)
		mockedK8sClient.AssertExpectations(t)
		mockedFileModder.AssertExpectations(t)
	})

	t.Run("should split yaml file into two parts and apply them each", func(t *testing.T) {
		// given
		yamlDoc1 := `yamlDoc1: 1
	namespace: aNamespaceToBeReplaced`
		yamlDoc2 := `yamlDoc1: 2
	namespace: aNamespaceToBeReplaced`
		yamlBytes := []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", etcdServerResourceURL).Return(yamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", yamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", yamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", []byte(yamlDoc1+"\n"), testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", []byte(yamlDoc2+"\n"), testTargetNamespaceName).Return(nil)

		installer := etcdServerInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            etcdServerResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mockedFileClient.AssertExpectations(t)
		mockedK8sClient.AssertExpectations(t)
		mockedFileModder.AssertExpectations(t)
	})
	t.Run("should fail on second apply", func(t *testing.T) {
		// given
		yamlDoc1 := `yamlDoc1: 1
	namespace: aNamespaceToBeReplaced`
		yamlDoc2 := `yamlDoc1: 2
	namespace: aNamespaceToBeReplaced`
		yamlBytes := []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", etcdServerResourceURL).Return(yamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", yamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", yamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", []byte(yamlDoc1+"\n"), testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", []byte(yamlDoc2+"\n"), testTargetNamespaceName).Return(assert.AnError)

		installer := etcdServerInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            etcdServerResourceURL,
			fileClient:             mockedFileClient,
			k8sClient:              mockedK8sClient,
			fileContentModificator: mockedFileModder,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.Error(t, err)
		mockedFileClient.AssertExpectations(t)
		mockedK8sClient.AssertExpectations(t)
		mockedFileModder.AssertExpectations(t)
	})
}
