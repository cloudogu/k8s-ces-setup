package setup

import (
	"context"
	"fmt"
	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func newEtcdClientInstallerStep(clientSet kubernetes.Interface, setupCtx ctx.SetupContext) *etcdClientInstallerStep {
	return &etcdClientInstallerStep{
		ClientSet:       clientSet,
		targetNamespace: setupCtx.AppConfig.TargetNamespace,
		imageURL:        setupCtx.AppConfig.EtcdClientImageRepo,
	}
}

type etcdClientInstallerStep struct {
	ClientSet       kubernetes.Interface `json:"client_set"`
	targetNamespace string
	imageURL        string
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
	err := ecis.createPod()
	if err != nil {
		return err
	}

	return nil
}

func (ecis *etcdClientInstallerStep) createPod() error {
	etcdClientName := "etcd-client"
	etcdClientLabels := make(map[string]string, 0)
	etcdClientLabels["run"] = etcdClientName
	mountServiceAccountToken := true
	const etcdAPIVersion = "2"

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
							Name:  "ROOT_PASSWORD",
							Value: "",
						},
						{
							Name:  "ETCDCTL_ENDPOINTS",
							Value: "",
						},
					},
				},
			},
			AutomountServiceAccountToken: &mountServiceAccountToken,
		},
		Status: corev1.PodStatus{},
	}

	_, err := ecis.ClientSet.CoreV1().Pods(ecis.targetNamespace).Create(context.Background(), etcdPod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("cannot create etcd pod in namespace %s with clientset: %w", ecis.targetNamespace, err)
	}
	return nil
}
