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
	ClientSet kubernetes.Interface `json:"client_set"`
	Namespace string               `json:"namespace"`
}

// newNamespaceCreator creates a new object of type namespaceCreator.
func newNamespaceCreator(clientSet kubernetes.Interface, namespace string) *namespaceCreator {
	newProvisioner := &namespaceCreator{
		ClientSet: clientSet,
		Namespace: namespace,
	}
	return newProvisioner
}

func (n *namespaceCreator) createNamespace() error {
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: n.Namespace,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}

	_, err := n.ClientSet.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("cannot create namespace %s with clientset: %w", n.Namespace, err)
	}

	return nil
}

// GetStepDescription returns the description of the namespace creation step.
func (n *namespaceCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new namespace %s", n.Namespace)
}

// PerformSetupStep creates a namespace during setup execution.
func (n *namespaceCreator) PerformSetupStep() error {
	return n.createNamespace()
}