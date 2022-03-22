package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func newDoguOperatorInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{
		resourceURL: setupCtx.AppConfig.DoguOperatorURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
}

type doguOperatorInstallerStep struct {
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the dogu operator installation step.
func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator from %s", dois.resourceURL)
}

// PerformSetupStep installs the dogu operator.
func (dois *doguOperatorInstallerStep) PerformSetupStep() error {
	fileContent, err := dois.fileClient.Get(dois.resourceURL)
	if err != nil {
		return err
	}

	err = dois.k8sClient.Apply(fileContent)
	if err != nil {
		return err
	}

	return nil
}
