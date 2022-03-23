package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func newEtcdServerInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *etcdServerInstallerStep {
	return &etcdServerInstallerStep{
		namespace:   setupCtx.AppConfig.Namespace,
		resourceURL: setupCtx.AppConfig.EtcdServerResourceURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
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

	err = esis.k8sClient.Apply(fileContent, esis.namespace)
	if err != nil {
		return err
	}

	return nil
}
