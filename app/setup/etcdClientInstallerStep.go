package setup

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
)

func newEtcdClientInstallerStep(clientSet kubernetes.Interface, etcdClientVersion string) *etcdClientInstallerStep {
	return &etcdClientInstallerStep{ClientSet: clientSet, Version: etcdClientVersion}
}

type etcdClientInstallerStep struct {
	ClientSet kubernetes.Interface
	Version   string
}

func (ecis *etcdClientInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd client v%s", ecis.Version)
}

func (ecis *etcdClientInstallerStep) PerformSetupStep() error {
	panic("implement me")
}
