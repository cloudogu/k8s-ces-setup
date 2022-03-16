package setup

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
)

func newEtcdInstallerStep(clientSet kubernetes.Interface, etcdVersion string) *etcdInstallerStep {
	return &etcdInstallerStep{ClientSet: clientSet, Version: etcdVersion}
}

type etcdInstallerStep struct {
	ClientSet kubernetes.Interface
	Version   string
}

func (eis *etcdInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd server v%s", eis.Version)
}

func (eis *etcdInstallerStep) PerformSetupStep() error {
	panic("implement me")
}
