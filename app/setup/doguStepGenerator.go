package setup

import (
	"fmt"
	componentEcoSystem "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/strings/slices"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote"
	setupcontext "github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"
	"github.com/cloudogu/k8s-ces-setup/app/setup/dogus"
	"github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
)

const (
	serviceAccountKindDogu      = "dogu"
	serviceAccountKindComponent = "component"
)

const (
	v1LabelDogu = "dogu.name"
)

// doguStepGenerator is responsible to generate the steps to install a dogu, i.e., applying the dogu cr into the cluster
// and waiting for the dependencies before doing so.
type doguStepGenerator struct {
	Client          kubernetes.Interface
	EcoSystemClient ecoSystem.EcoSystemV2Interface
	Dogus           *[]*core.Dogu
	Registry        remote.Registry
	namespace       string
	components      []string
	componentClient componentEcoSystem.ComponentInterface
}

// NewDoguStepGenerator creates a new generator capable of generating dogu installation steps.
func NewDoguStepGenerator(client kubernetes.Interface, clusterConfig *rest.Config, dogus setupcontext.Dogus, registry remote.Registry, namespace string, components []string) (*doguStepGenerator, error) {
	ecoSystemClient, err := ecoSystem.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create K8s EcoSystem client: %w", err)
	}
	componentClient, err := componentEcoSystem.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create component client: %w", err)
	}

	var doguList []*core.Dogu
	for _, doguString := range dogus.Install {
		dogu, err := getDoguByString(registry, doguString)
		if err != nil {
			return nil, err
		}

		doguList = append(doguList, dogu)
	}

	return &doguStepGenerator{Client: client, EcoSystemClient: ecoSystemClient, Dogus: &doguList, Registry: registry, namespace: namespace, components: components, componentClient: componentClient.Components(namespace)}, nil
}

// GenerateSteps generates dogu installation steps for all configured dogus.
func (dsg *doguStepGenerator) GenerateSteps() ([]ExecutorStep, error) {
	steps := []ExecutorStep{}

	installedDogus, err := core.SortDogusByDependencyWithError(*dsg.Dogus)
	if err != nil {
		return nil, fmt.Errorf("sorting dogus by dependency failed: %w", err)
	}
	waitList := map[string]bool{}

	for _, dogu := range installedDogus {
		// create wait step if needing a service account from a certain dogu
		steps = dsg.appendDoguWaitStepsIfNeeded(dogu, installedDogus, steps, waitList)
		installStep := dogus.NewInstallDogusStep(dsg.EcoSystemClient, dogu, dsg.namespace)
		steps = append(steps, installStep)
	}

	return steps, nil
}

func (dsg *doguStepGenerator) appendDoguWaitStepsIfNeeded(dogu *core.Dogu, installedDogus []*core.Dogu, steps []ExecutorStep, waitList map[string]bool) []ExecutorStep {
	for _, serviceAccountDependency := range dogu.ServiceAccounts {
		switch serviceAccountDependency.Kind {
		case "":
			fallthrough
		case serviceAccountKindDogu:
			if !shouldDoguWaitForSADogu(dogu, serviceAccountDependency, installedDogus) {
				logrus.Infof("skipping wait step for optional dogu %s service account creation", serviceAccountDependency.Type)
				continue
			}
			steps = dsg.createWaitStepForDogu(serviceAccountDependency, waitList, steps)
		case serviceAccountKindComponent:
			if !shouldDoguWaitForSAComponent(dogu, serviceAccountDependency, dsg.components) {
				logrus.Infof("skipping wait step for optional component %s service account creation", serviceAccountDependency.Type)
				continue
			}
			steps = dsg.createWaitStepForK8sComponent(serviceAccountDependency, waitList, steps)
		default:
			logrus.Errorf("unknown service account kind %s from dogu %s. skipping wait step creation for service account creation", serviceAccountDependency.Kind, dogu.GetSimpleName())
			continue
		}
	}

	return steps
}

func shouldDoguWaitForSAComponent(dogu *core.Dogu, serviceAccount core.ServiceAccount, configureComponents []string) bool {
	return !(isOptionalServiceAccount(dogu, serviceAccount) && !slices.Contains(configureComponents, serviceAccount.Type))
}

func shouldDoguWaitForSADogu(dogu *core.Dogu, serviceAccount core.ServiceAccount, configuredDogus []*core.Dogu) bool {
	return !(isOptionalServiceAccount(dogu, serviceAccount) && !isDoguConfigured(configuredDogus, serviceAccount.Type))
}

func isDoguConfigured(dogus []*core.Dogu, simpleName string) bool {
	for _, dogu := range dogus {
		if dogu.GetSimpleName() == simpleName {
			return true
		}
	}

	return false
}

func isOptionalServiceAccount(dogu *core.Dogu, serviceAccount core.ServiceAccount) bool {
	for _, optionalDependency := range dogu.OptionalDependencies {
		if optionalDependency.Name == serviceAccount.Type {
			return true
		}
	}

	return false
}

func (dsg *doguStepGenerator) createWaitStepForDogu(serviceAccountDependency core.ServiceAccount, waitList map[string]bool, steps []ExecutorStep) []ExecutorStep {
	labelSelector := fmt.Sprintf("%s=%s", v1LabelDogu, serviceAccountDependency.Type)

	if waitList[labelSelector] {
		return steps
	}

	waitForDependencyStep := dogus.NewWaitForPodStep(dsg.Client, labelSelector, dsg.namespace, dogus.PodTimeoutInSeconds())
	steps = append(steps, waitForDependencyStep)
	waitList[labelSelector] = true

	return steps
}

func (dsg *doguStepGenerator) createWaitStepForK8sComponent(serviceAccountDependency core.ServiceAccount, waitList map[string]bool, steps []ExecutorStep) []ExecutorStep {
	labelSelector := component.CreateComponentLabelSelector(serviceAccountDependency.Type)
	if waitList[labelSelector] {
		return steps
	}

	waitForDependencyStep := component.NewWaitForComponentStep(dsg.componentClient, serviceAccountDependency.Type, dsg.namespace)
	steps = append(steps, waitForDependencyStep)
	waitList[labelSelector] = true

	return steps
}

func getDoguByString(registry remote.Registry, doguString string) (*core.Dogu, error) {
	namespaceName, version, found := strings.Cut(doguString, ":")
	if !found {
		// get latest version
		latest, err := registry.Get(namespaceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version of dogu [%s]: %w", namespaceName, err)
		}

		return latest, nil
	} else {
		// get specific version
		latest, err := registry.GetVersion(namespaceName, version)
		if err != nil {
			return nil, fmt.Errorf("failed to get version [%s] of dogu [%s]: %w", version, namespaceName, err)
		}

		return latest, nil
	}
}
