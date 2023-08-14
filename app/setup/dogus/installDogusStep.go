package dogus

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
)

type ecoSystemClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
}

type installDogusStep struct {
	client    ecoSystemClient
	dogu      *core.Dogu
	namespace string
}

// NewInstallDogusStep creates a new step responsible to apply a dogu resource to the cluster, and, thus, starting the dogu installation.
func NewInstallDogusStep(client ecoSystemClient, dogu *core.Dogu, namespace string) *installDogusStep {
	return &installDogusStep{client: client, dogu: dogu, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (ids *installDogusStep) GetStepDescription() string {
	return fmt.Sprintf("Installing dogu [%s]", ids.dogu.GetFullName())
}

// PerformSetupStep applies a dogu recource for the configured dogu to the cluster.
func (ids *installDogusStep) PerformSetupStep(ctx context.Context) error {
	doguVersion, err := ids.dogu.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to get version from dogu [%s]: %w", ids.dogu.GetFullName(), err)
	}

	cr := getDoguCr(ids.dogu.GetSimpleName(), ids.dogu.GetFullName(), doguVersion.Raw, ids.namespace)
	_, err = ids.client.Dogus(ids.namespace).Create(ctx, cr, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply dogu %s: %w", ids.dogu.GetSimpleName(), err)
	}

	return nil
}

func getDoguCr(name string, namespaceName string, version string, k8sNamespace string) *v1.Dogu {
	cr := &v1.Dogu{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["dogu.name"] = name
	cr.Name = name
	cr.Namespace = k8sNamespace
	cr.Spec.Name = namespaceName
	cr.Spec.Version = version
	cr.Labels = labels

	return cr
}
