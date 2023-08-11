package patch

import (
	"context"

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
	// Patch applies a JSON patch to a Kubernetes resource.
	Patch(ctx context.Context, jsonPatch []byte, gvk schema.GroupVersionKind, name string) error
}
