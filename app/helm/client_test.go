package helm

import (
	"context"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("should create new client", func(t *testing.T) {
		namespace := "ecosystem"

		// override default controller method to retrieve a kube config
		oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
		defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{}
		}

		client, err := NewClient(namespace, "http://helm.repo", false, nil)

		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestClient_InstallOrUpgradeChart(t *testing.T) {
	t.Run("should install or upgrade chart", func(t *testing.T) {
		helmRepo := "oci://staging.cloudogu.com"
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   fmt.Sprintf("%s/%s", helmRepo, "testing/testComponent"),
			Namespace:   "ecosystem",
			Version:     "0.1.1",
			Timeout:     time.Second * 300,
			Wait:        true,
		}
		ctx := context.TODO()

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chartSpec, mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoOciEndpoint: helmRepo}

		err := client.InstallOrUpgradeChart(ctx, "ecosystem", "testing/testComponent", "0.1.1")

		require.NoError(t, err)
	})

	t.Run("should fail to install or upgrade chart for error in chart-string", func(t *testing.T) {
		helmRepo := "oci://staging.cloudogu.com"
		ctx := context.TODO()

		mockHelmClient := NewMockHelmClient(t)

		client := &Client{helmClient: mockHelmClient, helmRepoOciEndpoint: helmRepo}

		err := client.InstallOrUpgradeChart(ctx, "ecosystem", "", "0.1.1")

		require.Error(t, err)
		assert.ErrorContains(t, err, "error reading chartname '': wrong format")
	})

	t.Run("should fail to install or upgrade chart for error in helmClient", func(t *testing.T) {
		helmRepo := "oci://staging.cloudogu.com"
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   fmt.Sprintf("%s/%s", helmRepo, "testing/testComponent"),
			Namespace:   "ecosystem",
			Version:     "0.1.1",
			Timeout:     time.Second * 300,
			Wait:        true,
		}
		ctx := context.TODO()

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chartSpec, mock.Anything).Return(nil, assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoOciEndpoint: helmRepo}

		err := client.InstallOrUpgradeChart(ctx, "ecosystem", "testing/testComponent", "0.1.1")

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while installing chart testing/testComponent")
	})
}
