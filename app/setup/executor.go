package setup

import (
	"fmt"

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
	// ClientSet is the actual k8s client responsible for the k8s API communication
	ClientSet kubernetes.Interface `json:"client_set"`
	// Steps contains all necessary steps for the setup
	Steps []ExecutorStep `json:"steps"`
}

// NewExecutor creates a new setup executor with the given app configuration
func NewExecutor(clientSet kubernetes.Interface) *Executor {
	return &Executor{ClientSet: clientSet}
}

// RegisterSetupStep adds a new step to the setup
func (e *Executor) RegisterSetupStep(step ExecutorStep) {
	logrus.Debugf("Register setup step [%s]", step.GetStepDescription())
	e.Steps = append(e.Steps, step)
}

// PerformSetup starts the setup and executes all registered setup steps
func (e *Executor) PerformSetup() (err error, errCausingAction string) {
	logrus.Print("Starting the setup process")

	for _, step := range e.Steps {
		logrus.Printf("Setup-Step: %s", step.GetStepDescription())

		err := step.PerformSetupStep()
		if err != nil {
			return fmt.Errorf("failed to perform step [%s]: %w", step.GetStepDescription(), err), step.GetStepDescription()
		}
	}

	return nil, ""
}
