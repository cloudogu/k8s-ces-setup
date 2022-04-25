package setup

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// instanceSecretValidator validates whether the target namespace contains necessary structure or data.
type instanceSecretValidator struct {
	clientSet       kubernetes.Interface
	targetNamespace string
}

// newInstanceSecretValidator creates a new object of type instanceSecretValidator.
func newInstanceSecretValidator(clientSet kubernetes.Interface, targetNamespace string) *instanceSecretValidator {
	v := &instanceSecretValidator{
		clientSet:       clientSet,
		targetNamespace: targetNamespace,
	}
	return v
}

func (n *instanceSecretValidator) validate() error {
	_, err := n.clientSet.CoreV1().Secrets(n.targetNamespace).Get(context.Background(), secretNameDoguRegistry, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("instance secret validation error: cannot read secret from target namespace %s: %w", n.targetNamespace, err)
	}
	_, err = n.clientSet.CoreV1().Secrets(n.targetNamespace).Get(context.Background(), secretNameImageRegistry, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("instance secret validation error: cannot read secret from target namespace %s: %w", n.targetNamespace, err)
	}

	return nil
}
