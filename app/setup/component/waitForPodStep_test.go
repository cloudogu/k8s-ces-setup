package component

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"testing"
	"time"
)

func TestNewWaitForPodStep(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		// when
		step := NewWaitForPodStep(&fake.Clientset{}, "selector", "namespace", 5)

		// then
		require.NotNil(t, step)
	})
}

func Test_waitForPodStep_GetStepDescription(t *testing.T) {
	t.Parallel()

	t.Run("success get description", func(t *testing.T) {
		// given
		step := &waitForPodStep{labelSelector: "selector"}

		// when
		description := step.GetStepDescription()

		// then
		assert.Equal(t, "Wait for pod with selector selector to be ready", description)
	})
}

func Test_waitForPodStep_PerformSetupStep(t *testing.T) {
	t.Parallel()

	selector := "app=test"
	labels := make(map[string]string)
	labels["app"] = "test"
	t.Run("successfull perform step with ready pod", func(t *testing.T) {
		// given
		pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: labels}}
		clientset := fake.NewSimpleClientset(pod)

		watcher := watch.NewFake()
		clientset.PrependWatchReactor("pods", k8stesting.DefaultWatchReactor(watcher, nil))

		go func() {
			time.Sleep(time.Second)
			pod.Status = v1.PodStatus{Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}}}
			watcher.Add(pod)
		}()

		step := &waitForPodStep{
			clientSet:     clientset,
			labelSelector: selector,
			namespace:     "namespace",
			timeout:       time.Second * 5,
		}

		// when
		err := step.PerformSetupStep()

		// then
		require.NoError(t, err)
	})

	t.Run("failed because reaching timeout", func(t *testing.T) {
		// given
		pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: labels}}
		clientset := fake.NewSimpleClientset(pod)

		watcher := watch.NewFake()
		clientset.PrependWatchReactor("pods", k8stesting.DefaultWatchReactor(watcher, nil))

		go func() {
			time.Sleep(time.Second)
			pod.Status = v1.PodStatus{Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionFalse}}}
			watcher.Add(pod)
		}()
		timeout := time.Second * 3
		step := &waitForPodStep{
			clientSet:     clientset,
			labelSelector: selector,
			namespace:     "namespace",
			timeout:       timeout,
		}

		// when
		before := time.Now()
		err := step.PerformSetupStep()
		executionTime := time.Since(before)

		// then
		require.Error(t, err)
		assert.Greater(t, executionTime, time.Second*3)
		assert.Greater(t, time.Second*4, executionTime)
		assert.Contains(t, "pod is not ready: timeout reached", err.Error())
	})

	t.Run("failed to cast runtime object (no pod)", func(t *testing.T) {
		// given
		pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test", Labels: labels}}
		clientset := fake.NewSimpleClientset(pod)

		watcher := watch.NewFake()
		clientset.PrependWatchReactor("pods", k8stesting.DefaultWatchReactor(watcher, nil))

		go func() {
			time.Sleep(time.Second)
			service := &v1.Service{}
			watcher.Add(service)
		}()

		step := &waitForPodStep{
			clientSet:     clientset,
			labelSelector: selector,
			namespace:     "namespace",
			timeout:       time.Second * 5,
		}

		// when
		err := step.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, "failed to cast event object to pod", err.Error())
	})
}
