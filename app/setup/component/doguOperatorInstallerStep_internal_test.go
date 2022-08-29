package component

import (
	"github.com/cloudogu/k8s-apply-lib/apply"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		mockedFileClient := &mockFileClient{}
		mockedResourceRegistryClient := &mocks.ResourceRegistryClient{}
		mockedResourceRegistryClient.On("GetResourceFileContent", doguOperatorURL).Return([]byte(doguOpYamlBytes), nil)
		mockedK8sClient := &mockK8sClient{}
		mockedK8sClient.On("ApplyWithOwner", doguOpYamlBytes, testTargetNamespaceName, mock.Anything).Return(nil)

		installer := doguOperatorInstallerStep{
			namespace:              testTargetNamespaceName,
			resourceURL:            doguOperatorURL,
			resourceRegistryClient: mockedResourceRegistryClient,
			k8sClient:              mockedK8sClient,
		}

		// when
		err := installer.PerformSetupStep()

		// then
		require.NoError(t, err)
		mockedFileClient.AssertExpectations(t)
		mockedK8sClient.AssertExpectations(t)
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

func (mkc *mockK8sClient) Apply(yamlResources apply.YamlDocument, namespace string) error {
	args := mkc.Called(yamlResources, namespace)
	return args.Error(0)
}

func (mkc *mockK8sClient) ApplyWithOwner(yamlResources apply.YamlDocument, namespace string, resource metav1.Object) error {
	args := mkc.Called(yamlResources, namespace, resource)
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
