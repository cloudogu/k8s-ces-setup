package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type defaultHttpClient struct {
	httpClient *http.Client
	version    string
}

// NewFileClient creates a http client to read files from the network.
func NewFileClient(appVersion string) *defaultHttpClient {
	return &defaultHttpClient{
		httpClient: &http.Client{},
		version:    appVersion,
	}
}

// Get retrieves a file over HTTP and returns its content.
func (dhc *defaultHttpClient) Get(url string, username string, password string) ([]byte, error) {
	logrus.Debugf("Getting resource from %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to %s: %w", url, err)
	}

	req.Header.Set("User-Agent", "Cloudia/"+dhc.version)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := dhc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during GET request to %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response for YAML file '%s' returned with non-200 reply (HTTP %d): identifying this as an error", url, resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading the response body of '%s'", url)
	}
	defer resp.Body.Close()

	return bytes, nil
}
