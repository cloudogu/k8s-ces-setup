package component

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type waitForPodStep struct {
	clientSet     kubernetes.Interface
	labelSelector string
	namespace     string
	timeout       time.Duration
}

// NewWaitForPodStep creates a new setup step which on waits for a pod with a specific label
func NewWaitForPodStep(client kubernetes.Interface, labelSelector string, namespace string, timeout time.Duration) *waitForPodStep {
	return &waitForPodStep{
		clientSet:     client,
		labelSelector: labelSelector,
		namespace:     namespace,
		timeout:       timeout,
	}
}

// GetStepDescription return the human-readable description of the step
func (wfps *waitForPodStep) GetStepDescription() string {
	return fmt.Sprintf("Wait for pod with selector %s to be ready", wfps.labelSelector)
}

// PerformSetupStep implements all actions in this step
func (wfps *waitForPodStep) PerformSetupStep() error {
	err := wfps.isPodReady()

	return err
}

// isPodReady does a watch on a pod and returns nil if the pod is ready and the configured timout is not reached
func (wfps *waitForPodStep) isPodReady() error {
	watch, err := wfps.clientSet.CoreV1().Pods(wfps.namespace).Watch(context.Background(), v1.ListOptions{LabelSelector: wfps.labelSelector})
	if err != nil {
		return fmt.Errorf("failed to create watch on pod: %w", err)
	}

	timer := time.NewTimer(wfps.timeout)
	go func() {
		<-timer.C
		watch.Stop()
	}()

	for event := range watch.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			return fmt.Errorf("failed to cast event object to pod")
		}

		for _, cd := range pod.Status.Conditions {
			if cd.Type == corev1.PodReady && cd.Status == corev1.ConditionTrue {
				timer.Stop()
				return nil
			}
		}
	}

	return fmt.Errorf("pod is not ready: timeout reached")
}
