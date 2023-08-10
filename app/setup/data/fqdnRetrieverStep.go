package data

import (
	gocontext "context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"strings"
	"time"
)

const (
	cesLoadbalancerName = "ces-loadbalancer"

	// DoguLabelName is used to select a dogu pod by name.
	DoguLabelName    = "dogu.name"
	nginxIngressName = "nginx-ingress"
)

var backoff = wait.Backoff{
	Duration: 5000 * time.Millisecond,
	Factor:   1,
	Jitter:   0,
	Steps:    25,
	Cap:      3 * time.Minute,
}

type fqdnRetrieverStep struct {
	config    *context.SetupJsonConfiguration
	clientSet kubernetes.Interface
	namespace string
}

// NewFQDNRetrieverStep creates a new setup step sets the FQDN
func NewFQDNRetrieverStep(config *context.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *fqdnRetrieverStep {
	return &fqdnRetrieverStep{config: config, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *fqdnRetrieverStep) GetStepDescription() string {
	return "Retrieving a new FQDN from the IP of a loadbalancer service"
}

// PerformSetupStep creates a loadbalancer service and sets the loadbalancer IP as the new FQDN.
func (fcs *fqdnRetrieverStep) PerformSetupStep() error {
	err := fcs.checkIfNginxWillBeInstalled()
	if err != nil {
		return err
	}

	ctx := gocontext.Background()
	err = fcs.createServiceResource(ctx)
	if err != nil {
		return err
	}

	return fcs.setFQDNFromLoadbalancerIP(ctx)
}

func (fcs *fqdnRetrieverStep) checkIfNginxWillBeInstalled() error {
	for _, dogu := range fcs.config.Dogus.Install {
		if strings.Contains(dogu, nginxIngressName) {
			return nil
		}
	}
	return fmt.Errorf("invalid configuration. FQDN can only be created with nginx-ingress installed")
}

func (fcs *fqdnRetrieverStep) createServiceResource(ctx gocontext.Context) error {
	ipSingleStackPolicy := corev1.IPFamilyPolicySingleStack
	serviceResource := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cesLoadbalancerName,
			Namespace: fcs.namespace,
			Labels:    map[string]string{"app": "ces"},
		},
		Spec: corev1.ServiceSpec{
			Type:           corev1.ServiceTypeLoadBalancer,
			IPFamilyPolicy: &ipSingleStackPolicy,
			IPFamilies:     []corev1.IPFamily{corev1.IPv4Protocol},
			Selector:       map[string]string{DoguLabelName: nginxIngressName},
			Ports:          []corev1.ServicePort{createNginxPortResource(80), createNginxPortResource(443)},
		},
	}

	logrus.Debug("Create load balancer service")
	_, err := fcs.clientSet.CoreV1().Services(fcs.namespace).Create(ctx, serviceResource, metav1.CreateOptions{})
	return err
}

func createNginxPortResource(port int) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       fmt.Sprintf("%s-%d", nginxIngressName, port),
		Protocol:   corev1.ProtocolTCP,
		Port:       int32(port),
		TargetPort: intstr.FromInt(port),
	}
}

func (fcs *fqdnRetrieverStep) setFQDNFromLoadbalancerIP(ctx gocontext.Context) error {
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
