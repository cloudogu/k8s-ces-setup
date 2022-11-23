package setup_test

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/remote/mocks"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestNewDoguStepGenerator(t *testing.T) {
	t.Run("creating new generator fails by creating rest client on client config", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		clusterConfig.Host = "?(?=)/()?;:#=)A=#?);#:########--------/-*/-*/*+23+435647864645a+5dfa+6523+5a6rt+5e+qA=%);=§"
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}

		mockRegistry := &mocks.Registry{}

		// when
		_, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "annot create kubernetes RestClient")
	})

	t.Run("creating new generator fails by creating rest client on AddToScheme", func(t *testing.T) {
		// given
		orifignalAddToScheme := setup.AddToScheme
		defer func() { setup.AddToScheme = orifignalAddToScheme }()
		setup.AddToScheme = func(s *runtime.Scheme) error {
			return assert.AnError
		}

		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}

		mockRegistry := &mocks.Registry{}

		// when
		_, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("creating new generator fails by retrieving dogu from registry", func(t *testing.T) {
		// given
		clientMock := fake.NewSimpleClientset()
		clusterConfig := &rest.Config{}
		dogus := context.Dogus{Install: []string{"official/ldap", "official/cas"}}

		mockRegistry := &mocks.Registry{}
		mockRegistry.On("Get", mock.Anything).Return(nil, assert.AnError)

		// when
		_, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
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
		generator, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

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
		generator, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// when
		doguSteps := generator.GenerateSteps()

		// then
		require.NoError(t, err)
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
		generator, err := setup.NewDoguStepGenerator(clientMock, clusterConfig, dogus, mockRegistry, "mynamespace")

		// when
		doguSteps := generator.GenerateSteps()

		// then
		require.NoError(t, err)
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
