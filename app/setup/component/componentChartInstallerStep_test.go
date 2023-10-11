package component

import (
	"context"
	"testing"
	"time"

	helmclient "github.com/cloudogu/k8s-component-operator/pkg/helm/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentChartInstallerStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)

		// when
		step := NewCesComponentChartInstallerStep("testNS", "testing/co:0.1", helmClientMock)

		// then
		assert.NotNil(t, step)
		assert.Equal(t, "testNS", step.namespace)
		assert.Equal(t, "testing/co:0.1", step.chart)
		assert.Equal(t, helmClientMock, step.helmClient)
	})
}

func TestComponentChartInstallerStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := &cesComponentChartInstallerStep{
			chart:     "testChart",
			namespace: "testNS",
		}

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Install component-chart from testChart in namespace testNS", desc)
	})
}

func TestComponentChartInstallerStep_PerformSetupStep(t *testing.T) {
	t.Run("should successfully perform setup", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		chartSpec := &helmclient.ChartSpec{
			ReleaseName:     "testChart",
			ChartName:       "foo/testChart",
			Namespace:       "testing",
			Version:         "0.1",
			Timeout:         time.Second * 300,
			Atomic:          true,
			CreateNamespace: true,
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgrade(testCtx, chartSpec).Return(nil)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart:0.1",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
	t.Run("should successfully perform setup for 'latest' version", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		chartSpec := &helmclient.ChartSpec{
			ReleaseName:     "testChart",
			ChartName:       "foo/testChart",
			Namespace:       "testing",
			Version:         "",
			Timeout:         time.Second * 300,
			Atomic:          true,
			CreateNamespace: true,
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgrade(testCtx, chartSpec).Return(nil)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart:latest",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to perform setup for error in helmClient when installing the crd chart", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		chartSpec := &helmclient.ChartSpec{
			ReleaseName:     "testChart",
			ChartName:       "foo/testChart",
			Namespace:       "testing",
			Version:         "0.1",
			Timeout:         time.Second * 300,
			Atomic:          true,
			CreateNamespace: true,
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgrade(testCtx, chartSpec).Return(assert.AnError)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart:0.1",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should fail to perform setup for wrong chart format", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		helmClientMock := newMockHelmClient(t)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart_0.1",
			helmClient: helmClientMock,
		}
		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "componentChart 'foo/testChart_0.1' has a wrong format. Must be '<chartName>:<version>'; e.g.: 'foo/bar:1.2.3'")
	})
	t.Run("should fail to perform setup for empty chartName", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		helmClientMock := newMockHelmClient(t)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/:latest",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "error reading chartname 'foo/': wrong format")
	})
	t.Run("should fail to perform setup for error in helmClient when installing the chart", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		chartSpec := &helmclient.ChartSpec{
			ReleaseName:     "testChart",
			ChartName:       "foo/testChart",
			Namespace:       "testing",
			Version:         "0.1",
			Timeout:         time.Second * 300,
			Atomic:          true,
			CreateNamespace: true,
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgrade(testCtx, chartSpec).Return(assert.AnError)

		step := &cesComponentChartInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart:0.1",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})
}
