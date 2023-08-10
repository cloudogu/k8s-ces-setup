package patch

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type gvrMapper interface {
	meta.RESTMapper
}

type dynClient interface {
	dynamic.Interface
}

type jsonPatchApplier interface {
	Patch(jsonPatch []byte, gvk schema.GroupVersionKind, name string) error
}
