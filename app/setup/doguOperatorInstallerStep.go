package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	v1rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func newDoguOperatorInstallerStep(clientSet kubernetes.Interface, resourceURL, version string) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{
		ClientSet:   clientSet,
		resourceURL: resourceURL,
		Version:     version,
		httpClient:  &http.Client{},
	}
}

type doguOperatorInstallerStep struct {
	ClientSet   kubernetes.Interface
	Version     string
	resourceURL string
	httpClient  *http.Client
}

func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator version %s", dois.Version)
}

func (dois *doguOperatorInstallerStep) PerformSetupStep() error {
	var yamlResourceFiles []string

	fileContent, err := fetchYaml(dois.resourceURL, dois.httpClient)
	yamlResourceFiles = splitYamlFileSections(fileContent)
	_, err = parseYamlResources(yamlResourceFiles)
	if err != nil {
		return err
	}

	return nil
}

func fetchYaml(url string, httpClient *http.Client) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not find YAML file '%s'", url)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bytes, nil
}

func splitYamlFileSections(resourceBytes []byte) []string {
	fileAsString := string(resourceBytes[:])

	const yamlFileSeparator = "---\n"
	preResult := strings.Split(fileAsString, yamlFileSeparator)

	cleanedResult := make([]string, 0)
	for _, section := range preResult {
		if section != "" {
			cleanedResult = append(cleanedResult, section)
		}
	}

	return cleanedResult
}

func parseYamlResources(yamlResourceFiles []string) (interface{}, error) {
	for _, f := range yamlResourceFiles {
		parseYamlResource(f)
	}

	return nil, nil
}

func parseYamlResource(resource string) {
	sch := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, _ /*groupVersionKind*/, err := decode([]byte(resource), nil, nil)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
	}

	switch o := obj.(type) {
	case *v1core.Namespace:
		fmt.Printf("Namespace: %+v\n", o)
	case *v1core.Pod:
		fmt.Printf("Pod: %+v\n", o)
	case *v1rbac.Role:
		fmt.Printf("Role: %+v\n", o)
	case *v1rbac.RoleBinding:
		fmt.Printf("RoleBinding: %+v\n", o)
	case *v1rbac.ClusterRole:
		fmt.Printf("ClusterRole: %+v\n", o)
	case *v1rbac.ClusterRoleBinding:
		fmt.Printf("ClusterRoleBinding: %+v\n", o)
	case *v1core.ServiceAccount:
		fmt.Printf("ServiceAccount: %+v\n", o)
	case *v1apps.Deployment:
		fmt.Printf("Deployment: %+v\n", o)
	default:
		//o is unknown for us
	}
}
