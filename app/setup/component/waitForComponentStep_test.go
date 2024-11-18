package component

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"testing"
)

func TestNewWaitForComponentStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		componentsClientMock := newMockComponentsClient(t)

		// when
		step := NewWaitForComponentStep(componentsClientMock, "k8s-ces-control", "ecosystem")

		// then
		assert.NotNil(t, step)
		assert.Equal(t, componentsClientMock, step.client)
		assert.Equal(t, "app.kubernetes.io/name=k8s-ces-control", step.labelSelector)
		assert.Equal(t, "k8s-ces-control", step.componentName)
		assert.Equal(t, "ecosystem", step.namespace)
	})
}

func TestWaitForComponentStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := NewWaitForComponentStep(newMockComponentsClient(t), "k8s-ces-control", "ecosystem")

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Wait for component with selector app.kubernetes.io/name=k8s-ces-control to be ready", desc)
	})
}

func TestWaitForComponentStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	testNamespace := "ecosystem"
	testComponentName := "k8s-ces-control"
	testSelector := "app.kubernetes.io/name=k8s-ces-control"

	testStartResourceVersion := "2771"
	testEndResourceVersion := "2772"
	testComponent := &v1.Component{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app.kubernetes.io/name": testComponentName}, Name: testComponentName, Namespace: testNamespace, ResourceVersion: testStartResourceVersion}}
	installedComponent := &v1.Component{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app.kubernetes.io/name": testComponentName}, Name: testComponentName, Namespace: testNamespace, ResourceVersion: testEndResourceVersion}, Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled, Health: v1.AvailableHealthStatus}}

	t.Run("should not start watch if the component is already ready", func(t *testing.T) {
		// given
		testCtx = context.Background()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(installedComponent, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should successfully end watch on modified event with ready component", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Modify(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should successfully end watch on multiple event sequence with ready component", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Add(testComponent)
			watcher.Add(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "wrong object", Namespace: testNamespace, ResourceVersion: "9999"}})
			watcher.Modify(testComponent)
			watcher.Modify(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should retry watch if retry channel is closed", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()
		watcher2 := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(nil, assert.AnError)
		componentsClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(watcher, nil)
		componentsClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(watcher2, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Add(testComponent)
			watcher.Stop()
			watcher2.Modify(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail getting the initial resource version", func(t *testing.T) {
		// given
		testCtx = context.Background()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(nil, assert.AnError)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to get initial component cr \"k8s-ces-control\"", testComponentName)
	})

	t.Run("should retry getting the initial resource version on not found error", func(t *testing.T) {
		// given
		testCtx = context.Background()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, "")).Times(1)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(nil, assert.AnError).Times(1)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to get initial component cr \"k8s-ces-control\"", testComponentName)
	})

	t.Run("should stop the watch on delete event", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Delete(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to wait for component with label \"app.kubernetes.io/name=k8s-ces-control\" with retry watch: abort watch because of component deletion")
	})

	t.Run("should not retry if the resource version is too old", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Error(&metav1.Status{Code: http.StatusGone, Message: "msg", Reason: "reason"})
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to wait for component with label \"app.kubernetes.io/name=k8s-ces-control\" with retry watch: watch error message: \"msg\", reason: \"reason\"")
	})

	t.Run("should retry by default on error events", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()
		watcher2 := watch.NewFake()
		watcher3 := watch.NewFake()
		watcher4 := watch.NewFake()

		componentsClientMock := newMockComponentsClient(t)
		componentsClientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).Return(testComponent, nil)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil).Times(1)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher2, nil).Times(1)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher3, nil).Times(1)
		componentsClientMock.EXPECT().Watch(testCtx, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher4, nil).Times(1)

		step := NewWaitForComponentStep(componentsClientMock, testComponentName, testNamespace)

		go func() {
			watcher.Error(&metav1.Status{Code: http.StatusGatewayTimeout})
			watcher2.Error(&metav1.Status{Code: http.StatusInternalServerError})
			watcher3.Error(&metav1.Status{})
			watcher4.Add(installedComponent)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
