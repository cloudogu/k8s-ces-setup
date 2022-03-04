package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ExecutorStep describes a valid step in the setup
type ExecutorStep interface {
	GetName() string
	PerformSetupStep() error
}

// Executor is responsible to perform the actual steps of the setup
type Executor struct {
	ClientSet *kubernetes.Clientset `json:"client_set"`
	Config    config.Config         `json:"config"`
	Steps     []ExecutorStep        `json:"steps"`
}

// NewExecutor creates a new setup executor with the given app configuration
func NewExecutor(appConfig config.Config) (*Executor, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot load in cluster configuration; %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes configuration; %w", err)
	}

	e := &Executor{ClientSet: clientSet, Config: appConfig}
	registerSetupSteps(e)
	return e, nil
}

func registerSetupSteps(executor *Executor) {
	// create namespace
	executor.registerSetupStep(NewNamespaceCreationStep(executor.ClientSet, executor.Config.Namespace))
	// do other things
}

func (e *Executor) registerSetupStep(step ExecutorStep) {
	logrus.Debugf("Register setup step [%s]", step.GetName())
	e.Steps = append(e.Steps, step)
}

func (e *Executor) performSetup() error {
	logrus.Print("Starting the setup process")
	for _, step := range e.Steps {
		logrus.Print("Setup-Step: %s", step.GetName())
		err := step.PerformSetupStep()
		if err != nil {
			return fmt.Errorf("failed to perform step [%s]; %w", step.GetName(), err)
		}
	}

	return nil
}
