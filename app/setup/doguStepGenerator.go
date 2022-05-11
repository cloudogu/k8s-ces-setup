package setup

import (
	"fmt"
	"strings"
	"time"

	"github.com/cloudogu/k8s-ces-setup/app/setup/dogus"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/cloudogu/k8s-ces-setup/app/setup/component"
	"k8s.io/client-go/kubernetes"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	schemeGroupVersion = schema.GroupVersion{Group: "k8s.cloudogu.com", Version: "v1"}
)

// doguStepGenerator is responsible to generate the steps to install a dogu, i.e., applying the dogu cr into the cluster
// and waiting for the dependencies before doing so.
type doguStepGenerator struct {
	Client     kubernetes.Interface
	RestClient *rest.RESTClient
	Dogus      *[]*core.Dogu
	Registry   remote.Registry
	namespace  string
}

// NewDoguStepGenerator creates a new generator capable of generating dogu installation steps.
func NewDoguStepGenerator(client kubernetes.Interface, clusterConfig *rest.Config, dogus context.Dogus, registry remote.Registry, namespace string) (*doguStepGenerator, error) {
	restClient, err := getDoguRestClient(clusterConfig)
	if err != nil {
		return nil, err
	}

	doguList := []*core.Dogu{}
	for _, doguString := range dogus.Install {
		dogu, err := getDoguByString(registry, doguString)
		if err != nil {
			return nil, err
		}

		doguList = append(doguList, dogu)
	}

	return &doguStepGenerator{Client: client, RestClient: restClient, Dogus: &doguList, Registry: registry, namespace: namespace}, nil
}

// GenerateSteps generates dogu installation steps for all configured dogus.
func (dsg *doguStepGenerator) GenerateSteps() []ExecutorStep {
	steps := []ExecutorStep{}

	installedDogus := core.SortDogusByDependency(*dsg.Dogus)
	waitForDoguList := ""
	for _, dogu := range installedDogus {
		// create wait step if needing a service account from a certain dogu
		for _, serviceAccountDepedency := range dogu.ServiceAccounts {
			labelSelector := fmt.Sprintf("dogu=%s", serviceAccountDepedency.Type)

			if !strings.Contains(waitForDoguList, fmt.Sprintf("[%s]", labelSelector)) {
				waitForDependencyStep := component.NewWaitForPodStep(dsg.Client, labelSelector, dsg.namespace, time.Second*300)
				steps = append(steps, waitForDependencyStep)
				waitForDoguList += fmt.Sprintf("[%s]", labelSelector)
			}
		}

		// create install step
		installStep := dogus.NewInstallDogusStep(dsg.RestClient, dogu, dsg.namespace)
		steps = append(steps, installStep)
	}

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

func getDoguRestClient(config *rest.Config) (*rest.RESTClient, error) {
	err := AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to add scheme: %w", err)
	}

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes RestClient: %w", err)
	}

	return client, nil
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(schemeGroupVersion,
		&v1.Dogu{},
		&v1.DoguList{},
	)

	metav1.AddToGroupVersion(scheme, schemeGroupVersion)
	return nil
}
