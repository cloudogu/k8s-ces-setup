package component

import (
	"context"
	"fmt"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	helmclient "github.com/cloudogu/k8s-component-operator/pkg/helm/client"
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
	crdChart   string
}

// NewComponentOperatorInstallerStep creates new instance of the component-operator
func NewComponentOperatorInstallerStep(setupCtx *appcontext.SetupContext, helmClient helmClient) *componentOperatorInstallerStep {
	return &componentOperatorInstallerStep{
		namespace:  setupCtx.AppConfig.TargetNamespace,
		chart:      setupCtx.AppConfig.ComponentOperatorChart,
		crdChart:   setupCtx.AppConfig.ComponentOperatorCrdChart,
		helmClient: helmClient,
	}
}

// GetStepDescription returns a human-readable description of the component-operator installation step.
func (cois *componentOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install component-operator from %s and component-crd from %s", cois.chart, cois.crdChart)
}

// PerformSetupStep installs the dogu operator.
func (cois *componentOperatorInstallerStep) PerformSetupStep(ctx context.Context) error {
	err := cois.installChart(ctx, cois.crdChart)
	if err != nil {
		return err
	}
	return cois.installChart(ctx, cois.chart)
}

func (cois *componentOperatorInstallerStep) installChart(ctx context.Context, chart string) error {
	fullChartName, chartVersion, err := splitChartString(chart)
	if err != nil {
		return err
	}

	chartName := fullChartName[strings.LastIndex(fullChartName, "/")+1:]
	if len(chartName) <= 0 {
		return fmt.Errorf("error reading chartname '%s': wrong format", fullChartName)
	}

	chartSpec := cois.createChartSpec(chartName, fullChartName, chartVersion)

	return cois.helmClient.InstallOrUpgrade(ctx, chartSpec)
}

func splitChartString(chart string) (string, string, error) {
	chartSplit := strings.Split(chart, ":")
	if len(chartSplit) != 2 {
		return "", "", fmt.Errorf("componentChart '%s' has a wrong format. Must be '<chartName>:<version>'; e.g.: 'foo/bar:1.2.3'", chart)
	}

	fullChartName := chartSplit[0]
	chartVersion := chartSplit[1]
	return fullChartName, chartVersion, nil
}

func (cois *componentOperatorInstallerStep) createChartSpec(chartName string, fullChartName string, chartVersion string) *helmclient.ChartSpec {
	return &helmclient.ChartSpec{
		ReleaseName: chartName,
		ChartName:   fullChartName,
		Namespace:   cois.namespace,
		Version:     patchHelmChartVersion(chartVersion),
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// Wait for the release to deployed and ready
		Atomic: true,
	}
}
