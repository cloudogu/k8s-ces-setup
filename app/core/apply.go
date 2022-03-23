package core

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
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

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources

		dr = dyn.Resource(mapping.Resource)
	}

	ctx := context.Background()
	resourceExists, err := existsResource(ctx, dr, k8sObjects)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with server-side-apply
	//     types.ApplyPatchType indicates server-side-apply.
	//     FieldManager specifies the field owner ID.
	return createOrUpdateResource(resourceExists, ctx, dr, k8sObjects)
}

func createOrUpdateResource(resourceExists bool, ctx context.Context, dr dynamic.ResourceInterface, k8sObjects *unstructured.Unstructured) error {
	if resourceExists {
		logrus.Debugf("Updating resource %s/%s/%s", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
		_, err := dr.Update(ctx, k8sObjects, metav1.UpdateOptions{
			FieldManager: "k8s-ces-setup",
		})
		if err != nil {
			return NewResourceError(err, "error while updating", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
		}

		return nil
	}

	logrus.Debugf("Creating resource %s/%s/%s", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
	_, err := dr.Create(ctx, k8sObjects, metav1.CreateOptions{
		FieldManager: "k8s-ces-setup",
	})
	if err != nil {
		return NewResourceError(err, "error while creating", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
	}

	return nil
}

func existsResource(ctx context.Context, dr dynamic.ResourceInterface, k8sObjects *unstructured.Unstructured) (bool, error) {
	_, err := dr.Get(ctx, k8sObjects.GetName(), metav1.GetOptions{})
	if err != nil {
		var typedErr *k8serr.StatusError
		if errors.As(err, &typedErr) && typedErr.Status().Reason == metav1.StatusReasonNotFound {
			logrus.Debugf("Resource %s/%s/%s does not exist yet", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
			return false, nil
		}

		return false, NewResourceError(err, "error while getting", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
	}

	logrus.Debugf("Resource %s/%s/%s already exists", k8sObjects.GetKind(), k8sObjects.GetAPIVersion(), k8sObjects.GetName())
	return true, nil
}
