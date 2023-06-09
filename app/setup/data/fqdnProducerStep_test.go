package data

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func Test_fqdnCreatorStep_PerformSetupStep(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &context.SetupConfiguration{Naming: context.Naming{Fqdn: ""}, Dogus: context.Dogus{Install: []string{"nginx-ingress"}}}
		fakeClient := fake.NewSimpleClientset()

		step := NewFQDNCreatorStep(config, fakeClient, "ecosystem")

		// when
		//timer := time.NewTimer(time.Second * 2)
		//go func() {
		//	<-timer.C
		//	err := <mache update auf serivce>
		//		require.NoError(t, err)
		//}()
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
		assert.Equal(t, "cert", config.Naming.Certificate)
		assert.NotEmpty(t, "key", config.Naming.CertificateKey)
	})
}
