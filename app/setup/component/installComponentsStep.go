package component

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ecoSystemClient interface {
	Components(namespace string) ecosystem.ComponentInterface
}

type installComponentsStep struct {
	client             ecoSystemClient
	componentName      string
	componentNamespace string
	version            string
	namespace          string
}

// NewInstallComponentsStep creates a new step responsible to apply a component resource to the cluster, and, thus, starting the component installation.
func NewInstallComponentsStep(client ecoSystemClient, componentName string, componentNamespace string, version string, namespace string) *installComponentsStep {
	return &installComponentsStep{client: client, componentName: componentName, componentNamespace: componentNamespace, version: version, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (ics *installComponentsStep) GetStepDescription() string {
	return fmt.Sprintf("Installing component [%s:%s]", ics.componentName, ics.version)
}

// PerformSetupStep applies a component resource for the configured component to the cluster.
func (ics *installComponentsStep) PerformSetupStep(ctx context.Context) error {
	cr := getComponentCr(ics.componentName, ics.componentNamespace, ics.version, ics.namespace)
	_, err := ics.client.Components(cr.Namespace).Create(ctx, cr, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply component '%s:%s' : %w", ics.componentName, ics.version, err)
	}

	return nil
}

func getComponentCr(componentName string, componentNameSpace string, version string, k8sNamespace string) *v1.Component {
	cr := &v1.Component{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["app.kubernetes.io/name"] = componentName
	cr.Name = componentName
	cr.Namespace = k8sNamespace
	cr.Spec.Name = componentName
	cr.Spec.Namespace = componentNameSpace
	cr.Spec.Version = version
	cr.Labels = labels

	return cr
}
