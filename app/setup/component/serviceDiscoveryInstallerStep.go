package component

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
)

func NewServiceDiscoveryInstallerStep(setupCtx *context.SetupContext, k8sClient k8sClient) (*serviceDiscoveryInstallerStep, error) {
	return &serviceDiscoveryInstallerStep{
		namespace:              setupCtx.AppConfig.TargetNamespace,
		resourceURL:            setupCtx.AppConfig.ServiceDiscoveryURL,
		resourceRegistryClient: core.NewResourceRegistryClient(setupCtx.AppVersion, setupCtx.DoguRegistrySecret()),
		k8sClient:              k8sClient,
	}, nil
}

type serviceDiscoveryInstallerStep struct {
	namespace              string
	resourceURL            string
	resourceRegistryClient resourceRegistryClient
	k8sClient              k8sClient
}

// GetStepDescription returns a human-readable description of the service discovery installation step.
func (sdis *serviceDiscoveryInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install service discovery from %s", sdis.resourceURL)
}

// PerformSetupStep installs the service discovery.
func (sdis *serviceDiscoveryInstallerStep) PerformSetupStep() error {
	fileContent, err := sdis.resourceRegistryClient.GetResourceFileContent(sdis.resourceURL)
	if err != nil {
		return err
	}

	err = applyNamespacedYamlSection(sdis.k8sClient, fileContent, sdis.namespace)
	if err != nil {
		return err
	}

	return nil
}
