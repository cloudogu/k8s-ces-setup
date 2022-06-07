package component

import (
	"fmt"
	"github.com/cloudogu/k8s-apply-lib/apply"
	k8sv1 "github.com/cloudogu/k8s-dogu-operator/api/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
)

func NewServiceDiscoveryInstallerStep(clusterConfig *rest.Config, setupCtx *context.SetupContext) (*serviceDiscoveryInstallerStep, error) {
	k8sApplyClient, scheme, err := apply.New(clusterConfig, "TODO")
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s apply client: %w", err)
	}
	err = k8sv1.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("failed add applier scheme to dogu CRD scheme handling: %w", err)
	}

	return &serviceDiscoveryInstallerStep{
		namespace:              setupCtx.AppConfig.TargetNamespace,
		resourceURL:            setupCtx.AppConfig.ServiceDiscoveryURL,
		fileClient:             core.NewFileClient(setupCtx.AppVersion),
		k8sClient:              k8sApplyClient,
		fileContentModificator: &defaultFileContentModificator{},
	}, nil
}

type serviceDiscoveryInstallerStep struct {
	namespace              string
	resourceURL            string
	fileClient             fileClient
	fileContentModificator fileContentModificator
	k8sClient              k8sClient
}

// GetStepDescription returns a human-readable description of the service discovery installation step.
func (sdis *serviceDiscoveryInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install service discovery from %s", sdis.resourceURL)
}

// PerformSetupStep installs the service discovery.
func (sdis *serviceDiscoveryInstallerStep) PerformSetupStep() error {
	fileContent, err := sdis.fileClient.Get(sdis.resourceURL)
	if err != nil {
		return err
	}

	mod := sdis.fileContentModificator

	fileContent = mod.replaceNamespacedResources(fileContent, sdis.namespace)
	fileContent = mod.removeLegacyNamespaceFromResources(fileContent)

	sections := splitYamlFileSections(fileContent)

	err = sdis.applyYamlSections(sections)
	if err != nil {
		return err
	}

	return nil
}

func (sdis *serviceDiscoveryInstallerStep) applyYamlSections(sections [][]byte) error {
	for _, section := range sections {
		err := sdis.k8sClient.Apply(section, sdis.namespace)
		if err != nil {
			return err
		}
	}
	return nil
}
