package component

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const nodeMasterFileConfigMapName = "node-master-file"

type nodeMasterCreationStep struct {
	ClientSet       kubernetes.Interface `json:"client_set"`
	TargetNamespace string               `json:"target_namespace"`
}

// NewNodeMasterCreationStep create a new setup step responsible to create a node master config map pointing to the
// current etcd server.
func NewNodeMasterCreationStep(clusterConfig *rest.Config, targetNamespace string) (*nodeMasterCreationStep, error) {
	client, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	return &nodeMasterCreationStep{
		ClientSet:       client,
		TargetNamespace: targetNamespace,
	}, nil
}

// GetStepDescription returns a human-readable description of the node master file creation step.
func (nmcs *nodeMasterCreationStep) GetStepDescription() string {
	return "Setup node master file"
}

// PerformSetupStep creates a config map containing the node master address.
func (nmcs *nodeMasterCreationStep) PerformSetupStep() error {
	nodeMasterFileContent := fmt.Sprintf("etcd.%s.svc.cluster.local", nmcs.TargetNamespace)

	configMap, err := nmcs.ClientSet.CoreV1().ConfigMaps(nmcs.TargetNamespace).Get(context.Background(), nodeMasterFileConfigMapName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to get config map [%s]: %w", nodeMasterFileConfigMapName, err)
	}

	if errors.IsNotFound(err) {
		// create new
		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: nodeMasterFileConfigMapName},
			Data: map[string]string{
				"node_master": nodeMasterFileContent,
			},
		}

		_, err = nmcs.ClientSet.CoreV1().ConfigMaps(nmcs.TargetNamespace).Create(context.Background(), configMap, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create config map [%s]: %w", nodeMasterFileConfigMapName, err)
		}
	} else {
		configMap.Data = map[string]string{
			"node_master": nodeMasterFileContent,
		}

		_, err = nmcs.ClientSet.CoreV1().ConfigMaps(nmcs.TargetNamespace).Update(context.Background(), configMap, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update config map [%s]: %w", nodeMasterFileConfigMapName, err)
		}
	}

	return nil
}
