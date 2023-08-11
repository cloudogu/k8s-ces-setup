package context

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-apply-lib/apply"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

// DoguRegistrySecret defines the credentials and the endpoint for the dogu registry.
type DoguRegistrySecret struct {
	Endpoint  string `yaml:"endpoint"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	URLSchema string `yaml:"urlschema"`
}

// ReadDoguRegistrySecretFromCluster reads the dogu registry credentials from the kubernetes secret.
func ReadDoguRegistrySecretFromCluster(ctx context.Context, client kubernetes.Interface, namespace string) (*DoguRegistrySecret, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, SecretDoguRegistry, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("dogu registry secret %s not found: %w", SecretDoguRegistry, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get dogu registry secret %s: %w", SecretDoguRegistry, err)
	}
	urlSchema := string(secret.Data["urlschema"])
	if urlSchema != "index" {
		logger := apply.GetLogger()
		logger.Info("URLSchema is not index. Setting it to default.")
		urlSchema = "default"
	}
	return &DoguRegistrySecret{
		Endpoint:  string(secret.Data["endpoint"]),
		Username:  string(secret.Data["username"]),
		Password:  string(secret.Data["password"]),
		URLSchema: urlSchema,
	}, nil
}

// ReadDoguRegistrySecretFromFile reads the dogu registry credentials from a yaml file.
func ReadDoguRegistrySecretFromFile(path string) (*DoguRegistrySecret, error) {
	doguRegistry := &DoguRegistrySecret{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return doguRegistry, fmt.Errorf("could not find registry secret at %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return doguRegistry, fmt.Errorf("failed to read registry secret %s: %w", path, err)
	}

	err = yaml.Unmarshal(data, doguRegistry)
	if err != nil {
		return doguRegistry, fmt.Errorf("failed to unmarshal registry secret %s: %w", path, err)
	}

	return doguRegistry, nil
}
