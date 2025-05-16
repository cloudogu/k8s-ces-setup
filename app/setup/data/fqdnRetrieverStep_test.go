package data

import (
	"bytes"
	"context"
	appctx "github.com/cloudogu/k8s-ces-setup/v4/app/context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"os"
	"testing"
	"time"
)

const (
	testNamespace = "ecosystem"
)

func Test_fqdnRetrieverStep_PerformSetupStep(t *testing.T) {
	t.Run("should successfully set FQDN to IP when the load-balancer receives an external IP address", func(t *testing.T) {
		// given
		mockedLoadBalancerResource := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cesLoadbalancerName,
				Namespace: testNamespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{DoguLabelName: nginxIngressName},
			},
		}
		fakeClient := fake.NewSimpleClientset(mockedLoadBalancerResource)
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}, Dogus: appctx.Dogus{Install: []string{nginxIngressName}}}

		sut := NewFQDNRetrieverStep(config, fakeClient, testNamespace)

		// simulate asynchronous IP setting by cluster provider
		timer := time.NewTimer(time.Second * 2)
		go func() {
			<-timer.C
			patch := []byte(`{"status":{"loadBalancer":{"ingress":[{"ip": "111.222.111.222"}]}}}`)
			service, err := fakeClient.CoreV1().Services(testNamespace).Patch(
				context.Background(),
				cesLoadbalancerName,
				types.MergePatchType,
				patch,
				metav1.PatchOptions{},
			)
			require.NoError(t, err)
			assert.NotNil(t, service)
		}()

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, "111.222.111.222", config.Naming.Fqdn)
	})
}
func TestCreateLoadBalancerStep_PerformSetupStep(t *testing.T) {
	t.Run("failed due to missing nginx-ingress", func(t *testing.T) {
		// given
		config := &appctx.SetupJsonConfiguration{Naming: appctx.Naming{Fqdn: ""}}

		step := NewCreateLoadBalancerStep(config, nil, testNamespace)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid configuration: FQDN can only be created if nginx-ingress will be installed")
	})
}

func TestNewFQDNRetrieverStep(t *testing.T) {
	// when
	step := NewFQDNRetrieverStep(nil, nil, "")

	// then
	require.NotNil(t, step)
}

func Test_fqdnRetrieverStep_GetStepDescription(t *testing.T) {
	// given
	step := NewFQDNRetrieverStep(nil, nil, "")

	// when
	description := step.GetStepDescription()

	// then
	assert.Equal(t, "Retrieving a new FQDN from the IP of a loadbalancer service", description)
}

func Test_readFqdnFromLoadBalancerWaitTimeoutMinsEnv(t *testing.T) {
	tests := []struct {
		name        string
		setEnvVar   bool
		envVarValue string
		want        time.Duration
		wantLogs    bool
		wantedLogs  string
		logLevel    logrus.Level
	}{
		{
			name:       "Environment variable not set",
			setEnvVar:  false,
			want:       time.Duration(15),
			wantLogs:   true,
			wantedLogs: "failed to read FQDN_FROM_LOAD_BALANCER_WAIT_TIMEOUT_MINS environment variable, using default value of 15",
			logLevel:   logrus.DebugLevel,
		},
		{
			name:        "Environment variable not set correctly",
			setEnvVar:   true,
			envVarValue: "15//",
			want:        time.Duration(15),
			wantLogs:    true,
			wantedLogs:  "failed to parse FQDN_FROM_LOAD_BALANCER_WAIT_TIMEOUT_MINS environment variable, using default value of 15",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "read negative environment variable",
			setEnvVar:   true,
			envVarValue: "-20",
			want:        time.Duration(15),
			wantLogs:    true,
			wantedLogs:  "parsed value (-20) is smaller than 0, using default value of 15",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "Successfully read environment variable",
			setEnvVar:   true,
			envVarValue: "20",
			want:        time.Duration(20),
			wantLogs:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvVar {
				err := os.Setenv(fqdnFromLoadBalancerWaitTimeoutMinsEnv, tt.envVarValue)
				require.NoError(t, err)
			}
			var result = time.Duration(0)

			var logOutput bytes.Buffer

			originalOutput := logrus.StandardLogger().Out
			originalLevel := logrus.StandardLogger().Level
			if tt.wantLogs {
				logrus.StandardLogger().SetOutput(&logOutput)
				logrus.StandardLogger().SetLevel(tt.logLevel)
			}

			result = readFqdnFromLoadBalancerWaitTimeoutMinsEnv()

			logrus.StandardLogger().SetOutput(originalOutput)
			logrus.StandardLogger().SetLevel(originalLevel)

			logs := logOutput.String()

			assert.Equalf(t, tt.want, result, "readFqdnFromLoadBalancerWaitTimeoutMinsEnv()")

			if tt.wantLogs {
				assert.Contains(t, logs, tt.wantedLogs)
			}
		})
	}
}
