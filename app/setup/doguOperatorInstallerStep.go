package setup

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core"
	"k8s.io/client-go/rest"
	"regexp"
)

// namespaces follow RFC 1123 DNS-label rules, see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
var namespaceRfc1123Regex, _ = regexp.Compile(`(\s+namespace:\s+)"?([a-z0-9][a-z0-9-]{0,61}[a-z0-9])"?`)

func newDoguOperatorInstallerStep(clusterConfig *rest.Config, setupCtx context.SetupContext) *doguOperatorInstallerStep {
	return &doguOperatorInstallerStep{
		namespace:   setupCtx.AppConfig.Namespace,
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
	fileContent = replaceNamespace(fileContent, dois.namespace)

	err = dois.k8sClient.Apply(fileContent, dois.namespace)
	if err != nil {
		return err
	}

	return nil
}

func replaceNamespace(content []byte, namespace string) []byte {
	// do not re-use possible quotation marks because DNS labels are also proper YAML values
	return namespaceRfc1123Regex.ReplaceAll(content, []byte("${1}"+namespace))
}
