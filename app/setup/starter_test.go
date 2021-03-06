package setup

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestStarter_StartSetup(t *testing.T) {
	setupContext := context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}, StartupConfiguration: &context.SetupConfiguration{Naming: context.Naming{CertificateType: "selfsigned"}}}
	starter := &Starter{}
	starter.SetupContext = &setupContext
	starter.ClientSet = &fake.Clientset{}
	starter.Namespace = "test"

	t.Run("success", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(nil)
		executorMock.On("RegisterValidationStep").Return(nil)
		executorMock.On("RegisterComponentSetupSteps").Return(nil)
		executorMock.On("RegisterDataSetupSteps", mock.Anything).Return(nil)
		executorMock.On("RegisterDoguInstallationSteps").Return(nil)
		executorMock.On("PerformSetup").Return(nil, "")
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, executorMock)
	})

	t.Run("failed because setup is busy", func(t *testing.T) {
		// given
		doneStarter := &Starter{}
		doneStarter.SetupContext = &setupContext
		doneStarter.Namespace = "test"
		data := make(map[string]string)
		data[context.SetupStateKey] = context.SetupStateInstalling
		configmap := &v1.ConfigMap{ObjectMeta: v12.ObjectMeta{Name: context.SetupStateConfigMap, Namespace: "test"}, Data: data}
		doneStarter.ClientSet = fake.NewSimpleClientset(configmap)

		// when
		err := doneStarter.StartSetup()

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
		doneStarter.ClientSet = fake.NewSimpleClientset(configmap)

		// when
		err := doneStarter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "setup is busy or already done")
	})

	t.Run("failed to register ssl generate step", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(assert.AnError)
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register ssl generation setup step")
		mock.AssertExpectationsForObjects(t, executorMock)
	})

	t.Run("failed to register validation steps", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(nil)
		executorMock.On("RegisterValidationStep").Return(assert.AnError)
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register validation setup steps")
		mock.AssertExpectationsForObjects(t, executorMock)
	})

	t.Run("failed to register component setup steps", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(nil)
		executorMock.On("RegisterValidationStep").Return(nil)
		executorMock.On("RegisterComponentSetupSteps").Return(assert.AnError)
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register component setup steps")
		mock.AssertExpectationsForObjects(t, executorMock)
	})

	t.Run("failed to register data setup steps", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(nil)
		executorMock.On("RegisterValidationStep").Return(nil)
		executorMock.On("RegisterComponentSetupSteps").Return(nil)
		executorMock.On("RegisterDataSetupSteps", mock.Anything).Return(assert.AnError)
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register data setup steps")
		mock.AssertExpectationsForObjects(t, executorMock)
	})

	t.Run("failed to register dogu installation steps", func(t *testing.T) {
		// given
		executorMock := &mocks.SetupExecutor{}
		executorMock.On("RegisterSSLGenerationStep").Return(nil)
		executorMock.On("RegisterValidationStep").Return(nil)
		executorMock.On("RegisterComponentSetupSteps").Return(nil)
		executorMock.On("RegisterDataSetupSteps", mock.Anything).Return(nil)
		executorMock.On("RegisterDoguInstallationSteps").Return(assert.AnError)
		starter.SetupExecutor = executorMock

		// when
		err := starter.StartSetup()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register dogu installation steps")
		mock.AssertExpectationsForObjects(t, executorMock)
	})
}
