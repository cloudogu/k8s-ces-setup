package setup

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

	"k8s.io/client-go/kubernetes"
)

const (
	namespaceEnv             = "NAMESPACE"
	secretNameDoguRegistry   = "dogu-cloudogu-com"
	secretNameDockerRegistry = "docker-image-pull"
)

// secretCreator contains necessary information to create a new secrets in the cluster.
type secretCreator struct {
	ClientSet kubernetes.Interface `json:"client_set"`
	Namespace string               `json:"namespace"`
}

// newSecretCreator creates a new object of type secretCreator.
func newSecretCreator(clientSet kubernetes.Interface, namespace string) *secretCreator {
	newProvisioner := &secretCreator{
		ClientSet: clientSet,
		Namespace: namespace,
	}
	return newProvisioner
}

func (n *secretCreator) createSecrets() error {
	actualNamespace, err := getEnvVar(namespaceEnv)
	if err != nil {
		return fmt.Errorf("failed to read actual namespace")
	}

	setupSecretInterface := n.ClientSet.CoreV1().Secrets(actualNamespace)
	options := metav1.GetOptions{}
	dccSecret, err := setupSecretInterface.Get(context.Background(), secretNameDoguRegistry, options)
	if err != nil {
		return fmt.Errorf("failed to get dogu registry secret: %w", err)
	}
	dockerSecret, err := setupSecretInterface.Get(context.Background(), secretNameDockerRegistry, options)
	if err != nil {
		return fmt.Errorf("failed to get docker image pull secret: %w", err)
	}

	destDccSecret := &v1.Secret{
		TypeMeta:   dccSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretNameDoguRegistry, Namespace: n.Namespace},
		Immutable:  dccSecret.Immutable,
		Data:       dccSecret.Data,
		StringData: dccSecret.StringData,
		Type:       dccSecret.Type,
	}
	_, err = n.ClientSet.CoreV1().Secrets(n.Namespace).Create(context.Background(), destDccSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create dogu registry secret: %w", err)
	}

	destDockerSecret := &v1.Secret{
		TypeMeta:   dockerSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretNameDockerRegistry, Namespace: n.Namespace},
		Immutable:  dockerSecret.Immutable,
		Data:       dockerSecret.Data,
		StringData: dockerSecret.StringData,
		Type:       dockerSecret.Type,
	}
	_, err = n.ClientSet.CoreV1().Secrets(n.Namespace).Create(context.Background(), destDockerSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create docker image pull secret: %w", err)
	}

	return nil
}

// GetStepDescription returns the description of the namespace creation step.
func (n *secretCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new Secrets in namespace %s", n.Namespace)
}

// PerformSetupStep creates a namespace during setup execution.
func (n *secretCreator) PerformSetupStep() error {
	return n.createSecrets()
}

// getEnvVar returns the namespace the operator should be watching for changes
func getEnvVar(name string) (string, error) {
	ns, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	return ns, nil
}
