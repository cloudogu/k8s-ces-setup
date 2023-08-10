package patch

import (
	"context"
	"errors"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type Applier struct {
	gvrMapper gvrMapper
	dynClient dynClient
	namespace string
}

func New(clusterConfig *rest.Config, crdSchemeBuilders []runtime.SchemeBuilder) (*Applier, error) {
	gvrMapper, err := createGVRMapper(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating GVR mapper: %w", err)
	}

	dynCli, err := createDynamicClient(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating dynamic client: %w", err)
	}

	err = handleCrds(crdSchemeBuilders)
	if err != nil {
		return nil, fmt.Errorf("failed to add crds to scheme: %w", err)
	}

	namespace := os.Getenv("POD_NAMESPACE")

	return &Applier{
		gvrMapper: gvrMapper,
		dynClient: dynCli,
		namespace: namespace,
	}, nil
}

func handleCrds(crdSchemeBuilders []runtime.SchemeBuilder) error {
	schemeForCrdHandling := runtime.NewScheme()
	var errs []error
	for _, builder := range crdSchemeBuilders {
		err := builder.AddToScheme(schemeForCrdHandling)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func createGVRMapper(config *rest.Config) (meta.RESTMapper, error) {
	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc)), nil
}

func createDynamicClient(config *rest.Config) (dynamic.Interface, error) {
	// 2. Prepare the dynamic client
	return dynamic.NewForConfig(config)
}

func (ac *Applier) Patch(jsonPatch []byte, gvk schema.GroupVersionKind, name string) error {
	// 4. Map GVK to GVR
	// a resource can be uniquely identified by GroupVersionResource, but we need the GVK to find the corresponding GVR
	gvr, err := ac.gvrMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could not find GVK mapper for GroupKind=%v,Version=%s: %w", gvk.GroupKind(), gvk.Version, err)
	}

	var dr dynamic.ResourceInterface
	if gvr.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = ac.dynClient.Resource(gvr.Resource).Namespace(ac.namespace)
	} else {
		// for cluster-wide resources
		dr = ac.dynClient.Resource(gvr.Resource)
	}

	err = ac.patchResource(context.Background(), name, jsonPatch, dr)
	if err != nil {
		return fmt.Errorf("failed to patch resource %s of kind %s with json patch '%s'", name, gvk, jsonPatch)
	}

	return nil
}

func (ac *Applier) patchResource(ctx context.Context, name string, jsonPatch []byte, dr dynamic.ResourceInterface) error {
	_, err := dr.Patch(ctx, name, types.JSONPatchType, jsonPatch, v1.PatchOptions{})
	return err
}
