package setup

import (
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestStarter_StartSetup(t *testing.T) {
	setupContext := context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}, SetupJsonConfiguration: &context.SetupJsonConfiguration{Naming: context.Naming{CertificateType: "selfsigned"}}}
	starter := &Starter{}
	starter.SetupContext = &setupContext
	defaultSA := &v1.ServiceAccount{
		ObjectMeta: v12.ObjectMeta{
			Name:      "default",
			Namespace: "test",
		},
	}
	starter.Namespace = "test"

	t.Run("successful run without FQDN", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(nil)
		expect.RegisterComponentSetupSteps().Return(nil)
		expect.RegisterDataSetupSteps(mock.Anything, mock.Anything).Return(nil)
		expect.RegisterDoguInstallationSteps(mock.Anything).Return(nil)
		expect.PerformSetup(testCtx).Return(nil, "")
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.NoError(t, err)
		// test default service account token automount deactivation
		sa, err := starter.ClientSet.CoreV1().ServiceAccounts(starter.Namespace).Get(testCtx, "default", v12.GetOptions{})
		require.NoError(t, err)
		assert.False(t, *sa.AutomountServiceAccountToken)
	})

	t.Run("successful run with FQDN", func(t *testing.T) {
		// given
		setupContext.SetupJsonConfiguration.Naming.Fqdn = "My-Test-FQDN"
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(nil)
		expect.RegisterDataSetupSteps(mock.Anything, mock.Anything).Return(nil)
		expect.RegisterComponentSetupSteps().Return(nil)
		expect.RegisterDoguInstallationSteps(mock.Anything).Return(nil)
		expect.PerformSetup(testCtx).Return(nil, "")
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.NoError(t, err)
		// test default service account token automount deactivation
		sa, err := starter.ClientSet.CoreV1().ServiceAccounts(starter.Namespace).Get(testCtx, "default", v12.GetOptions{})
		require.NoError(t, err)
		assert.False(t, *sa.AutomountServiceAccountToken)
	})

	t.Run("failed because setup is busy", func(t *testing.T) {
		// given
		doneStarter := &Starter{}
		doneStarter.SetupContext = &setupContext
		doneStarter.Namespace = "test"
		data := make(map[string]string)
		data[context.SetupStateKey] = context.SetupStateInstalling
		configmap := &v1.ConfigMap{ObjectMeta: v12.ObjectMeta{Name: context.SetupStateConfigMap, Namespace: "test"}, Data: data}
		doneStarter.ClientSet = fake.NewClientset(configmap)

		// when
		err := doneStarter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "setup is busy or already done")
	})

	t.Run("failed because setup is done", func(t *testing.T) {
		// given
		doneStarter := &Starter{}
		doneStarter.SetupContext = &setupContext
		doneStarter.Namespace = "test"
		data := make(map[string]string)
		data[context.SetupStateKey] = context.SetupStateInstalled
		configmap := &v1.ConfigMap{ObjectMeta: v12.ObjectMeta{Name: context.SetupStateConfigMap, Namespace: "test"}, Data: data}
		doneStarter.ClientSet = fake.NewClientset(configmap)

		// when
		err := doneStarter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "setup is busy or already done")
	})

	t.Run("failed because setup could not get default service account", func(t *testing.T) {
		// given
		starter := &Starter{}
		starter.SetupContext = &setupContext
		starter.Namespace = "test"
		client := fake.NewClientset()
		client.PrependReactor("get", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, assert.AnError
		})
		starter.ClientSet = client

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get default service account")
	})

	t.Run("failed because setup could not update default service account", func(t *testing.T) {
		// given
		starter := &Starter{}
		starter.SetupContext = &setupContext
		starter.Namespace = "test"
		client := fake.NewClientset()
		client.PrependReactor("get", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, &v1.ServiceAccount{}, nil
		})
		client.PrependReactor("update", "serviceaccounts", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, assert.AnError
		})
		starter.ClientSet = client

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to deactivate token automount on default service account")
	})

	t.Run("failed to register loadbalancer fqdn retriever steps", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		executorMock.EXPECT().RegisterLoadBalancerFQDNRetrieverSteps().Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register steps for creating loadbalancer and retrieving its ip as fqdn")
	})

	t.Run("failed to register ssl generate step", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		executorMock.EXPECT().RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		executorMock.EXPECT().RegisterSSLGenerationStep().Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register ssl generation setup step")
	})

	t.Run("failed to register validation steps", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register validation setup steps")
	})

	t.Run("failed to register component setup steps", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(nil)
		expect.RegisterDataSetupSteps(mock.Anything, mock.Anything).Return(nil)
		expect.RegisterComponentSetupSteps().Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register component setup steps")
	})

	t.Run("failed to register data setup steps", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(nil)
		expect.RegisterDataSetupSteps(mock.Anything, mock.Anything).Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register data setup steps")
	})

	t.Run("failed to register dogu installation steps", func(t *testing.T) {
		// given
		executorMock := NewMockSetupExecutor(t)
		expect := executorMock.EXPECT()
		expect.RegisterLoadBalancerFQDNRetrieverSteps().Return(nil)
		expect.RegisterSSLGenerationStep().Return(nil)
		expect.RegisterValidationStep().Return(nil)
		expect.RegisterComponentSetupSteps().Return(nil)
		expect.RegisterDataSetupSteps(mock.Anything, mock.Anything).Return(nil)
		expect.RegisterDoguInstallationSteps(mock.Anything).Return(assert.AnError)
		starter.SetupExecutor = executorMock
		starter.ClientSet = fake.NewClientset(defaultSA)

		// when
		err := starter.StartSetup(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register dogu installation steps")
	})
}
