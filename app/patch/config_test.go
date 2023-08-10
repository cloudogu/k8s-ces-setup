package patch

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestResourceReference_GroupVersionKind(t *testing.T) {
	t.Run("should return GVK for a core API", func(t *testing.T) {
		// given

		sut := ResourceReference{
			ApiVersion: "v1",
			Kind:       "Pod",
			Name:       "my-pod",
		}

		// when
		actual := sut.GroupVersionKind()

		// then
		assert.Equal(t, schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}, actual)
	})

	t.Run("should return GVK for a grouped API", func(t *testing.T) {
		// given

		sut := ResourceReference{
			ApiVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       "my-deploy",
		}

		// when
		actual := sut.GroupVersionKind()

		// then
		assert.Equal(t, schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    "Deployment",
		}, actual)
	})
}
