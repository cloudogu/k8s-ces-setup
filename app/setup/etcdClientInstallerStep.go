package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func newEtcdClientInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *etcdClientInstallerStep {
	return &etcdClientInstallerStep{
		resourceURL: setupCtx.AppConfig.EtcdClientResourceURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
}

type etcdClientInstallerStep struct {
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the etcd client step.
func (ecis *etcdClientInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd client from %s", ecis.resourceURL)
}

// PerformSetupStep installs an etcd client.
func (ecis *etcdClientInstallerStep) PerformSetupStep() error {
	fileContent, err := ecis.fileClient.Get(ecis.resourceURL)
	if err != nil {
		return err
	}

	err = ecis.k8sClient.Apply(fileContent)
	if err != nil {
		return err
	}

	return nil
}
