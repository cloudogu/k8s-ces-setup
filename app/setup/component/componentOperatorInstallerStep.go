package component

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	helmclient "github.com/mittwald/go-helm-client"
	"strings"
	"time"
)

type helmClient interface {
	// InstallOrUpgrade takes a component and applies the corresponding helmChart.
	InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error
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

	fullChartName := chartSplit[0]
	chartVersion := chartSplit[1]

	chartName := fullChartName[strings.LastIndex(fullChartName, "/")+1:]
	if len(chartName) <= 0 {
		return fmt.Errorf("error reading chartname '%s': wrong format", fullChartName)
	}

	chartSpec := &helmclient.ChartSpec{
		ReleaseName: chartName,
		ChartName:   fullChartName,
		Namespace:   cois.namespace,
		Version:     chartVersion,
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// Wait for the release to deployed and ready
		Wait: true,
	}

	return cois.helmClient.InstallOrUpgrade(ctx, chartSpec)
}
