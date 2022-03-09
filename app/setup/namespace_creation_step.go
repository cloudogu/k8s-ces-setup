package setup

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespaceCreator contains necessary information to create a new namespace in the cluster
type NamespaceCreator struct {
	ClientSet kubernetes.Interface `json:"clientSet"`
	Namespace string               `json:"namespace"`
}

// NewNamespaceCreator creates a new object of type NamespaceCreator
func NewNamespaceCreator(clientSet kubernetes.Interface, namespace string) NamespaceCreator {
	newProvisioner := NamespaceCreator{
		ClientSet: clientSet,
		Namespace: namespace,
	}
	return newProvisioner
}

func (n NamespaceCreator) createNamespace() error {
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
		return fmt.Errorf("cannot create namespace %s with clientset; %w", n.Namespace, err)
	}

	return nil
}

func (n NamespaceCreator) GetStepDescription() string {
	return fmt.Sprintf("Create new namespace %s", n.Namespace)
}

func (n NamespaceCreator) PerformSetupStep() error {
	return n.createNamespace()
}
