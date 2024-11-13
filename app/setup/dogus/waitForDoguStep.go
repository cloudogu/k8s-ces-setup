package dogus

import (
	"context"
	"fmt"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	retrywatch "k8s.io/client-go/tools/watch"
	"k8s.io/client-go/util/retry"
	"os"
	"strconv"
	"time"
)

const (
	v1LabelDogu = "dogu.name"
)

var DefaultDoguWaitTimeOut5Minutes = time.Second * 300

// podTimeOutInSecondsEnvVar contains the name of the environment variable that may replace the default pod wait timeout.
// An environment variable with this name must contain the seconds as reasonably sized integer (=< int64)
const podTimeOutInSecondsEnvVar = "POD_TIMEOUT_SECS"

type waitForDoguStep struct {
	client        doguClient
	labelSelector string
	namespace     string
	doguName      string
	timeout       time.Duration
}

// NewWaitForDoguStep creates a new setup step which on waits for a dogu with a specific label
func NewWaitForDoguStep(client doguClient, doguName string, namespace string, timeout time.Duration) *waitForDoguStep {
	return &waitForDoguStep{
		client:        client,
		namespace:     namespace,
		doguName:      doguName,
		labelSelector: CreateDoguLabelSelector(doguName),
		timeout:       timeout,
	}
}

// GetStepDescription return the human-readable description of the step
func (wfds *waitForDoguStep) GetStepDescription() string {
	return fmt.Sprintf("Wait for dogu with selector %s to be ready", wfds.labelSelector)
}

// PerformSetupStep implements all actions in this step
func (wfds *waitForDoguStep) PerformSetupStep(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, wfds.timeout)
	defer cancel()
	return wfds.isDoguReady(timeoutCtx)
}

// isDoguReady does a watch on a dogu and returns nil if the dogu is installed
func (wfds *waitForDoguStep) isDoguReady(ctx context.Context) error {
	var get *v2.Dogu
	err := retry.OnError(retry.DefaultBackoff, errors.IsNotFound, func() error {
		var getErr error
		get, getErr = wfds.client.Get(ctx, wfds.doguName, metav1.GetOptions{})
		if getErr != nil && !errors.IsNotFound(getErr) {
			return fmt.Errorf("failed to get initial dogu cr %q: %w", wfds.doguName, getErr)
		}

		return getErr
	})

	if err != nil {
		return err
	}

	if isDoguStatusReady(get) {
		return nil
	}

	watcher := doguReadyWatcher{client: wfds.client, doguName: wfds.doguName, labelSelector: wfds.labelSelector}
	_, err = retrywatch.Until(ctx, get.ResourceVersion, watcher, watcher.checkDoguStatus)
	if err != nil {
		return fmt.Errorf("failed to wait for dogu with label %q with retry watch: %w", wfds.labelSelector, err)
	}

	return nil
}

type doguReadyWatcher struct {
	client        doguClient
	doguName      string
	labelSelector string
}

// Watch creates a watch for the dogu defined in this step.
// This function will be called initially and on every retry if the watch gets canceled from a recoverable error.
func (drw doguReadyWatcher) Watch(options metav1.ListOptions) (watch.Interface, error) {
	logrus.Debugf("creating initial or retry watch for dogu %q", drw.doguName)
	options.LabelSelector = drw.labelSelector
	w, err := drw.client.Watch(context.Background(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to create watch for label %q: %w", drw.labelSelector, err)
	}

	return w, nil
}

// checkDoguStatus is a condition function that will be called on every watch event received from the retry watcher.
// If it returns true, nil the watch will end.
// If it returns false, nil the watch will continue and check further events.
// If it returns and error the watch will end and don't retry.
func (crw doguReadyWatcher) checkDoguStatus(event watch.Event) (bool, error) {
	logrus.Debugf("received %q watch event for checking dogu ready status", event.Type)
	switch event.Type {
	case watch.Error:
		status, ok := event.Object.(*metav1.Status)
		if !ok {
			return false, fmt.Errorf("failed to cast event object to status")
		} else {
			return false, fmt.Errorf("watch error message: %q, reason: %q", status.Message, status.Reason)
		}
	case watch.Added, watch.Modified:
		dogu, ok := event.Object.(*v2.Dogu)
		if !ok {
			logrus.Errorf("failed to cast event object to dogu: selector=[%s] type=[%s]; object=[%+v]", crw.labelSelector, event.Type, event.Object)
			return false, nil
		}
		if isDoguStatusReady(dogu) {
			return true, nil
		}
		return false, nil
	case watch.Deleted:
		return false, fmt.Errorf("abort watch because of dogu deletion")
	default:
		return false, nil
	}
}

func isDoguStatusReady(dogu *v2.Dogu) bool {
	if dogu.Status.Status == v2.DoguStatusInstalled && dogu.Status.Health == v2.AvailableHealthStatus {
		logrus.Infof("dogu %q is installed and available", dogu.Spec.Name)
		return true
	}
	logrus.Debugf("dogu %q is not installed and not available", dogu.Spec.Name)
	return false
}

func CreateDoguLabelSelector(name string) string {
	return fmt.Sprintf("%s=%s", v1LabelDogu, name)
}

// DoguTimeoutInSeconds returns either DefaultDoguWaitTimeOut5Minutes or a positive integer if set as EnvVar
// POD_TIMEOUT_SECS. See also podTimeOutInSecondsEnvVar
func DoguTimeoutInSeconds() time.Duration {
	if podTimeoutRaw, ok := os.LookupEnv(podTimeOutInSecondsEnvVar); ok {
		logrus.Infof("Custom pod timeout found")

		podTimeout, err := strconv.ParseInt(podTimeoutRaw, 10, 32)
		if err != nil {
			logrus.Errorf("Failed to parse seconds into pod timeout %s=%s (fallback to %0.f): %s",
				podTimeOutInSecondsEnvVar, podTimeoutRaw, DefaultDoguWaitTimeOut5Minutes.Seconds(), err.Error())
			return DefaultDoguWaitTimeOut5Minutes
		}

		if podTimeout < 0 {
			logrus.Errorf("Found negative pod timeout %s=%d (fallback to %0.f)",
				podTimeOutInSecondsEnvVar, podTimeout, DefaultDoguWaitTimeOut5Minutes.Seconds())
			return DefaultDoguWaitTimeOut5Minutes
		}

		return time.Duration(podTimeout) * time.Second
	}

	return DefaultDoguWaitTimeOut5Minutes
}
