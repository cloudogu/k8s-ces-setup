package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
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
func (dkc *k8sApplyClient) Apply(yamlResources []byte) error {
	logrus.Debug("Applying K8s resources")
	logrus.Debug(string(yamlResources))

	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	ctx := context.Background()

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

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(k8sObjects.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Create(ctx, k8sObjects, metav1.CreateOptions{
		FieldManager: "sample-controller",
	})

	return err
}
