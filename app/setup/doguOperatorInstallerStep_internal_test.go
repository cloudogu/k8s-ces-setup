package setup

import (
	"fmt"
	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const doguOperatorURL = "http://url.server.com/dogu/operator.yaml"

var doguOperatorSetupCtx = ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: ctx.Config{
		TargetNamespace: testTargetNamespaceName,
		DoguOperatorURL: doguOperatorURL,
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
	installer := newDoguOperatorInstallerStep(nil, doguOperatorSetupCtx)

	// when
	description := installer.GetStepDescription()

	// then
	assert.Equal(t, "Install dogu operator from http://url.server.com/dogu/operator.yaml", description)
}

func TestDoguOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("should perform an installation without resource modification", func(t *testing.T) {
		// given
		doguOpYamlBytes := []byte("yaml result goes here")

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", doguOperatorURL).Return(doguOpYamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", doguOpYamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", doguOpYamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", doguOpYamlBytes, testTargetNamespaceName).Return(nil)

		installer := doguOperatorInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            doguOperatorURL,
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
		doguOpYamlBytes := []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", doguOperatorURL).Return(doguOpYamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", doguOpYamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", doguOpYamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", []byte(yamlDoc1+"\n"), testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", []byte(yamlDoc2+"\n"), testTargetNamespaceName).Return(nil)

		installer := doguOperatorInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            doguOperatorURL,
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
		doguOpYamlBytes := []byte(fmt.Sprintf(`---
%v
---
%v
`, yamlDoc1, yamlDoc2))

		mockedFileClient := &mockFileClient{}
		mockedFileClient.On("Get", doguOperatorURL).Return(doguOpYamlBytes, nil)
		mockedFileModder := &mockFileModder{}
		mockedFileModder.On("replaceNamespacedResources", doguOpYamlBytes, testTargetNamespaceName)
		mockedFileModder.On("removeLegacyNamespaceFromResources", doguOpYamlBytes)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("Apply", []byte(yamlDoc1+"\n"), testTargetNamespaceName).Return(nil)
		mockedK8sClient.On("Apply", []byte(yamlDoc2+"\n"), testTargetNamespaceName).Return(assert.AnError)

		installer := doguOperatorInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            doguOperatorURL,
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

type mockFileClient struct {
	mock.Mock
}

func (mfc *mockFileClient) Get(url string) ([]byte, error) {
	args := mfc.Called(url)
	return args.Get(0).([]byte), args.Error(1)
}

type mockK8sClient struct {
	mock.Mock
}

func (mkc *mockK8sClient) Apply(yamlResources []byte, namespace string) error {
	args := mkc.Called(yamlResources, namespace)
	return args.Error(0)
}

type mockFileModder struct {
	mock.Mock
}

func (mfm *mockFileModder) replaceNamespacedResources(content []byte, namespace string) []byte {
	mfm.Called(content, namespace)
	return content
}

func (mfm *mockFileModder) removeLegacyNamespaceFromResources(content []byte) []byte {
	mfm.Called(content)
	return content
}
