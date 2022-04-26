package core

import (
	"context"
	"fmt"
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
	gvrMapper meta.RESTMapper
	dynClient dynamic.Interface
}

// NewK8sClient creates a `kubectl`-like client which operates on the K8s API with YAML resources.
func NewK8sClient(clusterConfig *rest.Config) (*k8sApplyClient, error) {
	gvrMapper, err := createGVRMapper(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating k8s apply client: %w", err)
	}
	dynCli, err := createDynamicClient(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating k8s apply client: %w", err)
	}

	return &k8sApplyClient{gvrMapper: gvrMapper, dynClient: dynCli}, nil
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
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dyn, nil
}

// Apply sends a request to the K8s API with the provided YAML resources in order to apply them to the current cluster.
func (dkc *k8sApplyClient) Apply(yamlResources []byte, namespace string) error {
	logrus.Debug("Applying K8s resources")
	logrus.Debug(string(yamlResources))

	// 3. Decode YAML manifest into unstructured.Unstructured
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	k8sObjects := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yamlResources, nil, k8sObjects)
	if err != nil {
		return err
	}

	// 4. Map GVK to GVR
	// a resource can be uniquely identified by GroupVersionResource, but we need the GVK to find the corresponding GVR
	gvr, err := dkc.gvrMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if gvr.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dkc.dynClient.Resource(gvr.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources
		dr = dkc.dynClient.Resource(gvr.Resource)
	}

	return createOrUpdateResource(context.Background(), k8sObjects, dr)
}

func createOrUpdateResource(ctx context.Context, desiredResource *unstructured.Unstructured, dr dynamic.ResourceInterface) error {
	const fieldManager = "k8s-ces-setup"

	logrus.Debugf("Patching resource %s/%s/%s", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())
	// 6. marshal unstructured resource into proper JSON
	jsondata, err := json.Marshal(desiredResource)
	if err != nil {
		return NewResourceError(err, "error while parsing resource to json", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())
	}
	// 7. Update the object with server-side-apply
	//    types.ApplyPatchType indicates server-side-apply.
	//    FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, desiredResource.GetName(), types.ApplyPatchType, jsondata, metav1.PatchOptions{
		FieldManager: fieldManager,
	})
	if err != nil {
		return NewResourceError(err, "error while patching", desiredResource.GetKind(), desiredResource.GetAPIVersion(), desiredResource.GetName())
	}

	return nil
}
