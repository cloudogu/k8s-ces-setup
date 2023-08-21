package component

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefaultComponentWaitTimeOut5Minutes contains the default period  of 5 minutes to wait for components to get installed
var DefaultComponentWaitTimeOut5Minutes = time.Second * 300

type waitForComponentStep struct {
	client        componentsClient
	labelSelector string
	namespace     string
	timeout       time.Duration
}

// NewWaitForComponentStep creates a new setup step which on waits for a component with a specific label
func NewWaitForComponentStep(client componentsClient, labelSelector string, namespace string, timeout time.Duration) *waitForComponentStep {
	return &waitForComponentStep{
		client:        client,
		labelSelector: labelSelector,
		namespace:     namespace,
		timeout:       timeout,
	}
}

// GetStepDescription return the human-readable description of the step
func (wfcs *waitForComponentStep) GetStepDescription() string {
	return fmt.Sprintf("Wait for component with selector %s to be installed", wfcs.labelSelector)
}

// PerformSetupStep implements all actions in this step
func (wfcs *waitForComponentStep) PerformSetupStep(ctx context.Context) error {
	return wfcs.isComponentInstalled(ctx)
}

// isComponentInstalled does a watch on a component and returns nil if the component is installed and the configured timout is not reached
func (wfcs *waitForComponentStep) isComponentInstalled(ctx context.Context) error {
	watch, err := wfcs.client.Watch(ctx, metav1.ListOptions{LabelSelector: wfcs.labelSelector})
	if err != nil {
		return fmt.Errorf("failed to create watch on component: %w", err)
	}

	timer := time.NewTimer(wfcs.timeout)
	go func() {
		<-timer.C
		watch.Stop()
	}()

	for event := range watch.ResultChan() {
		component, ok := event.Object.(*v1.Component)
		if !ok {
			timer.Stop()
			return fmt.Errorf("error wait for component: failed to cast event to component: selector=[%s] type=[%s]; object=[%+v]",
				wfcs.labelSelector, event.Type, event.Object)
		}

		if component.Status.Status == v1.ComponentStatusInstalled {
			timer.Stop()
			return nil
		}
	}

	return fmt.Errorf("component is not ready: timeout reached")
}
