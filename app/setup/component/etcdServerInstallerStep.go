package component

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
)

func NewEtcdServerInstallerStep(setupCtx *context.SetupContext, k8sClient k8sClient) (*etcdServerInstallerStep, error) {
	return &etcdServerInstallerStep{
		namespace:   setupCtx.AppConfig.TargetNamespace,
		resourceURL: setupCtx.AppConfig.EtcdServerResourceURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   k8sClient,
	}, nil
}

type etcdServerInstallerStep struct {
	namespace   string
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the etcd installation step.
func (esis *etcdServerInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd server from %s", esis.resourceURL)
}

// PerformSetupStep installs the CES etcd.
func (esis *etcdServerInstallerStep) PerformSetupStep() error {
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
