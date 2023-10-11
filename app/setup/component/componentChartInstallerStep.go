package component

import (
	"context"
	"fmt"
	helmclient "github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"strings"
	"time"
)

type helmClient interface {
	// InstallOrUpgrade takes a component and applies the corresponding helmChart.
	InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error
}

type cesComponentChartInstallerStep struct {
	namespace  string
	chart      string
	helmClient helmClient
}

// NewCesComponentChartInstallerStep creates new instance of a k8s component chart
func NewCesComponentChartInstallerStep(namespace string, chartUrl string, helmClient helmClient) *cesComponentChartInstallerStep {
	return &cesComponentChartInstallerStep{
		namespace:  namespace,
		chart:      chartUrl,
		helmClient: helmClient,
	}
}

// GetStepDescription returns a human-readable description of the component-chart installation step.
func (s *cesComponentChartInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install component-chart from %s in namespace %s", s.chart, s.namespace)
}

// PerformSetupStep installs the component chart.
func (s *cesComponentChartInstallerStep) PerformSetupStep(ctx context.Context) error {
	if s.chart == "" {
		return fmt.Errorf("error install component chart: chart url is empty")
	}

	return s.installChart(ctx, s.chart)
}

func (s *cesComponentChartInstallerStep) installChart(ctx context.Context, chart string) error {
	fullChartName, chartVersion, err := SplitChartString(chart)
	if err != nil {
		return err
	}

	chartName := fullChartName[strings.LastIndex(fullChartName, "/")+1:]
	if len(chartName) <= 0 {
		return fmt.Errorf("error reading chartname '%s': wrong format", fullChartName)
	}

	chartSpec := s.createChartSpec(chartName, fullChartName, chartVersion)

	return s.helmClient.InstallOrUpgrade(ctx, chartSpec)
}

func SplitChartString(chart string) (string, string, error) {
	chartSplit := strings.Split(chart, ":")
	if len(chartSplit) != 2 {
		return "", "", fmt.Errorf("componentChart '%s' has a wrong format. Must be '<chartName>:<version>'; e.g.: 'foo/bar:1.2.3'", chart)
	}

	fullChartName := chartSplit[0]
	chartVersion := chartSplit[1]
	return fullChartName, chartVersion, nil
}

func SplitHelmNamespaceFromChartString(chartString string) (string, string) {
	split := strings.Split(chartString, "/")
	return split[0], split[1]
}

func (s *cesComponentChartInstallerStep) createChartSpec(chartName string, fullChartName string, chartVersion string) *helmclient.ChartSpec {
	return &helmclient.ChartSpec{
		ReleaseName: chartName,
		ChartName:   fullChartName,
		Namespace:   s.namespace,
		Version:     patchHelmChartVersion(chartVersion),
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// Wait for the release to deployed and ready
		Atomic:          true,
		CreateNamespace: true,
	}
}
