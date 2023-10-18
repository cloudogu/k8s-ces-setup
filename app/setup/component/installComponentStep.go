package component

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type componentsClient interface {
	ecosystem.ComponentInterface
}

type installComponentStep struct {
	client             componentsClient
	componentName      string
	componentNamespace string
	version            string
	namespace          string
	deployNamespace    string
}

// NewInstallComponentStep creates a new step responsible to apply a component resource to the cluster, and, thus, starting the component installation.
func NewInstallComponentStep(client componentsClient, componentName string, componentNamespace string, version string, namespace string, deployNamespace string) *installComponentStep {
	return &installComponentStep{client: client, componentName: componentName, componentNamespace: componentNamespace, version: version, namespace: namespace, deployNamespace: deployNamespace}
}

// GetStepDescription return the human-readable description of the step
func (ics *installComponentStep) GetStepDescription() string {
	return fmt.Sprintf("Installing component '%s/%s:%s'", ics.componentNamespace, ics.componentName, ics.version)
}

// PerformSetupStep applies a component resource for the configured component to the cluster.
func (ics *installComponentStep) PerformSetupStep(ctx context.Context) error {
	cr := getComponentCr(ics.componentName, ics.componentNamespace, ics.version, ics.namespace, ics.deployNamespace)
	_, err := ics.client.Create(ctx, cr, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply component '%s/%s:%s' : %w", ics.componentNamespace, ics.componentName, ics.version, err)
	}

	return nil
}

func getComponentCr(componentName string, componentNameSpace string, version string, k8sNamespace string, deployNamespace string) *v1.Component {
	cr := &v1.Component{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["app.kubernetes.io/name"] = componentName
	cr.Name = componentName
	cr.Namespace = k8sNamespace
	cr.Spec.Name = componentName
	cr.Spec.Namespace = componentNameSpace
	cr.Spec.Version = patchHelmChartVersion(version)
	cr.Spec.DeployNamespace = deployNamespace
	cr.Labels = labels

	return cr
}

func patchHelmChartVersion(version string) string {
	if version == "latest" {
		return ""
	}

	return version
}
