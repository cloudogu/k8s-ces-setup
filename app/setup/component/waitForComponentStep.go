package component

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/setup/dogus"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/retry-lib/retry"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	retrywatch "k8s.io/client-go/tools/watch"
	"os"
	"time"
)

const (
	v1LabelK8sComponent = "app.kubernetes.io/name"
)

// DefaultComponentWaitTimeOut30Minutes is the default timeout.
// Since the components are actual applied unordered we use a 30 Minutes timeout as default.
var DefaultComponentWaitTimeOut30Minutes = time.Second * 1800

// componentTimeOutInSecondsEnvVar contains the name of the environment variable that may replace the default component wait timeout.
// An environment variable with this name must contain the seconds as reasonably sized integer (=< int64)
const componentTimeOutInSecondsEnvVar = "COMPONENT_TIMEOUT_SECS"

type waitForComponentStep struct {
	client        componentsClient
	labelSelector string
	namespace     string
	componentName string
	timeout       time.Duration
}

// NewWaitForComponentStep creates a new setup step which on waits for a component with a specific label
func NewWaitForComponentStep(client componentsClient, componentName string, namespace string, timeout time.Duration) *waitForComponentStep {
	return &waitForComponentStep{
		client:        client,
		namespace:     namespace,
		componentName: componentName,
		labelSelector: CreateComponentLabelSelector(componentName),
		timeout:       timeout,
	}
}

// GetStepDescription return the human-readable description of the step
func (wfcs *waitForComponentStep) GetStepDescription() string {
	return fmt.Sprintf("Wait for component with selector %s to be ready", wfcs.labelSelector)
}

// PerformSetupStep implements all actions in this step
func (wfcs *waitForComponentStep) PerformSetupStep(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, wfcs.timeout)
	defer cancel()
	return wfcs.isComponentReady(timeoutCtx)
}

// isComponentStatusReady does a watch on a component and returns nil if the component is installed
func (wfcs *waitForComponentStep) isComponentReady(ctx context.Context) error {
	var get *v1.Component
	err := retry.OnErrorWithLimit(wfcs.timeout, errors.IsNotFound, func() error {
		var getErr error
		get, getErr = wfcs.client.Get(ctx, wfcs.componentName, metav1.GetOptions{})
		if getErr != nil && !errors.IsNotFound(getErr) {
			return fmt.Errorf("failed to get initial component cr %q: %w", wfcs.componentName, getErr)
		}

		return getErr
	})

	if err != nil {
		return err
	}

	if isComponentStatusReady(get) {
		return nil
	}

	watcher := componentReadyWatcher{client: wfcs.client, componentName: wfcs.componentName, labelSelector: wfcs.labelSelector}
	_, err = retrywatch.Until(ctx, get.ResourceVersion, watcher, watcher.checkComponentStatus)
	if err != nil {
		return fmt.Errorf("failed to wait for component with label %q with retry watch: %w", wfcs.labelSelector, err)
	}

	return nil
}

type componentReadyWatcher struct {
	client        componentsClient
	componentName string
	labelSelector string
}

// Watch creates a watch for the component defined in this step.
// This function will be called initially and on every retry if the watch gets canceled from a recoverable error.
func (crw componentReadyWatcher) Watch(options metav1.ListOptions) (watch.Interface, error) {
	logrus.Debugf("creating initial or retry watch for component %q", crw.componentName)
	options.LabelSelector = crw.labelSelector
	w, err := crw.client.Watch(context.Background(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to create watch for label %q: %w", crw.labelSelector, err)
	}

	return w, nil
}

// checkComponentStatus is a condition function that will be called on every watch event received from the retry watcher.
// If it returns true, nil the watch will end.
// If it returns false, nil the watch will continue and check further events.
// If it returns and error the watch will end and don't retry.
func (crw componentReadyWatcher) checkComponentStatus(event watch.Event) (bool, error) {
	logrus.Debugf("received %q watch event for checking component ready status", event.Type)
	switch event.Type {
	case watch.Error:
		status, ok := event.Object.(*metav1.Status)
		if !ok {
			return false, fmt.Errorf("failed to cast event object to status")
		} else {
			return false, fmt.Errorf("watch error message: %q, reason: %q", status.Message, status.Reason)
		}
	case watch.Added, watch.Modified:
		component, ok := event.Object.(*v1.Component)
		if !ok {
			logrus.Errorf("failed to cast event object to component: selector=[%s] type=[%s]; object=[%+v]", crw.labelSelector, event.Type, event.Object)
			return false, nil
		}
		if isComponentStatusReady(component) {
			return true, nil
		}
		return false, nil
	case watch.Deleted:
		return false, fmt.Errorf("abort watch because of component deletion")
	default:
		return false, nil
	}
}

func isComponentStatusReady(component *v1.Component) bool {
	if component.Status.Status == v1.ComponentStatusInstalled && component.Status.Health == v1.AvailableHealthStatus {
		logrus.Infof("component %q is installed and available", component.Spec.Name)
		return true
	}
	logrus.Debugf("component %q is not installed and not available", component.Spec.Name)
	return false
}

func CreateComponentLabelSelector(name string) string {
	return fmt.Sprintf("%s=%s", v1LabelK8sComponent, name)
}

// TimeoutInSeconds returns either DefaultComponentWaitTimeOut30Minutes or a positive integer if set as EnvVar
// COMPONENT_TIMEOUT_SECS. See also componentTimeOutInSecondsEnvVar
func TimeoutInSeconds() time.Duration {
	defaultTimeout := DefaultComponentWaitTimeOut30Minutes
	if podTimeoutRaw, ok := os.LookupEnv(componentTimeOutInSecondsEnvVar); ok {
		logrus.Infof("Custom component timeout found")
		return dogus.ParseTimeoutString(podTimeoutRaw, defaultTimeout)
	}

	return defaultTimeout
}
