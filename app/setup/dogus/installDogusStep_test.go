package dogus_test

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-ces-setup/app/setup/dogus"

	"k8s.io/client-go/rest/fake"
)

type fakeCodec struct{}

func (c *fakeCodec) Decode([]byte, *schema.GroupVersionKind, runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	return nil, nil, nil
}

func (c *fakeCodec) Encode(_ runtime.Object, _ io.Writer) error {
	return nil
}

func (c *fakeCodec) Identifier() runtime.Identifier {
	return "fake"
}

type fakeNegotiatedSerializer struct {
	mediaType string
}

func (n *fakeNegotiatedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	myMediaTypes := []runtime.SerializerInfo{
		{MediaType: "application/json"},
	}
	return myMediaTypes
}

func (n *fakeNegotiatedSerializer) EncoderForVersion(_ runtime.Encoder, _ runtime.GroupVersioner) runtime.Encoder {
	return &fakeCodec{}
}

func (n *fakeNegotiatedSerializer) DecoderToVersion(_ runtime.Decoder, _ runtime.GroupVersioner) runtime.Decoder {
	return &fakeCodec{}
}

func TestNewInstallDogusStep(t *testing.T) {
	t.Run("create new dogu install step", func(t *testing.T) {
		// given
		fakeRestClient := &fake.RESTClient{}
		restClientMock := fake.CreateHTTPClient(func(request *http.Request) (res *http.Response, err error) {

			return
		})
		fakeRestClient.Client = restClientMock
		myDogu := &core.Dogu{Name: "MyName"}

		// when
		installStep := dogus.NewInstallDogusStep(fakeRestClient, myDogu, "namespace")

		// then
		require.NotNil(t, installStep)
	})
}

func Test_installDogusStep_GetStepDescription(t *testing.T) {
	t.Run("create new dogu install step", func(t *testing.T) {
		// given
		fakeRestClient := &fake.RESTClient{}
		restClientMock := fake.CreateHTTPClient(func(request *http.Request) (res *http.Response, err error) {

			return
		})
		fakeRestClient.Client = restClientMock
		myDogu := &core.Dogu{Name: "MyName"}
		installStep := dogus.NewInstallDogusStep(fakeRestClient, myDogu, "namespace")

		// when
		description := installStep.GetStepDescription()

		// then
		assert.Equal(t, "Installing dogu [MyName]", description)
	})
}

func Test_installDogusStep_PerformSetupStep(t *testing.T) {
	t.Run("failed to get version", func(t *testing.T) {
		// given
		fakeRestClient := &fake.RESTClient{}
		restClientMock := fake.CreateHTTPClient(func(request *http.Request) (res *http.Response, err error) {

			return
		})
		fakeRestClient.Client = restClientMock
		myDogu := &core.Dogu{Name: "MyName", Version: "-----------"}
		installStep := dogus.NewInstallDogusStep(fakeRestClient, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get version from dogu")
	})

	t.Run("failed on doing post request", func(t *testing.T) {
		// given
		fakeRestClient := &fake.RESTClient{}
		restClientMock := fake.CreateHTTPClient(func(request *http.Request) (res *http.Response, err error) {
			res = &http.Response{Status: strconv.Itoa(http.StatusOK)}
			return
		})
		fakeRestClient.Client = restClientMock
		fakeRestClient.NegotiatedSerializer = &fakeNegotiatedSerializer{}

		//fakeRestClient.N
		myDogu := &core.Dogu{Name: "MyName", Version: "1.1.1-1"}
		installStep := dogus.NewInstallDogusStep(fakeRestClient, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to apply dogu MyName")
	})

	t.Run("successfully apply dogu cr", func(t *testing.T) {
		// given
		fakeRestClient := &fake.RESTClient{}
		restClientMock := fake.CreateHTTPClient(func(request *http.Request) (res *http.Response, err error) {
			res = &http.Response{Status: strconv.Itoa(http.StatusOK)}
			return
		})
		fakeRestClient.Client = restClientMock
		fakeRestClient.NegotiatedSerializer = &fakeNegotiatedSerializer{}

		//fakeRestClient.N
		myDogu := &core.Dogu{Name: "MyName", Version: "1.1.1-1"}
		installStep := dogus.NewInstallDogusStep(fakeRestClient, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep()

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to apply dogu MyName")
	})
}
