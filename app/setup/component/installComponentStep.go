package component

import (
	"context"
	"fmt"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appcontext "github.com/cloudogu/k8s-ces-setup/v4/app/context"
)

type installComponentStep struct {
	client              componentsClient
	componentName       string
	componentNamespace  string
	version             string
	namespace           string
	deployNamespace     string
	valuesYamlOverwrite string
}

// NewInstallComponentStep creates a new step responsible to apply a component resource to the cluster, and, thus, starting the component installation.
func NewInstallComponentStep(client componentsClient, componentName string, attributes appcontext.ComponentAttributes, namespace string) *installComponentStep {
	return &installComponentStep{
		client:              client,
		componentName:       componentName,
		componentNamespace:  attributes.HelmRepositoryNamespace,
		version:             attributes.Version,
		namespace:           namespace,
		deployNamespace:     attributes.DeployNamespace,
		valuesYamlOverwrite: attributes.ValuesYamlOverwrite,
	}
}

// GetStepDescription return the human-readable description of the step
func (ics *installComponentStep) GetStepDescription() string {
	return fmt.Sprintf("Installing component '%s/%s:%s'", ics.componentNamespace, ics.componentName, ics.version)
}

// PerformSetupStep applies a component resource for the configured component to the cluster.
func (ics *installComponentStep) PerformSetupStep(ctx context.Context) error {
	cr := ics.createComponentCr()
	_, err := ics.client.Create(ctx, cr, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply component '%s/%s:%s' : %w", ics.componentNamespace, ics.componentName, ics.version, err)
	}

	return nil
}

func (ics *installComponentStep) createComponentCr() *v1.Component {
	cr := &v1.Component{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["app.kubernetes.io/name"] = ics.componentName
	cr.Name = ics.componentName
	cr.Namespace = ics.namespace
	cr.Spec.Name = ics.componentName
	cr.Spec.Namespace = ics.componentNamespace
	cr.Spec.Version = patchHelmChartVersion(ics.version)
	cr.Spec.DeployNamespace = ics.deployNamespace
	cr.Spec.ValuesYamlOverwrite = ics.valuesYamlOverwrite
	cr.Labels = labels

	return cr
}

func patchHelmChartVersion(version string) string {
	if version == "latest" {
		return ""
	}

	return version
}
