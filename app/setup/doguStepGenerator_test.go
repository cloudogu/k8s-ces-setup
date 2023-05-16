package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote/mocks"

	"github.com/cloudogu/k8s-ces-setup/app/context"
)

func TestNewDoguStepGenerator(t *testing.T) {
	t.Run("creating new generator fails by creating ecosystem client from client config", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		clusterConfig.Host = "?(?=)/()?;:#=)A=#?);#:########--------/-*/-*/*+23+435647864645a+5dfa+6523+5a6rt+5e+qA=%);=§"
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}

		mockRegistry := &mocks.Registry{}

		// when
		_, err := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create K8s EcoSystem client")
	})

	t.Run("creating new generator fails by retrieving dogu from registry", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("Get", "official/ldap").Return(nil, assert.AnError)

		// when
		_, err := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get latest version of dogu [official/ldap]")
	})

	t.Run("creating new generator fails by retrieving versioned dogu from registry", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap:1.2.3-4", "official/cas:4.3.2-1"}}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("GetVersion", "official/ldap", "1.2.3-4").Return(nil, assert.AnError)

		// when
		_, err := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get version [1.2.3-4] of dogu [official/ldap]")
	})

	t.Run("create new generator", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}
		doguCas := &core.Dogu{Name: "cas"}
		doguLdap := &core.Dogu{Name: "ldap"}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("Get", "official/ldap").Return(doguLdap, nil)
		mockRegistry.On("Get", "official/cas").Return(doguCas, nil)

		// when
		generator, err := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.NoError(t, err)
		assert.NotNil(t, generator)
		assert.Len(t, *generator.Dogus, 2)
	})
}

func Test_doguStepGenerator_GenerateSteps(t *testing.T) {
	t.Run("generate dogu steps", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas", "official/postfix:1.0.0-1", "official/postgres", "official/redmine:10.0.0-5"}}
		doguCas := &core.Dogu{Name: "cas", Version: "6.5.4-2"}
		doguLdap := &core.Dogu{Name: "ldap", Version: "2.1.0-1"}
		doguPostfix := &core.Dogu{Name: "postfix", Version: "1.0.0-1"}
		doguPostgres := &core.Dogu{Name: "postgres", Version: "0.3.4-0"}
		doguRedmine := &core.Dogu{Name: "redmine", Version: "10.0.0-5"}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("Get", "official/ldap").Return(doguLdap, nil)
		mockRegistry.On("Get", "official/cas").Return(doguCas, nil)
		mockRegistry.On("Get", "official/postgres").Return(doguPostgres, nil)
		mockRegistry.On("GetVersion", "official/postfix", "1.0.0-1").Return(doguPostfix, nil)
		mockRegistry.On("GetVersion", "official/redmine", "10.0.0-5").Return(doguRedmine, nil)
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// when
		doguSteps, _ := generator.GenerateSteps()

		// then
		assert.NotNil(t, generator)
		assert.Len(t, doguSteps, 5)
		assert.Equal(t, "Installing dogu [cas]", doguSteps[0].GetStepDescription())
		assert.Equal(t, "Installing dogu [ldap]", doguSteps[1].GetStepDescription())
		assert.Equal(t, "Installing dogu [postfix]", doguSteps[2].GetStepDescription())
		assert.Equal(t, "Installing dogu [postgres]", doguSteps[3].GetStepDescription())
		assert.Equal(t, "Installing dogu [redmine]", doguSteps[4].GetStepDescription())
	})

	t.Run("generate dogu steps with service account dependencies", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas", "official/postfix:1.0.0-1", "official/postgres", "official/redmine:10.0.0-5"}}
		doguCas := &core.Dogu{Name: "cas", Version: "6.5.4-2", ServiceAccounts: []core.ServiceAccount{{Type: "ldap"}}, Dependencies: []core.Dependency{{Type: "dogu", Name: "ldap"}}}
		doguLdap := &core.Dogu{Name: "ldap", Version: "2.1.0-1"}
		doguPostfix := &core.Dogu{Name: "postfix", Version: "1.0.0-1"}
		doguPostgres := &core.Dogu{Name: "postgres", Version: "0.3.4-0", ServiceAccounts: []core.ServiceAccount{{Type: "cas"}, {Type: "ldap"}}, Dependencies: []core.Dependency{{Type: "dogu", Name: "cas"}, {Type: "dogu", Name: "ldap"}}}
		doguRedmine := &core.Dogu{Name: "redmine", Version: "10.0.0-5", ServiceAccounts: []core.ServiceAccount{{Type: "postgres"}, {Type: "postfix"}}, Dependencies: []core.Dependency{{Type: "dogu", Name: "postgres"}, {Type: "dogu", Name: "postfix"}, {Type: "dogu", Name: "cas"}}}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("Get", "official/ldap").Return(doguLdap, nil)
		mockRegistry.On("Get", "official/cas").Return(doguCas, nil)
		mockRegistry.On("Get", "official/postgres").Return(doguPostgres, nil)
		mockRegistry.On("GetVersion", "official/postfix", "1.0.0-1").Return(doguPostfix, nil)
		mockRegistry.On("GetVersion", "official/redmine", "10.0.0-5").Return(doguRedmine, nil)
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// when
		doguSteps, _ := generator.GenerateSteps()

		// then
		assert.NotNil(t, generator)
		assert.Len(t, doguSteps, 9)
		assert.Equal(t, "Installing dogu [ldap]", doguSteps[0].GetStepDescription())
		assert.Equal(t, "Installing dogu [postfix]", doguSteps[1].GetStepDescription())
		assert.Equal(t, "Wait for pod with selector dogu.name=ldap to be ready", doguSteps[2].GetStepDescription())
		assert.Equal(t, "Installing dogu [cas]", doguSteps[3].GetStepDescription())
		assert.Equal(t, "Wait for pod with selector dogu.name=cas to be ready", doguSteps[4].GetStepDescription())
		assert.Equal(t, "Installing dogu [postgres]", doguSteps[5].GetStepDescription())
		assert.Equal(t, "Wait for pod with selector dogu.name=postgres to be ready", doguSteps[6].GetStepDescription())
		assert.Equal(t, "Wait for pod with selector dogu.name=postfix to be ready", doguSteps[7].GetStepDescription())
		assert.Equal(t, "Installing dogu [redmine]", doguSteps[8].GetStepDescription())
	})
}

func Test_doguStepGenerator_createWaitStepForDogu(t *testing.T) {
	singleFakeStep := &fakeExecutorStep{}
	t.Run("generates wait step to wait for dogus", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		dogus := context.Dogus{}

		clusterConfig := &rest.Config{}
		serviceAccount := core.ServiceAccount{
			Type: "postfix",
			Kind: "",
		}
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, nil, "ns")
		waitList := map[string]bool{"dogu.name=ldap": true}
		allStepsTillNow := []ExecutorStep{singleFakeStep}

		// when
		actualSteps := generator.createWaitStepForDogu(serviceAccount, waitList, allStepsTillNow)

		// then
		assert.NotNil(t, actualSteps)
		assert.Len(t, actualSteps, 2)
		assert.Contains(t, "Wait for pod with selector dogu.name=your-most-favorite to be ready", actualSteps[0].GetStepDescription())
		assert.Contains(t, "Wait for pod with selector dogu.name=postfix to be ready", actualSteps[1].GetStepDescription())
	})
	t.Run("generates wait step to wait for dogus (explicit kind)", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		dogus := context.Dogus{}

		clusterConfig := &rest.Config{}
		serviceAccount := core.ServiceAccount{
			Type: "postfix",
			Kind: "dogu",
		}
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, nil, "ns")
		waitList := map[string]bool{"dogu.name=ldap": true}
		allStepsTillNow := []ExecutorStep{singleFakeStep}

		// when
		actualSteps := generator.createWaitStepForDogu(serviceAccount, waitList, allStepsTillNow)

		// then
		assert.NotNil(t, actualSteps)
		assert.Len(t, actualSteps, 2)
		assert.Contains(t, "Wait for pod with selector dogu.name=your-most-favorite to be ready", actualSteps[0].GetStepDescription())
		assert.Contains(t, "Wait for pod with selector dogu.name=postfix to be ready", actualSteps[1].GetStepDescription())
	})
	t.Run("does not generate wait step because there is already a similar waiting step", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		dogus := context.Dogus{}

		clusterConfig := &rest.Config{}
		serviceAccount := core.ServiceAccount{
			Type: "postfix",
			Kind: "",
		}
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, nil, "ns")
		waitList := map[string]bool{"dogu.name=postfix": true}
		allStepsTillNow := []ExecutorStep{singleFakeStep}

		// when
		actualSteps := generator.createWaitStepForDogu(serviceAccount, waitList, allStepsTillNow)

		// then
		assert.NotNil(t, actualSteps)
		assert.Len(t, actualSteps, 1)
		assert.Contains(t, "Wait for pod with selector dogu.name=your-most-favorite to be ready", actualSteps[0].GetStepDescription())
	})
}

func Test_doguStepGenerator_createWaitStepForK8sComponent(t *testing.T) {
	singleFakeStep := &fakeExecutorStep{}
	t.Run("generates wait step to wait for dogus", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		dogus := context.Dogus{}

		clusterConfig := &rest.Config{}
		serviceAccount := core.ServiceAccount{
			Type: "k8s-dogu-operator",
			Kind: "k8s",
		}
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, nil, "ns")
		waitList := map[string]bool{"dogu.name=ldap": true}
		allStepsTillNow := []ExecutorStep{singleFakeStep}

		// when
		actualSteps := generator.createWaitStepForK8sComponent(serviceAccount, waitList, allStepsTillNow)

		// then
		assert.NotNil(t, actualSteps)
		assert.Len(t, actualSteps, 2)
		assert.Contains(t, "Wait for pod with selector dogu.name=your-most-favorite to be ready", actualSteps[0].GetStepDescription())
		assert.Contains(t, "Wait for pod with selector app.kubernetes.io/name=k8s-dogu-operator to be ready", actualSteps[1].GetStepDescription())
	})
	t.Run("does not generate wait step because there is already a similar waiting step", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		dogus := context.Dogus{}

		clusterConfig := &rest.Config{}
		serviceAccount := core.ServiceAccount{
			Type: "k8s-dogu-operator",
			Kind: "k8s",
		}
		generator, _ := NewDoguStepGenerator(clientMock, clusterConfig, dogus, nil, "ns")
		waitList := map[string]bool{"app.kubernetes.io/name=k8s-dogu-operator": true}
		allStepsTillNow := []ExecutorStep{singleFakeStep}

		// when
		actualSteps := generator.createWaitStepForK8sComponent(serviceAccount, waitList, allStepsTillNow)

		// then
		assert.NotNil(t, actualSteps)
		assert.Len(t, actualSteps, 1)
		assert.Contains(t, "Wait for pod with selector dogu.name=your-most-favorite to be ready", actualSteps[0].GetStepDescription())
	})
}

type fakeExecutorStep struct {
}

func (f *fakeExecutorStep) GetStepDescription() string {
	return "Wait for pod with selector dogu.name=your-most-favorite to be ready"
}

func (f *fakeExecutorStep) PerformSetupStep() error {
	return nil
}
