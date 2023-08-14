package component

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-apply-lib/apply"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
)

type resourceRegistryClient interface {
	// GetResourceFileContent retrieves a file identified by its URL and returns the contents.
	GetResourceFileContent(resourceURL string) ([]byte, error)
}

type k8sClient interface {
	// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster's namespace.
	Apply(yamlResources apply.YamlDocument, namespace string) error
	// ApplyWithOwner provides a testable method
	ApplyWithOwner(doc apply.YamlDocument, namespace string, resource metav1.Object) error
}

type doguOperatorInstallerStep struct {
	namespace              string
	resourceURL            string
	k8sClient              k8sClient
	resourceRegistryClient resourceRegistryClient
}

// NewDoguOperatorInstallerStep creates new instance of the dogu operator and creates an unversioned client for apply dogu cr's
func NewDoguOperatorInstallerStep(setupCtx *appcontext.SetupContext, k8sClient k8sClient) (*doguOperatorInstallerStep, error) {
	return &doguOperatorInstallerStep{
		namespace:              setupCtx.AppConfig.TargetNamespace,
		resourceURL:            setupCtx.AppConfig.DoguOperatorURL,
		k8sClient:              k8sClient,
		resourceRegistryClient: core.NewResourceRegistryClient(setupCtx.AppVersion, setupCtx.DoguRegistryConfiguration),
	}, nil
}

// GetStepDescription returns a human-readable description of the dogu operator installation step.
func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator from %s", dois.resourceURL)
}

// PerformSetupStep installs the dogu operator.
func (dois *doguOperatorInstallerStep) PerformSetupStep(context.Context) error {
	fileContent, err := dois.resourceRegistryClient.GetResourceFileContent(dois.resourceURL)
	if err != nil {
		return err
	}

	err = applyNamespacedYamlSection(dois.k8sClient, fileContent, dois.namespace)
	if err != nil {
		return err
	}

	return nil
}

func applyNamespacedYamlSection(k8sClient k8sClient, fileContent []byte, namespace string) error {
	namespaceTemplate := struct {
		Namespace string
	}{
		Namespace: namespace,
	}

	err := apply.NewBuilder(k8sClient).
		WithYamlResource("input", fileContent).
		WithNamespace(namespace).
		WithTemplate("input", namespaceTemplate).
		ExecuteApply()

	return err
}
