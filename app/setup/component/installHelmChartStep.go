package component

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	helmclient "github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/labels"
)

type helmClient interface {
	// InstallOrUpgrade takes a component and applies the corresponding helmChart.
	InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error
	// GetLatestVersion tries to get the latest version identifier for the chart with the given name.
	GetLatestVersion(chartName string) (string, error)
}

type installHelmChartStep struct {
	namespace  string
	chart      string
	helmClient helmClient
}

// NewInstallHelmChartStep creates new instance of a k8s component chart
func NewInstallHelmChartStep(namespace string, chartUrl string, helmClient helmClient) *installHelmChartStep {
	return &installHelmChartStep{
		namespace:  namespace,
		chart:      chartUrl,
		helmClient: helmClient,
	}
}

// GetStepDescription returns a human-readable description of the component-chart installation step.
func (s *installHelmChartStep) GetStepDescription() string {
	return fmt.Sprintf("Install component-chart from %s in namespace %s", s.chart, s.namespace)
}

// PerformSetupStep installs the component chart.
func (s *installHelmChartStep) PerformSetupStep(ctx context.Context) error {
	if s.chart == "" {
		return fmt.Errorf("error install component chart: chart url is empty")
	}

	return s.installChart(ctx, s.chart)
}

func (s *installHelmChartStep) installChart(ctx context.Context, chart string) error {
	fullChartName, chartVersion, err := SplitChartString(chart)
	if err != nil {
		return err
	}

	releaseName := fullChartName[strings.LastIndex(fullChartName, "/")+1:]
	if len(releaseName) <= 0 {
		return fmt.Errorf("error reading chartname '%s': wrong format", fullChartName)
	}

	if chartVersion == "latest" {
		chartVersion, err = s.helmClient.GetLatestVersion(fullChartName)
		if err != nil {
			return fmt.Errorf("error fetching latest version of chart %q: %w", fullChartName, err)
		}
	}

	chartSpec := s.createChartSpec(releaseName, fullChartName, chartVersion)

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

func (s *installHelmChartStep) createChartSpec(releaseName string, fullChartName string, chartVersion string) *helmclient.ChartSpec {
	return &helmclient.ChartSpec{
		ReleaseName: releaseName,
		ChartName:   fullChartName,
		Namespace:   s.namespace,
		Version:     chartVersion,
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// Wait for the release to deployed and ready
		Atomic:          true,
		CreateNamespace: true,
		PostRenderer: labels.NewPostRenderer(map[string]string{
			v1.ComponentNameLabelKey:    releaseName,
			v1.ComponentVersionLabelKey: chartVersion,
		}),
	}
}
