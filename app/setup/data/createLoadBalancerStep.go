package data

import (
	gocontext "context"
	"fmt"
	"strings"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	cesLoadbalancerName = "ces-loadbalancer"

	// DoguLabelName is used to select a dogu pod by name.
	DoguLabelName    = "dogu.name"
	nginxIngressName = "nginx-ingress"
)

type createLoadBalancerStep struct {
	config    *context.SetupJsonConfiguration
	clientSet kubernetes.Interface
	namespace string
}

// NewCreateLoadBalancerStep creates the external loadbalancer service for the Cloudogu EcoSystem
func NewCreateLoadBalancerStep(config *context.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *createLoadBalancerStep {
	return &createLoadBalancerStep{config: config, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *createLoadBalancerStep) GetStepDescription() string {
	return "Retrieving a new FQDN from the IP of a loadbalancer service"
}

// PerformSetupStep creates a loadbalancer service and sets the loadbalancer IP as the new FQDN.
func (fcs *createLoadBalancerStep) PerformSetupStep() error {
	err := fcs.checkIfNginxWillBeInstalled()
	if err != nil {
		return err
	}

	ctx := gocontext.Background()
	return fcs.createServiceResource(ctx)
}

func (fcs *createLoadBalancerStep) checkIfNginxWillBeInstalled() error {
	for _, dogu := range fcs.config.Dogus.Install {
		if strings.Contains(dogu, nginxIngressName) {
			return nil
		}
	}
	return fmt.Errorf("invalid configuration: FQDN can only be created if nginx-ingress will be installed")
}

func (fcs *createLoadBalancerStep) createServiceResource(ctx gocontext.Context) error {
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
