package setup

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	secretNameDoguRegistry   = "dogu-cloudogu-com"
	secretNameDockerRegistry = "registry-cloudogu-com"
)

// secretCreator contains necessary information to create a new secrets in the cluster.
type secretCreator struct {
	ClientSet                 kubernetes.Interface `json:"client_set"`
	TargetNamespace           string               `json:"target_namespace"`
	CredentialSourceNamespace string               `yaml:"credential_source_namespace"`
}

// newSecretCreator creates a new object of type secretCreator.
func newSecretCreator(clientSet kubernetes.Interface, targetNamespace, credentialSourceNamespace string) *secretCreator {
	newProvisioner := &secretCreator{
		ClientSet:                 clientSet,
		TargetNamespace:           targetNamespace,
		CredentialSourceNamespace: credentialSourceNamespace,
	}
	return newProvisioner
}

func (sc *secretCreator) createSecrets() error {
	setupSecretInterface := sc.ClientSet.CoreV1().Secrets(sc.CredentialSourceNamespace)
	options := metav1.GetOptions{}
	dccSecret, err := setupSecretInterface.Get(context.Background(), secretNameDoguRegistry, options)
	if err != nil {
		return fmt.Errorf("failed to get dogu registry secret: %w", err)
	}
	dockerSecret, err := setupSecretInterface.Get(context.Background(), secretNameDockerRegistry, options)
	if err != nil {
		return fmt.Errorf("failed to get docker image pull secret: %w", err)
	}

	targetNamespace := sc.TargetNamespace
	destDccSecret := &v1.Secret{
		TypeMeta:   dccSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretNameDoguRegistry, Namespace: targetNamespace},
		Immutable:  dccSecret.Immutable,
		Data:       dccSecret.Data,
		StringData: dccSecret.StringData,
		Type:       dccSecret.Type,
	}
	_, err = sc.ClientSet.CoreV1().Secrets(targetNamespace).Create(context.Background(), destDccSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create dogu registry secret: %w", err)
	}

	destDockerSecret := &v1.Secret{
		TypeMeta:   dockerSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretNameDockerRegistry, Namespace: targetNamespace},
		Immutable:  dockerSecret.Immutable,
		Data:       dockerSecret.Data,
		StringData: dockerSecret.StringData,
		Type:       dockerSecret.Type,
	}
	_, err = sc.ClientSet.CoreV1().Secrets(targetNamespace).Create(context.Background(), destDockerSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create docker image pull secret: %w", err)
	}

	return nil
}

// GetStepDescription returns the description of the namespace creation step.
func (sc *secretCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new Secrets in namespace %s", sc.TargetNamespace)
}

// PerformSetupStep creates a namespace during setup execution.
func (sc *secretCreator) PerformSetupStep() error {
	return sc.createSecrets()
}
