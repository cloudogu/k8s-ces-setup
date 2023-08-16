package component

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"strings"
)

type helmClient interface {
	// InstallOrUpgradeChart uses Helm to install the given chart in the given namespace.
	InstallOrUpgradeChart(ctx context.Context, namespace string, chart string, version string) error
}

type componentOperatorInstallerStep struct {
	namespace  string
	chart      string
	helmClient helmClient
}

// NewComponentOperatorInstallerStep creates new instance of the component-operator
func NewComponentOperatorInstallerStep(setupCtx *appcontext.SetupContext, helmClient helmClient) *componentOperatorInstallerStep {
	return &componentOperatorInstallerStep{
		namespace:  setupCtx.AppConfig.TargetNamespace,
		chart:      setupCtx.AppConfig.ComponentOperatorChart,
		helmClient: helmClient,
	}
}

// GetStepDescription returns a human-readable description of the component-operator installation step.
func (cois *componentOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install component-operator from %s", cois.chart)
}

// PerformSetupStep installs the dogu operator.
func (cois *componentOperatorInstallerStep) PerformSetupStep(ctx context.Context) error {
	chartSplit := strings.Split(cois.chart, ":")
	if len(chartSplit) != 2 {
		return fmt.Errorf("componentOperatorChart '%s' has a wrong format. Must be '<chartName>:<version>'; e.g.: 'foo/bar:1.2.3'", cois.chart)
	}

	return cois.helmClient.InstallOrUpgradeChart(ctx, cois.namespace, chartSplit[0], chartSplit[1])
}
