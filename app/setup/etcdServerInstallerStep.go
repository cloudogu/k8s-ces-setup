package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func newEtcdServerInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *etcdServerInstallerStep {
	return &etcdServerInstallerStep{
		namespace:   setupCtx.AppConfig.TargetNamespace,
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

	// avoid extra namespace validation: during the namespace creation earlier it was already validated by the K8s API
	fileContent = replaceNamespacedResources(fileContent, esis.namespace)
	fileContent = removeLegacyNamespaceFromResources(fileContent)

	sections := splitYamlFileSections(fileContent)

	for _, section := range sections {
		err = esis.k8sClient.Apply(section, esis.namespace)
		if err != nil {
			return err
		}
	}

	return nil
}
