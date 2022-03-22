package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

type fileClient interface {
	// Get retrieves a file identified by its URL and returns the contents.
	Get(url string) ([]byte, error)
}

type k8sClient interface {
	// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster.
	Apply(yamlResources []byte) error
}

func newDoguOperatorInstallerStep(clusterConfig *rest.Config, resourceURL, appVersion string) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{
		resourceURL: resourceURL,
		Version:     appVersion,
		fileClient:  core.NewFileClient(appVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
}

type doguOperatorInstallerStep struct {
	clusterConfig *rest.Config
	Version       string
	resourceURL   string
	fileClient    fileClient
	k8sClient     k8sClient
}

func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator version %s", dois.Version)
}

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
