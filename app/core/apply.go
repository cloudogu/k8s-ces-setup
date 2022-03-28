package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type k8sApplyClient struct {
	clusterConfig *rest.Config
}

// NewK8sClient creates a `kubectl`-like client which operates on the K8s API.
func NewK8sClient(clusterConfig *rest.Config) *k8sApplyClient {
	return &k8sApplyClient{clusterConfig: clusterConfig}
}

// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster.
func (dkc *k8sApplyClient) Apply(yamlResources []byte, namespace string) error {
	logrus.Debug("Applying K8s resources")
	logrus.Debug(string(yamlResources))

	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(dkc.clusterConfig)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(dkc.clusterConfig)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	k8sObjects := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yamlResources, nil, k8sObjects)
	if err != nil {
		return err
	}

	// 4. Map GVK to GVR
	// a resource can be uniquely identified by GroupVersionResource, but we need the GVK to find the corresponding GVR
	gvr, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if gvr.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(gvr.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(gvr.Resource)
	}

	ctx := context.Background()
	if err != nil {
		return err
	}

	return createOrUpdateResource(k8sObjects, ctx, dr)
}

func createOrUpdateResource(desiredResource *unstructured.Unstructured, ctx context.Context, dr dynamic.ResourceInterface) error {
	const fieldManager = "k8s-ces-setup"

	logrus.Debugf("Patching resource %s/%s/%s", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())

	jsondata, err := json.Marshal(desiredResource)
	if err != nil {
		return NewResourceError(err, "error while parsing resource to json", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())
	}
	// 6a. Update the object with server-side-apply
	//     types.ApplyPatchType indicates server-side-apply.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, desiredResource.GetName(), types.ApplyPatchType, jsondata, metav1.PatchOptions{
		FieldManager: fieldManager,
	})
	if err != nil {
		return NewResourceError(err, "error while patching", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())
	}

	return nil
}
