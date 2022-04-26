package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type defaultHttpClient struct {
	httpClient *http.Client
	version    string
}

//NewFileClient creates a http client to read files from the network.
func NewFileClient(appVersion string) *defaultHttpClient {
	return &defaultHttpClient{
		httpClient: &http.Client{},
		version:    appVersion,
	}
}

// Get retrieves a file over HTTP and returns its content.
func (dhc *defaultHttpClient) Get(url string) ([]byte, error) {
	logrus.Debugf("Getting resource from %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Cloudia/"+dhc.version)

	resp, err := dhc.httpClient.Do(req)
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
