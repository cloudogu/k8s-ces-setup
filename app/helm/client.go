package helm

import (
	"context"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

const (
	helmRepositoryCache    = "/tmp/.helmcache"
	helmRepositoryConfig   = "/tmp/.helmrepo"
	helmRegistryConfigFile = "/tmp/.helmregistry/config.json"
)

// HelmClient embeds the helmclient.Client interface for usage in this package.
type HelmClient interface {
	helmclient.Client
}

// Client wraps the HelmClient with config.HelmRepositoryData
type Client struct {
	helmClient          HelmClient
	helmRepoOciEndpoint string
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoOciEndpoint string, debug bool, debugLog action.DebugLog) (*Client, error) {
	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        namespace,
			RepositoryCache:  helmRepositoryCache,
			RepositoryConfig: helmRepositoryConfig,
			RegistryConfig:   helmRegistryConfigFile,
			Debug:            debug,
			DebugLog:         debugLog,
			Linting:          true,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	return &Client{helmClient: helmClient, helmRepoOciEndpoint: helmRepoOciEndpoint}, nil
}

// InstallOrUpgradeChart uses Helm to install the given chart in the given namespace.
func (c *Client) InstallOrUpgradeChart(ctx context.Context, namespace string, chart string, version string) error {
	chartName := chart[strings.LastIndex(chart, "/")+1:]
	if len(chartName) <= 0 {
		return fmt.Errorf("error reading chartname '%s': wrong format", chart)
	}

	chartSpec := &helmclient.ChartSpec{
		ReleaseName: chartName,
		ChartName:   fmt.Sprintf("%s/%s", c.helmRepoOciEndpoint, chart),
		Namespace:   namespace,
		Version:     version,
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// Wait for the release to deployed and ready
		Wait: true,
	}

	_, err := c.helmClient.InstallOrUpgradeChart(ctx, chartSpec, nil)
	if err != nil {
		return fmt.Errorf("error while installing chart %s: %w", chart, err)
	}
	return nil
}
