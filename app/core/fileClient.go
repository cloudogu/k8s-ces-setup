package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type defaultHttpClient struct {
	httpClient *http.Client
}

func NewFileClient(appVersion string) *defaultHttpClient {
	return &defaultHttpClient{
		httpClient: &http.Client{},
	}
}

// Get retrieves a file over HTTP returns its contents.
func (dhc *defaultHttpClient) Get(url string) ([]byte, error) {
	resp, err := dhc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not find YAML file '%s'", url)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bytes, nil
}
