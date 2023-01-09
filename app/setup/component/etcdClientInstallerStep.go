package component

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
)

func NewEtcdClientInstallerStep(clientSet kubernetes.Interface, setupCtx *ctx.SetupContext) *etcdClientInstallerStep {
	etcdServiceUrl := fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", setupCtx.AppConfig.TargetNamespace)

	return &etcdClientInstallerStep{
		clientSet:       clientSet,
		targetNamespace: setupCtx.AppConfig.TargetNamespace,
		imageURL:        setupCtx.AppConfig.EtcdClientImageRepo,
		etcdServiceUrl:  etcdServiceUrl,
	}
}

type etcdClientInstallerStep struct {
	clientSet       kubernetes.Interface
	targetNamespace string
	imageURL        string
	etcdServiceUrl  string
}

// GetStepDescription returns a human-readable description of the etcd client step.
func (ecis *etcdClientInstallerStep) GetStepDescription() string {
	return fmt.Sprintf("Install etcd client from %s", ecis.imageURL)
}

// PerformSetupStep installs an etcd client.
func (ecis *etcdClientInstallerStep) PerformSetupStep() error {
	err := ecis.installClient()
	if err != nil {
		return err
	}

	return nil
}

func (ecis *etcdClientInstallerStep) installClient() error {
	err := ecis.createPod(ecis.etcdServiceUrl)
	if err != nil {
		return err
	}

	return nil
}

func (ecis *etcdClientInstallerStep) createPod(etcdServiceUrl string) error {
	etcdClientName := "etcd-client"
	const etcdAPIVersion = "2"
	etcdClientLabels := make(map[string]string)
	etcdClientLabels["app"] = "ces"
	etcdClientLabels["app.kubernetes.io/name"] = "etcdClient"
	mountServiceAccountToken := true

	etcdPod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:   etcdClientName,
			Labels: etcdClientLabels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    etcdClientName,
					Image:   ecis.imageURL,
					Command: []string{"sleep", "infinity"},
					Env: []corev1.EnvVar{
						{
							Name:  "ETCDCTL_API",
							Value: etcdAPIVersion,
						},
						{
							Name:  "ETCDCTL_ENDPOINTS",
							Value: etcdServiceUrl,
						},
					},
				},
			},
			AutomountServiceAccountToken: &mountServiceAccountToken,
		},
		Status: corev1.PodStatus{},
	}

	_, err := ecis.clientSet.CoreV1().Pods(ecis.targetNamespace).Create(context.Background(), etcdPod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("cannot create etcd pod in namespace %s with clientset: %w", ecis.targetNamespace, err)
	}

	return nil
}
