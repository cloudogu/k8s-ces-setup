package setup

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
)

func newDoguOperatorInstallerStep(clientSet kubernetes.Interface, doguOperatorVersion string) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{ClientSet: clientSet, Version: doguOperatorVersion}
}

type doguOperatorInstallerStep struct {
	ClientSet kubernetes.Interface
	Version   string
}

func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator version %s", dois.Version)
}

func (dois *doguOperatorInstallerStep) PerformSetupStep() error {
	panic("implement me")
}
