package component

import (
	"context"
	"testing"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentOperatorInstallerStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		setupCtx := &ctx.SetupContext{
			AppConfig: &ctx.Config{
				TargetNamespace:        "testNS",
				ComponentOperatorChart: "testing/co:0.1",
			},
		}
		helmClientMock := newMockHelmClient(t)

		// when
		step := NewComponentOperatorInstallerStep(setupCtx, helmClientMock)

		// then
		assert.NotNil(t, step)
		assert.Equal(t, "testNS", step.namespace)
		assert.Equal(t, "testing/co:0.1", step.chart)
		assert.Equal(t, helmClientMock, step.helmClient)
	})
}

func TestComponentOperatorInstallerStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := &componentOperatorInstallerStep{
			chart: "testChart",
		}

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Install component-operator from testChart", desc)
	})
}

func TestComponentOperatorInstallerStep_PerformSetupStep(t *testing.T) {
	t.Run("should successfully perform setup", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgradeChart(testCtx, "testing", "foo/testChart", "0.1").Return(nil)

		step := &componentOperatorInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart:0.1",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
	t.Run("should fail to perform setup for wrong chart format", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		helmClientMock := newMockHelmClient(t)

		step := &componentOperatorInstallerStep{
			namespace:  "testing",
			chart:      "foo/testChart_0.1",
			helmClient: helmClientMock,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "componentOperatorChart 'foo/testChart_0.1' has a wrong format. Must be '<chartName>:<version>'; e.g.: 'foo/bar:1.2.3'")
	})
	t.Run("should fail to perform setup for error in helmClient", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().InstallOrUpgradeChart(testCtx, "testing", "foo/testChart", "0.1").Return(assert.AnError)

		step := &componentOperatorInstallerStep{
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
