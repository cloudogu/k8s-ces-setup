package component

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	ctx "github.com/cloudogu/k8s-ces-setup/app/context"
)

// NewEtcdClientInstallerStep creates step to install the etcd client.
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
func (ecis *etcdClientInstallerStep) PerformSetupStep(ctx context.Context) error {
	err := ecis.installClient(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ecis *etcdClientInstallerStep) installClient(ctx context.Context) error {
	err := ecis.createDeployment(ctx, ecis.etcdServiceUrl)
	if err != nil {
		return err
	}

	return nil
}

func (ecis *etcdClientInstallerStep) createDeployment(ctx context.Context, etcdServiceUrl string) error {
	etcdClientName := "etcd-client"
	const etcdAPIVersion = "2"
	etcdClientLabels := make(map[string]string)
	etcdClientLabels["app"] = "ces"
	etcdClientLabels["app.kubernetes.io/name"] = etcdClientName

	deployment := ecis.getEtcdClientDeployment(etcdServiceUrl, etcdClientName, etcdClientLabels, etcdAPIVersion)

	_, err := ecis.clientSet.AppsV1().Deployments(ecis.targetNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("cannot create etcd deployment in namespace %s with clientset: %w", ecis.targetNamespace, err)
	}

	return nil
}

func (ecis *etcdClientInstallerStep) getEtcdClientDeployment(etcdServiceUrl string, etcdClientName string,
	etcdClientLabels map[string]string, etcdAPIVersion string) *appsv1.Deployment {
	deploymentObjMeta := metav1.ObjectMeta{
		Name:      etcdClientName,
		Namespace: ecis.targetNamespace,
		Labels:    etcdClientLabels,
	}

	podObjMeta := metav1.ObjectMeta{
		Name:   etcdClientName,
		Labels: etcdClientLabels,
	}

	replicas := int32(1)

	podSpec := corev1.PodSpec{
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
	}

	deployment := &appsv1.Deployment{ObjectMeta: deploymentObjMeta, Spec: appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{MatchLabels: etcdClientLabels},
		Replicas: &replicas,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: podObjMeta,
			Spec:       podSpec,
		},
	}}

	return deployment
}
