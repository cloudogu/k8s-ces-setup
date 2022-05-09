package data

import (
	gocontext "context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"strings"
)

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	schemeGroupVersion = schema.GroupVersion{Group: "k8s.cloudogu.com", Version: "v1"}
)

type installDogusStep struct {
	client    *rest.RESTClient
	dogus     context.Dogus
	Registry  remote.Registry
	namespace string
}

func NewInstallDogusStep(clusterConfig *rest.Config, dogus context.Dogus, registry remote.Registry, namespace string) (*installDogusStep, error) {
	client, err := getDoguRestClient(clusterConfig)
	if err != nil {
		return nil, err
	}

	return &installDogusStep{client: client, dogus: dogus, Registry: registry, namespace: namespace}, nil
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
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
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

// GetStepDescription return the human-readable description of the step
func (ids *installDogusStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogus")
}

func (ids *installDogusStep) PerformSetupStep() error {
	// generate dogu cr from dogus
	for _, dogu := range ids.dogus.Install {
		namespaceName, version, found := strings.Cut(dogu, ":")
		if !found {
			latest, err := ids.Registry.GetVersion(namespaceName, "")
			if err != nil {
				return fmt.Errorf("failed to get latest version of dogu %s: %w", dogu, err)
			}
			version = latest.Version
		}

		_, name, found := strings.Cut(namespaceName, "/")
		cr := getDoguCr(name, namespaceName, version, ids.namespace)
		result := ids.client.Post().Namespace(ids.namespace).Resource("dogus").Body(cr).Do(gocontext.Background())
		err := result.Error()
		if err != nil {
			return fmt.Errorf("failed to apply dogu %s: %w", name, err)
		}
	}

	return nil
}

func getDoguCr(name string, namespaceName string, version string, k8sNamespace string) *v1.Dogu {
	cr := &v1.Dogu{}
	labels := make(map[string]string)
	labels["app"] = "ces"
	labels["dogu"] = name
	cr.Name = name
	cr.Namespace = k8sNamespace
	cr.Spec.Name = namespaceName
	cr.Spec.Version = version
	cr.Labels = labels

	return cr
}
