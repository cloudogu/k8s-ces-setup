package dogus

import (
	"context"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"testing"
)

func TestNewWaitForDoguStep(t *testing.T) {
	t.Run("create without error", func(t *testing.T) {
		// given
		doguClientMock := newMockDoguClient(t)

		// when
		step := NewWaitForDoguStep(doguClientMock, "cas", "ecosystem", DefaultDoguWaitTimeOut5Minutes)

		// then
		assert.NotNil(t, step)
		assert.Equal(t, doguClientMock, step.client)
		assert.Equal(t, "dogu.name=cas", step.labelSelector)
		assert.Equal(t, "cas", step.doguName)
		assert.Equal(t, "ecosystem", step.namespace)
	})
}

func TestWaitForDoguStep_GetStepDescription(t *testing.T) {
	t.Run("should get description", func(t *testing.T) {
		// given
		step := NewWaitForDoguStep(newMockDoguClient(t), "cas", "ecosystem", DefaultDoguWaitTimeOut5Minutes)

		// when
		desc := step.GetStepDescription()

		// then
		assert.Equal(t, "Wait for dogu with selector dogu.name=cas to be ready", desc)
	})
}

func assertTimeoutCtx(t *testing.T, ctx context.Context) {
	_, ok := ctx.Deadline()
	assert.True(t, ok)
}

func TestWaitForDoguStep_PerformSetupStep(t *testing.T) {
	t.Parallel()
	var testCtx = context.Background()

	testNamespace := "ecosystem"
	testDoguName := "cas"
	testSelector := "dogu.name=cas"

	testStartResourceVersion := "2771"
	testEndResourceVersion := "2772"
	testDogu := &v2.Dogu{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"dogu.name": testDoguName}, Name: testDoguName, Namespace: testNamespace, ResourceVersion: testStartResourceVersion}}
	installedDogu := &v2.Dogu{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"dogu.name": testDoguName}, Name: testDoguName, Namespace: testNamespace, ResourceVersion: testEndResourceVersion}, Status: v2.DoguStatus{Status: v2.DoguStatusInstalled, Health: v2.AvailableHealthStatus}}

	t.Run("should not start watch if the dogu is already ready", func(t *testing.T) {
		// given
		testCtx = context.Background()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(installedDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should successfully end watch on modified event with ready dogu", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Modify(installedDogu)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should successfully end watch on multiple event sequence with ready dogu", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Add(testDogu)
			watcher.Add(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "wrong object", Namespace: testNamespace, ResourceVersion: "9999"}})
			watcher.Modify(testDogu)
			watcher.Modify(installedDogu)
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

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(nil, assert.AnError)
		doguClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(watcher, nil)
		doguClientMock.EXPECT().Watch(context.Background(), metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Times(1).Return(watcher2, nil)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Add(testDogu)
			watcher.Stop()
			watcher2.Modify(installedDogu)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail getting the initial resource version", func(t *testing.T) {
		// given
		testCtx = context.Background()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(nil, assert.AnError).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to get initial dogu cr \"cas\"", testDoguName)
	})

	t.Run("should retry getting the initial resource version on not found error", func(t *testing.T) {
		// given
		testCtx = context.Background()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, "")).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		}).Times(1)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(nil, assert.AnError).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		}).Times(1)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		require.ErrorContains(t, err, "failed to get initial dogu cr \"cas\"", testDoguName)
	})

	t.Run("should stop the watch on delete event", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Delete(installedDogu)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to wait for dogu with label \"dogu.name=cas\" with retry watch: abort watch because of dogu deletion")
	})

	t.Run("should not retry if the resource version is too old", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Error(&metav1.Status{Code: http.StatusGone, Message: "msg", Reason: "reason"})
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to wait for dogu with label \"dogu.name=cas\" with retry watch: watch error message: \"msg\", reason: \"reason\"")
	})

	t.Run("should retry by default on error events", func(t *testing.T) {
		// given
		testCtx = context.Background()

		watcher := watch.NewFake()
		watcher2 := watch.NewFake()
		watcher3 := watch.NewFake()
		watcher4 := watch.NewFake()

		doguClientMock := newMockDoguClient(t)
		doguClientMock.EXPECT().Get(mock.Anything, testDoguName, metav1.GetOptions{}).Return(testDogu, nil).Run(func(ctx context.Context, _ string, _ metav1.GetOptions) {
			assertTimeoutCtx(t, ctx)
		})
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher, nil).Times(1)
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher2, nil).Times(1)
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher3, nil).Times(1)
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{LabelSelector: testSelector, ResourceVersion: testStartResourceVersion, AllowWatchBookmarks: true}).Return(watcher4, nil).Times(1)

		step := NewWaitForDoguStep(doguClientMock, testDoguName, testNamespace, DefaultDoguWaitTimeOut5Minutes)

		go func() {
			watcher.Error(&metav1.Status{Code: http.StatusGatewayTimeout})
			watcher2.Error(&metav1.Status{Code: http.StatusInternalServerError})
			watcher3.Error(&metav1.Status{})
			watcher4.Add(installedDogu)
		}()

		// when
		err := step.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
