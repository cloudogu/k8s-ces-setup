package component

import (
	"fmt"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

const testTargetNamespaceName = "myfavouritenamespace-1"
const doguOperatorURL = "http://url.server.com/dogu/operator.yaml"

var doguOperatorSetupCtx = &ctx.SetupContext{
	AppVersion: "1.2.3",
	AppConfig: ctx.Config{
		TargetNamespace: testTargetNamespaceName,
		DoguOperatorURL: doguOperatorURL,
	},
}

func TestDefaultFileContentModificator_replaceNamespace(t *testing.T) {
	t.Run("should replace namespace within simple namespaced resource", func(t *testing.T) {
		input := simpleTestRoleBinding()
		sut := &defaultFileContentModificator{}

		// when
		actual := sut.replaceNamespacedResources(input, testTargetNamespaceName)

		// then
		expected := fmt.Sprintf(`apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8s-dogu-operator-leader-election-rolebinding
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8s-dogu-operator-leader-election-role
subjects:
- kind: ServiceAccount
  name: k8s-dogu-operator-controller-manager
  namespace: %s
`, testTargetNamespaceName, testTargetNamespaceName)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("should replace namespace in complex resources", func(t *testing.T) {
		input := twoNamespacedResources()
		sut := &defaultFileContentModificator{}

		// when
		actual := sut.replaceNamespacedResources(input, testTargetNamespaceName)

		// then
		expected := fmt.Sprintf(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-dogu-operator-controller-manager
  namespace: %s
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8s-dogu-operator-leader-election-role
  namespace: %s
rules:
- apiGroups:
`,
			testTargetNamespaceName, testTargetNamespaceName)

		assert.Equal(t, expected, string(actual))
	})
}

func TestDefaultFileContentModificator_removeLegacyNamespaceFromResources(t *testing.T) {
	input := []byte(`apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: ecosystem
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.8.0
  creationTimestamp: null
  name: dogus.k8s.cloudogu.com`)
	sut := &defaultFileContentModificator{}

	// when
	actual := sut.removeLegacyNamespaceFromResources(input)

	// then
	expectedNoResource := `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.8.0
  creationTimestamp: null
  name: dogus.k8s.cloudogu.com`
	assert.Equal(t, expectedNoResource, string(actual))
}

func Test_splitYamlFileSections(t *testing.T) {
	t.Run("should return two sections (with leading delimiter)", func(t *testing.T) {
		const simpleMultiLineYaml = `---
test:
---
anotherTest:
`
		input := []byte(simpleMultiLineYaml)

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, "test:\n", string(sections[0]))
		assert.Equal(t, "anotherTest:\n", string(sections[1]))
	})
	t.Run("should return two sections (without leading delimiter)", func(t *testing.T) {
		const simpleMultiLineYaml = `test:
---
anotherTest:
`
		input := []byte(simpleMultiLineYaml)

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, "test:\n", string(sections[0]))
		assert.Equal(t, "anotherTest:\n", string(sections[1]))
	})
	t.Run("should return sections for complex YAML", func(t *testing.T) {
		input := []byte(multiFileYaml())

		// when
		sections := splitYamlFileSections(input)

		// then
		assert.Len(t, sections, 2)
		assert.Equal(t, `# A comment for the service
apiVersion: v1
kind: Service
metadata:
  name: your-app
  app.kubernetes.io/name: your-app
  labels:
    app: your-app
spec:
  type: NodePort
  ports:
`, string(sections[0]))
		assert.Equal(t, `# a comment for the deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: your-app
  name: your-app
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: your-app
  template:
    metadata:
      labels:
        app: your-app
        app.kubernetes.io/name: your-app
    spec:
`, string(sections[1]))
	})
}

func TestNewDoguOperatorInstallerStep(t *testing.T) {
	t.Parallel()

	// when
	actual, _ := NewDoguOperatorInstallerStep(&rest.Config{}, doguOperatorSetupCtx)

	// then
	assert.NotNil(t, actual)
}

func TestDoguOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	// given
	installer, _ := NewDoguOperatorInstallerStep(&rest.Config{}, doguOperatorSetupCtx)

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

func multiFileYaml() string {
	return `---
# A comment for the service
apiVersion: v1
kind: Service
metadata:
  name: your-app
  app.kubernetes.io/name: your-app
  labels:
    app: your-app
spec:
  type: NodePort
  ports:
---
# a comment for the deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: your-app
  name: your-app
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: your-app
  template:
    metadata:
      labels:
        app: your-app
        app.kubernetes.io/name: your-app
    spec:
`
}

func simpleTestRoleBinding() []byte {
	return []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8s-dogu-operator-leader-election-rolebinding
  namespace: ecosystem
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8s-dogu-operator-leader-election-role
subjects:
- kind: ServiceAccount
  name: k8s-dogu-operator-controller-manager
  namespace: ecosystem
`)
}

func twoNamespacedResources() []byte {
	return []byte(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-dogu-operator-controller-manager
  namespace: ecosystem
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8s-dogu-operator-leader-election-role
  namespace: ecosystem
rules:
- apiGroups:
`)
}
