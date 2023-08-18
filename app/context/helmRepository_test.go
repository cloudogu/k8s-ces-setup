package context

import (
	"context"
	_ "embed"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmRepositoryData_GetOciEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		Endpoint string
		want     string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "success getOciEndpoint",
			Endpoint: "https://staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "success getOciEndpoint with Path",
			Endpoint: "https://staging-registry.cloudogu.com/foo/bar",
			want:     "oci://staging-registry.cloudogu.com/foo/bar",
			wantErr:  assert.NoError,
		},
		{
			name:     "success getOciEndpoint with other protocol",
			Endpoint: "ftp://staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "success no protocol",
			Endpoint: "staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "error empty string",
			Endpoint: "",
			want:     "",
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hrd := &HelmRepositoryData{
				Endpoint: tt.Endpoint,
			}
			got, err := hrd.GetOciEndpoint()
			if !tt.wantErr(t, err, fmt.Sprintf("GetOciEndpoint()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOciEndpoint()")
		})
	}
}

func TestReadHelmRepositoryDataFromCluster(t *testing.T) {
	t.Run("read without error", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		configMap := &v1.ConfigMap{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: "component-operator-helm-repository", Namespace: "testNS"},
			Immutable:  nil,
			Data: map[string]string{
				"endpoint": "testendpoint",
			},
		}
		clientset := fake.NewSimpleClientset(configMap)

		// when
		repoData, err := ReadHelmRepositoryDataFromCluster(testCtx, clientset, "testNS")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, repoData)
		assert.Equal(t, "testendpoint", repoData.Endpoint)
	})

	t.Run("not found error", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		configMap := &v1.ConfigMap{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "testNS"},
			Immutable:  nil,
			Data: map[string]string{
				"endpoint": "testendpoint",
			},
		}
		clientset := fake.NewSimpleClientset(configMap)

		// when
		_, err := ReadHelmRepositoryDataFromCluster(testCtx, clientset, "testNS")

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "helm repository configMap component-operator-helm-repository not found:")
	})
}

func TestReadHelmRepositoryDataFromFile(t *testing.T) {
	t.Run("read without error", func(t *testing.T) {
		// given

		// when
		repoData, err := ReadHelmRepositoryDataFromFile("testdata/testHelmRepoData.yaml")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, repoData)
		assert.Equal(t, "http://192.168.56.3:30100", repoData.Endpoint)
	})

	t.Run("error not found", func(t *testing.T) {
		// given

		// when
		_, err := ReadHelmRepositoryDataFromFile("invalid.yaml")

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "could not find configuration at invalid.yaml")
	})

	t.Run("error failes to unmarshal", func(t *testing.T) {
		// given

		// when
		_, err := ReadHelmRepositoryDataFromFile("helmRepository_test.go")

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal configuration helmRepository_test.go")
	})
}
