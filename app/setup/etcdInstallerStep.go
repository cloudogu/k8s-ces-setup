package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func newEtcdInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *etcdInstallerStep {
	return &etcdInstallerStep{
		resourceURL: setupCtx.AppConfig.DoguOperatorURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
}

type etcdInstallerStep struct {
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the etcd installation step.
func (eis *etcdInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd server from %s", eis.resourceURL)
}

// PerformSetupStep installs the CES etcd.
func (eis *etcdInstallerStep) PerformSetupStep() error {
	fileContent, err := eis.fileClient.Get(eis.resourceURL)
	if err != nil {
		return err
	}

	err = eis.k8sClient.Apply(fileContent)
	if err != nil {
		return err
	}

	return nil
}
