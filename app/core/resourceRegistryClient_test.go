package core

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewResourceRegistryClient(t *testing.T) {
	t.Run("successful init", func(t *testing.T) {
		// given
		doguRegistrySecret := &context.DoguRegistrySecret{
			Endpoint: "endpoint",
			Username: "username",
			Password: "password",
		}

		// when
		sut := NewResourceRegistryClient("1.0.0", doguRegistrySecret)

		// then
		require.NotNil(t, sut)
		assert.Equal(t, doguRegistrySecret, sut.doguRegistrySecret)
		assert.NotNil(t, sut.fileClient)
	})
}

func TestResourceRegistryClient_GetResourceFileContent(t *testing.T) {
	t.Run("should return byte slice on success with same host as dogu registry", func(t *testing.T) {
		// given
		doguRegistrySecret := &context.DoguRegistrySecret{
			Endpoint: "https://url.de/api/v23/k8s",
			Username: "username",
			Password: "password",
		}
		sut := NewResourceRegistryClient("1.0.0", doguRegistrySecret)
		fileClientMock := &mocks.FileClient{}
		byteResult := []byte("test")
		fileClientMock.On("Get", "https://url.de/api/v23/k8s/component/0.23.0", "username", "password").Return(byteResult, nil)
		sut.fileClient = fileClientMock

		// when
		content, err := sut.GetResourceFileContent("https://url.de/api/v23/k8s/component/0.23.0")

		// then
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, byteResult, content)
		fileClientMock.AssertExpectations(t)
	})

	t.Run("should return byte slice on success with different host as dogu registry", func(t *testing.T) {
		// given
		doguRegistrySecret := &context.DoguRegistrySecret{
			Endpoint: "https://url123.de/api/v23/k8s",
			Username: "username",
			Password: "password",
		}
		sut := NewResourceRegistryClient("1.0.0", doguRegistrySecret)
		fileClientMock := &mocks.FileClient{}
		byteResult := []byte("test")
		fileClientMock.On("Get", "https://url.de/api/v23/k8s/component/0.23.0", "", "").Return(byteResult, nil)
		sut.fileClient = fileClientMock

		// when
		content, err := sut.GetResourceFileContent("https://url.de/api/v23/k8s/component/0.23.0")

		// then
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.Equal(t, byteResult, content)
		fileClientMock.AssertExpectations(t)
	})

	t.Run("should return error when get an error from the file client", func(t *testing.T) {
		// given
		doguRegistrySecret := &context.DoguRegistrySecret{
			Endpoint: "https://url123.de/api/v23/k8s",
			Username: "username",
			Password: "password",
		}
		sut := NewResourceRegistryClient("1.0.0", doguRegistrySecret)
		fileClientMock := &mocks.FileClient{}
		fileClientMock.On("Get", "https://url.de/api/v23/k8s/component/0.23.0", "", "").Return(nil, assert.AnError)
		sut.fileClient = fileClientMock

		// when
		content, err := sut.GetResourceFileContent("https://url.de/api/v23/k8s/component/0.23.0")

		// then
		require.Error(t, err)
		assert.Nil(t, content)
		assert.Contains(t, err.Error(), "failed to get file content from https://url.de/api/v23/k8s/component/0.23.0")
		fileClientMock.AssertExpectations(t)
	})
}
