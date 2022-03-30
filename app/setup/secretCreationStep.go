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
	err := sc.copySecretToTargetNamespace(secretNameDoguRegistry)
	if err != nil {
		return fmt.Errorf("failed to copy dogu registry secret: %w", err)
	}
	err = sc.copySecretToTargetNamespace(secretNameDockerRegistry)
	if err != nil {
		return fmt.Errorf("failed to copy docker image pull secret: %w", err)
	}

	return nil
}

// copySecretToTargetNamespace copies a secret from one namespace to another by rewriting the namespace.
func (sc *secretCreator) copySecretToTargetNamespace(secretName string) error {
	setupSecretInterface := sc.ClientSet.CoreV1().Secrets(sc.CredentialSourceNamespace)

	originalSecret, err := setupSecretInterface.Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	secretCopy := &v1.Secret{
		TypeMeta:   originalSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: sc.TargetNamespace},
		Immutable:  originalSecret.Immutable,
		Data:       originalSecret.Data,
		StringData: originalSecret.StringData,
		Type:       originalSecret.Type,
	}
	_, err = sc.ClientSet.CoreV1().Secrets(sc.TargetNamespace).Create(context.Background(), secretCopy, metav1.CreateOptions{})
	if err != nil {
		return err
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
