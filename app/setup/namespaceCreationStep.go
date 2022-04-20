package setup

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// namespaceCreator contains necessary information to create a new namespace in the cluster.
type namespaceCreator struct {
	clientSet       kubernetes.Interface
	targetNamespace string
}

// newNamespaceCreator creates a new object of type namespaceCreator.
func newNamespaceCreator(clientSet kubernetes.Interface, targetNamespace string) *namespaceCreator {
	newProvisioner := &namespaceCreator{
		clientSet:       clientSet,
		targetNamespace: targetNamespace,
	}
	return newProvisioner
}

func (n *namespaceCreator) createNamespace() error {
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: n.targetNamespace,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}

	_, err := n.clientSet.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("cannot create namespace %s with clientset: %w", n.targetNamespace, err)
	}

	return nil
}

// GetStepDescription returns the description of the namespace creation step.
func (n *namespaceCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new namespace %s", n.targetNamespace)
}

// PerformSetupStep creates a namespace during setup execution.
func (n *namespaceCreator) PerformSetupStep() error {
	return n.createNamespace()
}
