package component

import (
	"fmt"
	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
)

// NewStaticWebserverInstallerStep create a new instance of an installer step which installs a static webserver.
func NewStaticWebserverInstallerStep(setupCtx *ctx.SetupContext, k8sClient k8sClient) (*staticWebserverInstallerStep, error) {
	return &staticWebserverInstallerStep{
		namespace:   setupCtx.AppConfig.TargetNamespace,
		resourceURL: setupCtx.AppConfig.StaticWebserverURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   k8sClient,
	}, nil
}

type staticWebserverInstallerStep struct {
	namespace   string
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the static webserver installation step.
func (esis *staticWebserverInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install static webserver from %s", esis.resourceURL)
}

// PerformSetupStep installs the static webserver.
func (esis *staticWebserverInstallerStep) PerformSetupStep() error {
	fileContent, err := esis.fileClient.Get(esis.resourceURL)
	if err != nil {
		return err
	}

	err = applyNamespacedYamlSection(esis.k8sClient, fileContent, esis.namespace)
	if err != nil {
		return err
	}

	return nil
}
