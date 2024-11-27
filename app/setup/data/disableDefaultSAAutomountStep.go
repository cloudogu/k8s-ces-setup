package data

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type disableDefaultSAAutomountStep struct {
	clientSet kubernetes.Interface
	namespace string
}

// NewDisableDefaultSAAutomountStep disables automounting the token for the default service account in the ecosystem namespace
func NewDisableDefaultSAAutomountStep(clientSet kubernetes.Interface, namespace string) *disableDefaultSAAutomountStep {
	return &disableDefaultSAAutomountStep{clientSet: clientSet, namespace: namespace}
}

// GetStepDescription return the human-readable description of the step
func (fcs *disableDefaultSAAutomountStep) GetStepDescription() string {
	return "Disable automounting the token for the default service account in the ecosystem namespace"
}

// PerformSetupStep disables automounting the token for the default service account in the ecoSystem namespace.
// No pod in the ecosystem namespace should mount the default service account, with this we prevent accidental mounting.
func (fcs *disableDefaultSAAutomountStep) PerformSetupStep(ctx context.Context) error {
	serviceAccount, err := fcs.clientSet.CoreV1().ServiceAccounts(fcs.namespace).Get(ctx, "default", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to get default service account: %w", err)
	}
	var automountServiceAccountToken = false
	serviceAccount.AutomountServiceAccountToken = &automountServiceAccountToken
	_, err = fcs.clientSet.CoreV1().ServiceAccounts(fcs.namespace).Update(ctx, serviceAccount, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("unable to deactivate token automount on default service account: %w", err)
	}
	return nil
}
