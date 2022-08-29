package core

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"net/url"
)

type fileClient interface {
	// Get retrieves a file identified by its URL and returns the contents.
	Get(url string, username string, password string) ([]byte, error)
}

// ResourceRegistryClient is used to get resources with a file client from a dogu registry.
type ResourceRegistryClient struct {
	fileClient         fileClient
	doguRegistrySecret *context.DoguRegistrySecret
}

// NewResourceRegistryClient creates a new instance of ResourceRegistryClient
func NewResourceRegistryClient(appVersion string, secret *context.DoguRegistrySecret) *ResourceRegistryClient {
	return &ResourceRegistryClient{fileClient: NewFileClient(appVersion), doguRegistrySecret: secret}
}

// GetResourceFileContent gets the bytes from the given resource url. If the host of the url is the same as the host
// of the dogu registry, basic auth will be used.
func (rrc *ResourceRegistryClient) GetResourceFileContent(resourceURL string) ([]byte, error) {
	registryEndpoint := rrc.doguRegistrySecret.Endpoint
	registryUrl, err := url.Parse(registryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse registry endpoint %s to url: %w", rrc.doguRegistrySecret.Endpoint, err)
	}
	resourceUrl, err := url.Parse(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource url %s to url: %w", resourceURL, err)
	}

	username := ""
	password := ""
	if registryUrl.Host == resourceUrl.Host {
		username = rrc.doguRegistrySecret.Username
		password = rrc.doguRegistrySecret.Password
	}

	fileContent, err := rrc.fileClient.Get(resourceURL, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content from %s: %w", resourceURL, err)
	}

	return fileContent, nil
}
