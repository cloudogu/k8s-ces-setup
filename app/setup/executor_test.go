package setup

import (
	"errors"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-ces-setup/app/context"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mySimpleSetupStep struct {
	PerformedStep  bool
	ErrorOnPerform bool
	Description    string
}

func newSimpleSetupStep(description string, errorOnPerform bool) *mySimpleSetupStep {
	return &mySimpleSetupStep{
		PerformedStep:  false,
		Description:    description,
		ErrorOnPerform: errorOnPerform,
	}
}

func (m *mySimpleSetupStep) GetStepDescription() string {
	return m.Description
}

func (m *mySimpleSetupStep) PerformSetupStep() error {
	if m.ErrorOnPerform {
		return errors.New("failed to do nothing")
	}

	m.PerformedStep = true
	return nil
}

func TestNewExecutor(t *testing.T) {
	t.Parallel()

	// given
	restConfigMock := &rest.Config{}
	clientSetMock := &fake.Clientset{}
	testContext := &context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}, DoguRegistryConfiguration: &context.DoguRegistrySecret{
		Endpoint: "endpoint",
		Username: "username",
		Password: "password",
	}}

	// when
	executor, err := NewExecutor(restConfigMock, clientSetMock, testContext)

	// then
	require.Nil(t, err)
	require.NotNil(t, executor)
}

func TestExecutor_RegisterSetupStep(t *testing.T) {
	t.Parallel()

	t.Run("Register multiple setup steps", func(t *testing.T) {
		// given
		executor := Executor{}
		step1 := newSimpleSetupStep("Step1", false)
		step2 := newSimpleSetupStep("Step2", false)
		step3 := newSimpleSetupStep("Step3", false)

		// when
		executor.RegisterSetupStep(step1)
		executor.RegisterSetupStep(step3)
		executor.RegisterSetupStep(step2)

		// then
		require.NotNil(t, executor.Steps)
		assert.Len(t, executor.Steps, 3)

		assert.Equal(t, step1, executor.Steps[0])
		assert.Equal(t, "Step1", executor.Steps[0].GetStepDescription())

		assert.Equal(t, step3, executor.Steps[1])
		assert.Equal(t, "Step3", executor.Steps[1].GetStepDescription())

		assert.Equal(t, step2, executor.Steps[2])
		assert.Equal(t, "Step2", executor.Steps[2].GetStepDescription())
	})
}

func TestExecutor_PerformSetup(t *testing.T) {
	t.Parallel()

	t.Run("Perform setup with multiple successful setup steps", func(t *testing.T) {
		// given
		executor := Executor{}

		step1 := newSimpleSetupStep("Step1", false)
		step2 := newSimpleSetupStep("Step2", false)
		step3 := newSimpleSetupStep("Step3", false)

		executor.RegisterSetupStep(step1)
		executor.RegisterSetupStep(step3)
		executor.RegisterSetupStep(step2)

		// when
		err, _ := executor.PerformSetup()

		// then
		require.NoError(t, err)
		assert.True(t, step1.PerformedStep)
		assert.True(t, step2.PerformedStep)
		assert.True(t, step3.PerformedStep)
	})

	t.Run("Perform setup with error on setup step", func(t *testing.T) {
		// given
		executor := Executor{}

		step1 := newSimpleSetupStep("Step1", false)
		step2 := newSimpleSetupStep("Step2", true)
		step3 := newSimpleSetupStep("Step3", false)

		executor.RegisterSetupStep(step1)
		executor.RegisterSetupStep(step2)
		executor.RegisterSetupStep(step3)

		// when
		err, uiCause := executor.PerformSetup()

		// then
		require.Error(t, err)
		assert.Equal(t, "Step2", uiCause)
		assert.Equal(t, "failed to perform step [Step2]: failed to do nothing", err.Error())
		assert.True(t, step1.PerformedStep)
		assert.True(t, !step2.PerformedStep)
		assert.True(t, !step3.PerformedStep) // not performed because step 2 could not perform
	})
}

func TestExecutor_RegisterFQDNRetrieverStep(t *testing.T) {
	t.Run("successfully register FQDN retriever step", func(t *testing.T) {
		// given
		testContext := &context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}}
		executor := &Executor{
			ClusterConfig: &rest.Config{},
			SetupContext:  testContext,
		}

		// when
		executor.RegisterFQDNRetrieverStep()

		// then
		assert.True(t, len(executor.Steps) == 1)
		assert.Equal(t, "Retrieving a new FQDN from the IP of a loadbalancer service", executor.Steps[0].GetStepDescription())
	})
}

func TestExecutor_RegisterComponentSetupSteps(t *testing.T) {
	t.Run("successfully register steps", func(t *testing.T) {
		// given
		testContext := &context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}}
		executor := &Executor{
			ClusterConfig: &rest.Config{},
			SetupContext:  testContext,
		}

		// when
		err := executor.RegisterComponentSetupSteps()

		// then
		require.NoError(t, err)
	})

	t.Run("failed to create applier", func(t *testing.T) {
		// given
		testContext := &context.SetupContext{AppConfig: &context.Config{TargetNamespace: "test"}}
		executor := &Executor{
			SetupContext:  testContext,
			ClusterConfig: &rest.Config{ExecProvider: &api.ExecConfig{}, AuthProvider: &api.AuthProviderConfig{}},
		}

		// when
		err := executor.RegisterComponentSetupSteps()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create k8s apply client")
	})
}

func Test_getRemoteConfig(t *testing.T) {
	type args struct {
		endpoint  string
		urlSchema string
	}
	tests := []struct {
		name string
		args args
		want *core.Remote
	}{
		{
			name: "test default url schema",
			args: args{endpoint: "https://example.com/", urlSchema: "default"},
			want: &core.Remote{Endpoint: "https://example.com", URLSchema: "default", CacheDir: "/tmp"},
		},
		{
			name: "test default url schema with 'dogus' suffix",
			args: args{endpoint: "https://example.com/dogus", urlSchema: "default"},
			want: &core.Remote{Endpoint: "https://example.com", URLSchema: "default", CacheDir: "/tmp"},
		},
		{
			name: "test default url schema with 'dogus/' suffix",
			args: args{endpoint: "https://example.com/dogus/", urlSchema: "default"},
			want: &core.Remote{Endpoint: "https://example.com", URLSchema: "default", CacheDir: "/tmp"},
		},
		{
			name: "test non-default url schema",
			args: args{endpoint: "https://example.com/", urlSchema: "index"},
			want: &core.Remote{Endpoint: "https://example.com", URLSchema: "index", CacheDir: "/tmp"},
		},
		{
			name: "test non-default url schema with 'dogus' suffix",
			args: args{endpoint: "https://example.com/dogus", urlSchema: "index"},
			want: &core.Remote{Endpoint: "https://example.com/dogus", URLSchema: "index", CacheDir: "/tmp"},
		},
		{
			name: "test non-default url schema with 'dogus/' suffix",
			args: args{endpoint: "https://example.com/dogus/", urlSchema: "index"},
			want: &core.Remote{Endpoint: "https://example.com/dogus", URLSchema: "index", CacheDir: "/tmp"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getRemoteConfig(tt.args.endpoint, tt.args.urlSchema), "getRemoteConfig(%v, %v)", tt.args.endpoint, tt.args.urlSchema)
		})
	}
}
