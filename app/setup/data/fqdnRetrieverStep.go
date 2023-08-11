package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

var backoff = wait.Backoff{
	Duration: 5000 * time.Millisecond,
	Factor:   1,
	Jitter:   0,
	Steps:    25,
	Cap:      3 * time.Minute,
}

type fqdnRetrieverStep struct {
	config    *appcontext.SetupJsonConfiguration
	clientSet kubernetes.Interface
	namespace string
}

// NewFQDNRetrieverStep creates a new setup step sets the FQDN
func NewFQDNRetrieverStep(config *appcontext.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *fqdnRetrieverStep {
	return &fqdnRetrieverStep{config: config, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *fqdnRetrieverStep) GetStepDescription() string {
	return "Retrieving a new FQDN from the IP of a loadbalancer service"
}

// PerformSetupStep creates a loadbalancer service and sets the loadbalancer IP as the new FQDN.
func (fcs *fqdnRetrieverStep) PerformSetupStep(ctx context.Context) error {
	return fcs.setFQDNFromLoadbalancerIP(ctx)
}

func (fcs *fqdnRetrieverStep) setFQDNFromLoadbalancerIP(ctx context.Context) error {
	return retry.OnError(backoff, serviceRetry, func() error {
		logrus.Debug("Try retrieving service...")
		service, err := fcs.clientSet.CoreV1().Services(fcs.namespace).Get(ctx, cesLoadbalancerName, metav1.GetOptions{})

		if errors.IsNotFound(err) || len(service.Status.LoadBalancer.Ingress) <= 0 {
			logrus.Debugf("wait for service %s to be instantiated", cesLoadbalancerName)
			return fmt.Errorf("service not yet ready %s: %w", cesLoadbalancerName, err)
		}
		if err != nil {
			return err
		}

		loadbalancerIP := service.Status.LoadBalancer.Ingress[0].IP
		fcs.config.Naming.Fqdn = loadbalancerIP
		logrus.Infof("Loadbalancer IP succesfully retrieved and set as new FQDN")
		return nil
	})
}

func serviceRetry(err error) bool {
	return strings.Contains(err.Error(), "service not yet ready")
}
