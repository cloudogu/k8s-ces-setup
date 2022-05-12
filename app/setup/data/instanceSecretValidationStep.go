package data

import (
	"context"
	"fmt"

	setupContext "github.com/cloudogu/k8s-ces-setup/app/context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// instanceSecretValidatorStep validates whether the target namespace contains necessary structure or data.
type instanceSecretValidatorStep struct {
	clientSet       kubernetes.Interface
	targetNamespace string
}

// NewInstanceSecretValidatorStep creates a new object of type instanceSecretValidatorStep.
func NewInstanceSecretValidatorStep(clientSet kubernetes.Interface, targetNamespace string) *instanceSecretValidatorStep {
	v := &instanceSecretValidatorStep{
		clientSet:       clientSet,
		targetNamespace: targetNamespace,
	}
	return v
}

// GetStepDescription returns a human-readable description of the instance secrets validation step.
func (isv *instanceSecretValidatorStep) GetStepDescription() string {
	return "Validate instance secrets"
}

// PerformSetupStep validates the current instance secrets.
func (isv *instanceSecretValidatorStep) PerformSetupStep() error {
	_, err := isv.clientSet.CoreV1().Secrets(isv.targetNamespace).Get(context.Background(), setupContext.SecretDoguRegistry, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("instance secret validation error: cannot read secret from target namespace %s: %w", isv.targetNamespace, err)
	}
	_, err = isv.clientSet.CoreV1().Secrets(isv.targetNamespace).Get(context.Background(), setupContext.SecretDockerRegistry, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("instance secret validation error: cannot read secret from target namespace %s: %w", isv.targetNamespace, err)
	}

	return nil
}
