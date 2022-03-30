package setup

import (
	"bytes"
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
	"regexp"
)

// namespaces follow RFC 1123 DNS-label rules, see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
var namespacedResourcesRfc1123Regex, _ = regexp.Compile(`(\s+namespace:\s+)"?([a-z0-9][a-z0-9-]{0,61}[a-z0-9])"?`)

func newDoguOperatorInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{
		namespace:   setupCtx.AppConfig.TargetNamespace,
		resourceURL: setupCtx.AppConfig.DoguOperatorURL,
		fileClient:  core.NewFileClient(setupCtx.AppVersion),
		k8sClient:   core.NewK8sClient(clusterConfig),
	}
}

type doguOperatorInstallerStep struct {
	namespace   string
	resourceURL string
	fileClient  fileClient
	k8sClient   k8sClient
}

// GetStepDescription returns a human-readable description of the dogu operator installation step.
func (dois *doguOperatorInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install dogu operator from %s", dois.resourceURL)
}

// PerformSetupStep installs the dogu operator.
func (dois *doguOperatorInstallerStep) PerformSetupStep() error {
	fileContent, err := dois.fileClient.Get(dois.resourceURL)
	if err != nil {
		return err
	}

	// avoid extra namespace validation: during the namespace creation earlier it was already validated by the K8s API
	fileContent = replaceNamespacedResources(fileContent, dois.namespace)
	fileContent = removeLegacyNamespaceFromResources(fileContent)

	sections := splitYamlFileSections(fileContent)

	for _, section := range sections {
		err = dois.k8sClient.Apply(section, dois.namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitYamlFileSections(resourceBytes []byte) [][]byte {
	yamlFileSeparator := []byte("---\n")

	preResult := bytes.Split(resourceBytes, yamlFileSeparator)

	cleanedResult := make([][]byte, 0)
	for _, section := range preResult {
		if len(section) > 0 {
			cleanedResult = append(cleanedResult, section)
		}
	}

	return cleanedResult
}

func replaceNamespacedResources(content []byte, namespace string) []byte {
	// do not re-use possible quotation marks because DNS labels are also proper YAML values
	return namespacedResourcesRfc1123Regex.ReplaceAll(content, []byte("${1}"+namespace))
}

func removeLegacyNamespaceFromResources(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte(`apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: ecosystem
---`), []byte("---"))
}
