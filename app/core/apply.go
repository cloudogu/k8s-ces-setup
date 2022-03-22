package core

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/apply"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"time"
)

type k8sApplyClient struct {
	clusterConfig  *rest.Config
	kubectlFactory cmdutil.Factory
	applyIOStreams genericclioptions.IOStreams
}

// NewK8sClient creates a `kubectl`-like client which operates on the K8s API.
func NewK8sClient(clusterConfig *rest.Config) *k8sApplyClient {

	kubectlDefaultConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	flagsFile := pflag.NewFlagSet("f", pflag.ExitOnError)
	err := flagsFile.Parse([]string{"-f", "-"})
	if err != nil {
		logrus.Error("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		logrus.Error(err.Error())
		time.Sleep(60 * time.Second)
		panic("oh no, not again")
	}
	kubectlDefaultConfigFlags.AddFlags(flagsFile)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubectlDefaultConfigFlags)
	kubectlFactory := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	applyIO := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	return &k8sApplyClient{
		clusterConfig:  clusterConfig,
		kubectlFactory: kubectlFactory,
		applyIOStreams: applyIO,
	}
}

// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster.
func (dkc *k8sApplyClient) Apply(yamlResources []byte) error {

	cmdApply := apply.NewCmdApply("kubectl", dkc.kubectlFactory, dkc.applyIOStreams)

	return cmdApply.Execute()
}
