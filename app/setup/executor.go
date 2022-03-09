package setup

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// ExecutorStep describes a valid step in the setup
type ExecutorStep interface {
	// GetStepDescription returns the description of the setup step. The Executor prints the description of every step
	// when executing the setup
	GetStepDescription() string
	// PerformSetupStep is called for every registered step when executing the setup
	PerformSetupStep() error
}

// Executor is responsible to perform the actual steps of the setup
type Executor struct {
	ClientSet kubernetes.Interface `json:"clientSet"`
	Config    config.Config        `json:"config"`
	Steps     []ExecutorStep       `json:"steps"`
}

// NewExecutor creates a new setup executor with the given app configuration
func NewExecutor(clientSet kubernetes.Interface, appConfig config.Config) Executor {
	e := Executor{
		ClientSet: clientSet,
		Config:    appConfig,
	}
	return e
}

// RegisterSetupStep adds a new step to the setup
func (e *Executor) RegisterSetupStep(step ExecutorStep) {
	logrus.Debugf("Register setup step [%s]", step.GetStepDescription())
	e.Steps = append(e.Steps, step)
}

// PerformSetup starts the setup and executes all registered setup steps
func (e *Executor) PerformSetup() error {
	logrus.Print("Starting the setup process")

	for _, step := range e.Steps {
		logrus.Printf("Setup-Step: %s", step.GetStepDescription())

		err := step.PerformSetupStep()
		if err != nil {
			return fmt.Errorf("failed to perform step [%s]; %w", step.GetStepDescription(), err)
		}
	}

	return nil
}
