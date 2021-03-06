package dogus

import (
	gocontext "context"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"

	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"k8s.io/client-go/rest"
)

type installDogusStep struct {
	client    rest.Interface
	dogu      *core.Dogu
	namespace string
}

// NewInstallDogusStep creates a new step responsible to apply a dogu resource to the cluster, and, thus, starting the dogu instllation.
func NewInstallDogusStep(client rest.Interface, dogu *core.Dogu, namespace string) *installDogusStep {
	return &installDogusStep{client: client, dogu: dogu, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (ids *installDogusStep) GetStepDescription() string {
	return fmt.Sprintf("Installing dogu [%s]", ids.dogu.GetFullName())
}

// PerformSetupStep applies a dogu recource for the configured dogu to the cluster.
func (ids *installDogusStep) PerformSetupStep() error {
	doguVersion, err := ids.dogu.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to get version from dogu [%s]: %w", ids.dogu.GetFullName(), err)
	}

	cr := getDoguCr(ids.dogu.GetSimpleName(), ids.dogu.GetFullName(), doguVersion.Raw, ids.namespace)
	result := ids.client.Post().Namespace(ids.namespace).Resource("dogus").Body(cr).Do(gocontext.Background())
	err = result.Error()
	if err != nil {
		return fmt.Errorf("failed to apply dogu %s: %w", ids.dogu.GetSimpleName(), err)
	}

	return nil
}

func getDoguCr(name string, namespaceName string, version string, k8sNamespace string) *v1.Dogu {
	cr := &v1.Dogu{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["dogu"] = name
	cr.Name = name
	cr.Namespace = k8sNamespace
	cr.Spec.Name = namespaceName
	cr.Spec.Version = version
	cr.Labels = labels

	return cr
}
