package setup

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceCreationStep struct {
	ClientSet *kubernetes.Clientset `json:"client_set"`
	Namespace string                `json:"namespace"`
}

func NewNamespaceCreationStep(clientSet *kubernetes.Clientset, namespace string) NamespaceCreationStep {
	newProvisioner := NamespaceCreationStep{
		ClientSet: clientSet,
		Namespace: namespace,
	}
	return newProvisioner
}

func (n NamespaceCreationStep) createNamespace() error {
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

func (n NamespaceCreationStep) GetName() string {
	return fmt.Sprintf("Create new namespace %s", n.Namespace)
}

func (n NamespaceCreationStep) PerformSetupStep() error {
	return n.createNamespace()
}
