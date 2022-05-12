package setup_test

import (
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/k8s-ces-setup/app/setup"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
)

func TestFinisher_FinishSetup(t *testing.T) {
	t.Run("finish setup with no resources", func(t *testing.T) {
		// given
		fakeclient := fake.NewSimpleClientset()
		finisher := setup.NewFinisher(fakeclient, "mytestnamespace")

		// when
		err := finisher.FinishSetup()

		// then
		require.NoError(t, err)

		cm, err := context.GetSetupConfigMap(fakeclient, "mytestnamespace")
		require.NoError(t, err)

		assert.Equal(t, "installed", cm.Data[context.SetupStateKey])
	})
}
