package core

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

// Apply sends a request to the K8s API with the provided YAML resource in order to apply them to the current cluster.
func (dkc *k8sApplyClient) Apply(yamlResource []byte, namespace string) error {
	return dkc.ApplyWithOwner(yamlResource, namespace, nil)
}

// ApplyWithOwner sends a request to the K8s API with the provided YAML resource in order to apply them to the current cluster.
func (dkc *k8sApplyClient) ApplyWithOwner(yamlResource []byte, namespace string, owningResource metav1.Object) error {
	logrus.Debug("Applying K8s resource")
	logrus.Debug(string(yamlResource))

	// 3. Decode YAML manifest into unstructured.Unstructured
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	k8sObjects := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yamlResource, nil, k8sObjects)
	if err != nil {
		return fmt.Errorf("could not decode YAML doccument '%s': %w", string(yamlResource), err)
	}

	// 4. Map GVK to GVR
	// a resource can be uniquely identified by GroupVersionResource, but we need the GVK to find the corresponding GVR
	gvr, err := dkc.gvrMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could find GVK mapper for GroupKind=%v,Version=%s and YAML doccument '%s': %w", gvk.GroupKind(), gvk.Version, string(yamlResource), err)
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if gvr.Scope.Name() == meta.RESTScopeNameNamespace {
		k8sObjects.SetNamespace(namespace)
		// namespaced resources should specify the namespace
		dr = dkc.dynClient.Resource(gvr.Resource).Namespace(namespace)

		if owningResource != nil {
			err = SetControllerReference(owningResource, k8sObjects, gvk)
			if err != nil {
				return fmt.Errorf("could not apply YAML doccument '%s': could not set controller reference: %w", string(yamlResource), err)
			}
		}
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

// SetControllerReference sets owner as a Controller OwnerReference on controlled.
// This is used for garbage collection of the controlled object and for
// reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).
// Since only one OwnerReference can be a controller, it returns an error if
// there is another OwnerReference with Controller flag set.
//
// This customized version was taken from https://github.com/kubernetes-sigs/controller-runtime/blob/196828e54e4210497438671b2b449522c004db5c/pkg/controller/controllerutil/controllerutil.go#L96
func SetControllerReference(owner, controlled metav1.Object, gvk *schema.GroupVersionKind) error {
	// Validate the owner.
	_, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime.Object, cannot call SetControllerReference", owner)
	}
	if err := validateOwner(owner, controlled); err != nil {
		return err
	}

	// Create a new controller ref.
	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
	}

	// Return early with an error if the object is already controlled.
	if existing := metav1.GetControllerOf(controlled); existing != nil && !referSameObject(*existing, ref) {
		return &controllerutil.AlreadyOwnedError{
			Object: controlled,
			Owner:  *existing,
		}
	}

	// Update owner references and return.
	upsertOwnerRef(ref, controlled)
	return nil
}

func upsertOwnerRef(ref metav1.OwnerReference, object metav1.Object) {
	owners := object.GetOwnerReferences()
	if idx := indexOwnerRef(owners, ref); idx == -1 {
		owners = append(owners, ref)
	} else {
		owners[idx] = ref
	}
	object.SetOwnerReferences(owners)
}

// indexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func indexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

func validateOwner(owner, object metav1.Object) error {
	objectName := object.GetName()
	ownerNs := owner.GetNamespace()
	if ownerNs != "" {
		objNs := object.GetNamespace()
		if objNs == "" {
			return fmt.Errorf("cluster-scoped resource %s must not have a namespace-scoped owner, owner's namespace %s", objectName, ownerNs)
		}
		if ownerNs != objNs {
			return fmt.Errorf("cross-namespace owner references are disallowed, owner's namespace %s, obj's name %s, obj's namespace %s", owner.GetNamespace(), objectName, object.GetNamespace())
		}
	}
	return nil
}

// Returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}
