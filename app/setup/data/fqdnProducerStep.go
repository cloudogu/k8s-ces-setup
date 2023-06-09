package data

import (
	gocontext "context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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
	DoguLabelName = "dogu.name"
)

var backoff = wait.Backoff{
	Duration: 1500 * time.Millisecond,
	Factor:   1.5,
	Jitter:   0,
	Steps:    25,
	Cap:      3 * time.Minute,
}

type fqdnCreatorStep struct {
	config    *context.SetupConfiguration
	clientSet kubernetes.Interface
	namespace string
}

// NewFQDNCreatorStep creates a new setup step sets the FQDN
func NewFQDNCreatorStep(config *context.SetupConfiguration, clientSet kubernetes.Interface, namespace string) *fqdnCreatorStep {
	return &fqdnCreatorStep{config: config, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *fqdnCreatorStep) GetStepDescription() string {
	return fmt.Sprintf("Creating a new FQDN from the IP of a loadbalancer service")
}

// PerformSetupStep creates a loadbalancer service and sets the loadbalancer IP as the new FQDN.
func (fcs *fqdnCreatorStep) PerformSetupStep() error {
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

func (fcs *fqdnCreatorStep) checkIfNginxWillBeInstalled() error {
	for _, dogu := range fcs.config.Dogus.Install {
		if strings.Contains(dogu, "nginx-ingress") {
			return nil
		}
	}
	// nginx-ingress not found
	return fmt.Errorf("invalid configuration. FQDN can only be created with nginx-ingress installed")
}

func (fcs *fqdnCreatorStep) createServiceResource(ctx gocontext.Context) error {
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
			Selector:       map[string]string{DoguLabelName: "nginx-ingress"},
			Ports:          []corev1.ServicePort{createNginxPortResource(80), createNginxPortResource(443)},
		},
	}

	logrus.Debug("Create load balancer service")
	_, err := fcs.clientSet.CoreV1().Services(fcs.namespace).Create(ctx, serviceResource, metav1.CreateOptions{})
	return err
}

func createNginxPortResource(port int) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       fmt.Sprintf("nginx-ingress-%d", port),
		Protocol:   "TCP",
		Port:       int32(port),
		TargetPort: intstr.FromInt(port),
	}
}

func (fcs *fqdnCreatorStep) setFQDNFromLoadbalancerIP(ctx gocontext.Context) error {
	return retry.OnError(backoff, serviceRetry, func() error {
		logrus.Debug("Try retrieving service...")
		service, err := fcs.clientSet.CoreV1().Services(fcs.namespace).Get(ctx, cesLoadbalancerName, metav1.GetOptions{})
		if err != nil || len(service.Status.LoadBalancer.Ingress) <= 0 {
			logrus.Debugf("wait for service %s to be instantiated", cesLoadbalancerName)
			return fmt.Errorf("service not yet ready %s: %w", cesLoadbalancerName, err)
		}

		loadbalancerIP := service.Status.LoadBalancer.Ingress[0].IP
		fcs.config.Naming.Fqdn = loadbalancerIP
		logrus.Infof("Loadbalancer Ip succesfully retrieved and set as new FQDN")
		return nil
	})
}

func serviceRetry(err error) bool {
	return strings.Contains(err.Error(), "service not yet ready")
}
