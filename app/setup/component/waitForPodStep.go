package component

import (
	"context"
	"fmt"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DefaultPodWaitTimeOut5Minutes contains the period  when a
var DefaultPodWaitTimeOut5Minutes = time.Second * 300

// podTimeOutInSecondsEnvVar contains the name of the environment variable that may replace the default pod wait timeout.
// An environment variable with this name must contain the seconds as reasonably sized integer (=< int64)
const podTimeOutInSecondsEnvVar = "POD_TIMEOUT_SECS"

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
	return wfps.isPodReady()
}

// isPodReady does a watch on a pod and returns nil if the pod is ready and the configured timout is not reached
func (wfps *waitForPodStep) isPodReady() error {
	backgroundCtx := context.Background()
	watch, err := wfps.clientSet.CoreV1().Pods(wfps.namespace).Watch(backgroundCtx, v1.ListOptions{LabelSelector: wfps.labelSelector})
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
		logger := log.FromContext(backgroundCtx)
		if !ok {
			logger.Error(fmt.Errorf("failed to cast event to pod: selector=[%s] type=[%s]; object=[%+v]",
				wfps.labelSelector, event.Type, event.Object), "error wait for pod")
			continue
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

// PodTimeoutInSeconds returns either DefaultPodWaitTimeOut5Minutes or a positive integer if set as EnvVar
// POD_TIMEOUT_SECS. See also podTimeOutInSecondsEnvVar
func PodTimeoutInSeconds() time.Duration {
	if podTimeoutRaw, ok := os.LookupEnv(podTimeOutInSecondsEnvVar); ok {
		logrus.Infof("Custom pod timeout found")

		podTimeout, err := strconv.ParseInt(podTimeoutRaw, 10, 32)
		if err != nil {
			logrus.Errorf("Failed to parse seconds into pod timeout %s=%s (fallback to %0.f): %s",
				podTimeOutInSecondsEnvVar, podTimeoutRaw, DefaultPodWaitTimeOut5Minutes.Seconds(), err.Error())
			return DefaultPodWaitTimeOut5Minutes
		}

		if podTimeout < 0 {
			logrus.Errorf("Found negative pod timeout %s=%d (fallback to %0.f)",
				podTimeOutInSecondsEnvVar, podTimeout, DefaultPodWaitTimeOut5Minutes.Seconds())
			return DefaultPodWaitTimeOut5Minutes
		}

		return time.Duration(podTimeout) * time.Second
	}

	return DefaultPodWaitTimeOut5Minutes
}
