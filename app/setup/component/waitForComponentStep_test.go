package component

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWaitForComponentStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		componentsClientMock := newMockComponentsClient(t)

		// when
		step := NewWaitForComponentStep(componentsClientMock, "app=test", "testNS", 5*time.Second)

		// then
		assert.NotNil(t, step)
		assert.Equal(t, componentsClientMock, step.client)
		assert.Equal(t, "app=test", step.labelSelector)
		assert.Equal(t, "testNS", step.namespace)
		assert.Equal(t, 5*time.Second, step.timeout)
	})
}

func TestWaitForComponentStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := &waitForComponentStep{
			labelSelector: "app=test",
		}

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Wait for component with selector app=test to be installed", desc)
	})
}

func TestWaitForComponentStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	namespace := "testNS"
	selector := "app=test"

	installedComponent := &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}}

	t.Run("should successfully perform setup", func(t *testing.T) {
		// given
		testCtx = context.TODO()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: selector}).Return(watcher, nil)

		step := &waitForComponentStep{
			client:        componentsClientMock,
			labelSelector: selector,
			namespace:     namespace,
			timeout:       5 * time.Second,
		}

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Add(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to perform setup on watch error", func(t *testing.T) {
		// given
		testCtx = context.TODO()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: selector}).Return(nil, assert.AnError)

		step := &waitForComponentStep{
			client:        componentsClientMock,
			labelSelector: selector,
			namespace:     namespace,
			timeout:       5 * time.Second,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to create watch on component:")
	})

	t.Run("should fail when watch returns an unexpected object", func(t *testing.T) {
		// given
		testCtx := context.TODO()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: selector}).Return(watcher, nil)

		step := &waitForComponentStep{
			client:        componentsClientMock,
			labelSelector: selector,
			namespace:     namespace,
			timeout:       50 * time.Second,
		}

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Add(&corev1.Namespace{})
			watcher.Add(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "error wait for component: failed to cast event to component: selector=[app=test] type=[ADDED];")
	})

	t.Run("should fail to perform setup on timeout", func(t *testing.T) {
		// given
		testCtx = context.TODO()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: selector}).Return(watcher, nil)

		step := &waitForComponentStep{
			client:        componentsClientMock,
			labelSelector: selector,
			namespace:     namespace,
			timeout:       5 * time.Second,
		}

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "component is not ready: timeout reached")
	})
}
