package setup

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	secretNameDoguRegistry  = "cloudogu-dogu-registry"
	secretNameImageRegistry = "cloudogu-image-registry"
)

// secretCreator contains necessary information to create a new secrets in the cluster.
type secretCreator struct {
	clientSet                 kubernetes.Interface
	targetNamespace           string
	credentialSourceNamespace string
}

// newSecretCreator creates a new object of type secretCreator.
func newSecretCreator(clientSet kubernetes.Interface, targetNamespace, credentialSourceNamespace string) *secretCreator {
	newProvisioner := &secretCreator{
		clientSet:                 clientSet,
		targetNamespace:           targetNamespace,
		credentialSourceNamespace: credentialSourceNamespace,
	}
	return newProvisioner
}

func (sc *secretCreator) createSecrets() error {
	err := sc.copySecretToTargetNamespace(secretNameDoguRegistry)
	if err != nil {
		return fmt.Errorf("failed to copy dogu registry secret: %w", err)
	}
	err = sc.copySecretToTargetNamespace(secretNameImageRegistry)
	if err != nil {
		return fmt.Errorf("failed to copy docker image pull secret: %w", err)
	}

	return nil
}

// copySecretToTargetNamespace copies a secret from one namespace to another by rewriting the namespace.
func (sc *secretCreator) copySecretToTargetNamespace(secretName string) error {
	setupSecretInterface := sc.clientSet.CoreV1().Secrets(sc.credentialSourceNamespace)

	originalSecret, err := setupSecretInterface.Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	secretCopy := &v1.Secret{
		TypeMeta:   originalSecret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: sc.targetNamespace},
		Immutable:  originalSecret.Immutable,
		Data:       originalSecret.Data,
		StringData: originalSecret.StringData,
		Type:       originalSecret.Type,
	}
	_, err = sc.clientSet.CoreV1().Secrets(sc.targetNamespace).Create(context.Background(), secretCopy, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GetStepDescription returns the description of the namespace creation step.
func (sc *secretCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new Secrets in namespace %s", sc.targetNamespace)
}

// PerformSetupStep creates a namespace during setup execution.
func (sc *secretCreator) PerformSetupStep() error {
	return sc.createSecrets()
}
