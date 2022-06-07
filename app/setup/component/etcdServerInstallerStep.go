package component

import (
	"fmt"
	"github.com/cloudogu/k8s-apply-lib/apply"
	k8sv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func NewEtcdServerInstallerStep(clusterConfig *rest.Config, setupCtx *context.SetupContext) (*etcdServerInstallerStep, error) {
	k8sApplyClient, scheme, err := apply.New(clusterConfig, K8sSetupFieldManagerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s apply client: %w", err)
	}
	err = k8sv1.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("failed add applier scheme to dogu CRD scheme handling: %w", err)
	}
	return &etcdServerInstallerStep{
		namespace:              setupCtx.AppConfig.TargetNamespace,
		resourceURL:            setupCtx.AppConfig.EtcdServerResourceURL,
		fileClient:             core.NewFileClient(setupCtx.AppVersion),
		k8sClient:              k8sApplyClient,
		fileContentModificator: &defaultFileContentModificator{},
	}, nil
}

type etcdServerInstallerStep struct {
	namespace              string
	resourceURL            string
	fileClient             fileClient
	k8sClient              k8sClient
	fileContentModificator fileContentModificator
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

	mod := esis.fileContentModificator

	fileContent = mod.replaceNamespacedResources(fileContent, esis.namespace)
	fileContent = mod.removeLegacyNamespaceFromResources(fileContent)

	sections := splitYamlFileSections(fileContent)

	err = esis.applyYamlSections(sections)
	if err != nil {
		return err
	}

	return nil
}

func (esis *etcdServerInstallerStep) applyYamlSections(sections [][]byte) error {
	for _, section := range sections {
		err := esis.k8sClient.Apply(section, esis.namespace)
		if err != nil {
			return err
		}
	}
	return nil
}
