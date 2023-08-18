package context

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
)

// HelmRepositoryData contains all necessary data for the helm repository.
type HelmRepositoryData struct {
	Endpoint string `json:"endpoint"`
}

// GetOciEndpoint returns the configured endpoint of the HelmRepositoryData with the OCI-protocol
func (hrd *HelmRepositoryData) GetOciEndpoint() (string, error) {
	split := strings.Split(hrd.Endpoint, "://")
	if len(split) == 1 && split[0] != "" {
		return fmt.Sprintf("oci://%s", split[0]), nil
	}
	if len(split) == 2 && split[1] != "" {
		return fmt.Sprintf("oci://%s", split[1]), nil
	}

	return "", fmt.Errorf("error creating oci-endpoint from '%s': wrong format", hrd.Endpoint)
}

// ReadHelmRepositoryDataFromCluster reads the helm repository data from the kubernetes configmap.
func ReadHelmRepositoryDataFromCluster(ctx context.Context, client kubernetes.Interface, namespace string) (*HelmRepositoryData, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, HelmRepositoryConfigMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("helm repository configMap %s not found: %w", HelmRepositoryConfigMapName, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get helm repository configMap %s: %w", HelmRepositoryConfigMapName, err)
	}

	return &HelmRepositoryData{
		Endpoint: configMap.Data["endpoint"],
	}, nil
}

// ReadHelmRepositoryDataFromFile reads the helm repository data from a yaml file.
func ReadHelmRepositoryDataFromFile(path string) (*HelmRepositoryData, error) {
	data := &HelmRepositoryData{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return data, fmt.Errorf("could not find configuration at %s", path)
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf("failed to read configuration %s: %w", path, err)
	}

	err = yaml.Unmarshal(fileData, data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal configuration %s: %w", path, err)
	}

	return data, nil
}
