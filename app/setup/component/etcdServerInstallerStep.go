package component

import (
	"context"
	"fmt"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
)

// NewEtcdServerInstallerStep creates a step to install the etcd server.
func NewEtcdServerInstallerStep(setupCtx *appcontext.SetupContext, k8sClient k8sClient) (*etcdServerInstallerStep, error) {
	return &etcdServerInstallerStep{
		namespace:              setupCtx.AppConfig.TargetNamespace,
		resourceURL:            setupCtx.AppConfig.EtcdServerResourceURL,
		resourceRegistryClient: core.NewResourceRegistryClient(setupCtx.AppVersion, setupCtx.DoguRegistryConfiguration),
		k8sClient:              k8sClient,
	}, nil
}

type etcdServerInstallerStep struct {
	namespace              string
	resourceURL            string
	resourceRegistryClient resourceRegistryClient
	k8sClient              k8sClient
}

// GetStepDescription returns a human-readable description of the etcd installation step.
func (esis *etcdServerInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd server from %s", esis.resourceURL)
}

// PerformSetupStep installs the CES etcd.
func (esis *etcdServerInstallerStep) PerformSetupStep(context.Context) error {
	fileContent, err := esis.resourceRegistryClient.GetResourceFileContent(esis.resourceURL)
	if err != nil {
		return err
	}

	err = applyNamespacedYamlSection(esis.k8sClient, fileContent, esis.namespace)
	if err != nil {
		return err
	}

	return nil
}
