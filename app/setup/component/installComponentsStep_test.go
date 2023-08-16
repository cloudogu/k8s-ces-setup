package component

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewInstallComponentsStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		componentsClientMock := newMockComponentsClient(t)

		// when
		step := NewInstallComponentsStep(componentsClientMock, "comp", "testing", "0.0.2", "testNS")

		// then
		assert.NotNil(t, step)
		assert.Equal(t, componentsClientMock, step.client)
		assert.Equal(t, "comp", step.componentName)
		assert.Equal(t, "testing", step.componentNamespace)
		assert.Equal(t, "0.0.2", step.version)
		assert.Equal(t, "testNS", step.namespace)
	})
}

func TestInstallComponentsStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := &installComponentsStep{
			componentName:      "testComp",
			componentNamespace: "testing",
			version:            "0.2.3",
		}

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Installing component 'testing/testComp:0.2.3'", desc)
	})
}

func TestInstallComponentsStep_PerformSetupStep(t *testing.T) {
	t.Run("should successfully perform setup", func(t *testing.T) {
		// given
		namespace := "testNS"
		testCtx := context.TODO()

		expectedComponent := &v1.Component{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testComponent",
				Namespace: namespace,
				Labels: map[string]string{
					"app":                    "ces",
					"app.kubernetes.io/name": "testComponent",
				},
			},
			Spec: v1.ComponentSpec{
				Name:      "testComponent",
				Namespace: "testing",
				Version:   "4.5.6",
			},
		}

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Create(testCtx, expectedComponent, metav1.CreateOptions{}).Return(nil, nil)

		step := &installComponentsStep{
			client:             componentsClientMock,
			namespace:          namespace,
			componentName:      "testComponent",
			componentNamespace: "testing",
			version:            "4.5.6",
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
	t.Run("should fail to perform setup for error in component client", func(t *testing.T) {
		// given
		namespace := "testNS"
		testCtx := context.TODO()

		expectedComponent := &v1.Component{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testComponent",
				Namespace: namespace,
				Labels: map[string]string{
					"app":                    "ces",
					"app.kubernetes.io/name": "testComponent",
				},
			},
			Spec: v1.ComponentSpec{
				Name:      "testComponent",
				Namespace: "testing",
				Version:   "4.5.6",
			},
		}

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Create(testCtx, expectedComponent, metav1.CreateOptions{}).Return(nil, assert.AnError)

		step := &installComponentsStep{
			client:             componentsClientMock,
			namespace:          namespace,
			componentName:      "testComponent",
			componentNamespace: "testing",
			version:            "4.5.6",
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to apply component 'testing/testComponent:4.5.6' :")
	})
}