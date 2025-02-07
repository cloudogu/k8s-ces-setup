package data

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"strings"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
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
	config    *appcontext.SetupJsonConfiguration
	clientSet kubernetes.Interface
	namespace string
}

// NewCreateLoadBalancerStep creates the external loadbalancer service for the Cloudogu EcoSystem
func NewCreateLoadBalancerStep(config *appcontext.SetupJsonConfiguration, clientSet kubernetes.Interface, namespace string) *createLoadBalancerStep {
	return &createLoadBalancerStep{config: config, clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *createLoadBalancerStep) GetStepDescription() string {
	return "Creating the main loadbalancer service for the Cloudogu EcoSystem"
}

// PerformSetupStep creates a loadbalancer service and sets the loadbalancer IP as the new FQDN.
func (fcs *createLoadBalancerStep) PerformSetupStep(ctx context.Context) error {
	err := fcs.checkIfNginxWillBeInstalled()
	if err != nil {
		return err
	}

	return fcs.upsertServiceResource(ctx)
}

func (fcs *createLoadBalancerStep) checkIfNginxWillBeInstalled() error {
	for _, dogu := range fcs.config.Dogus.Install {
		if strings.Contains(dogu, nginxIngressName) {
			return nil
		}
	}
	return fmt.Errorf("invalid configuration: FQDN can only be created if nginx-ingress will be installed")
}

func (fcs *createLoadBalancerStep) upsertServiceResource(ctx context.Context) error {
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

	get, err := fcs.clientSet.CoreV1().Services(fcs.namespace).Get(ctx, serviceResource.Name, metav1.GetOptions{})
	if err == nil {
		return fcs.updateServiceResource(ctx, get, serviceResource)
	}

	if err != nil && errors.IsNotFound(err) {
		logrus.Debug("Create load balancer service")
		_, err = fcs.clientSet.CoreV1().Services(fcs.namespace).Create(ctx, serviceResource, metav1.CreateOptions{})
	}

	return err
}

func (fcs *createLoadBalancerStep) updateServiceResource(ctx context.Context, actualService *corev1.Service, service *corev1.Service) error {
	logrus.Debug("Update existing load balancer service")
	// Update to ensure idempotence
	if actualService.Labels != nil {
		actualService.Labels["app"] = "ces"
	} else {
		actualService.Labels = service.Labels
	}
	actualService.Spec.Type = service.Spec.Type
	actualService.Spec.IPFamilyPolicy = service.Spec.IPFamilyPolicy
	actualService.Spec.IPFamilies = service.Spec.IPFamilies
	actualService.Spec.Selector = service.Spec.Selector
	actualService.Spec.Ports = service.Spec.Ports
	_, updateErr := fcs.clientSet.CoreV1().Services(fcs.namespace).Update(ctx, actualService, metav1.UpdateOptions{})
	return updateErr
}

func createNginxPortResource(port int) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       fmt.Sprintf("%s-%d", nginxIngressName, port),
		Protocol:   corev1.ProtocolTCP,
		Port:       int32(port),
		TargetPort: intstr.FromInt32(int32(port)),
	}
}
