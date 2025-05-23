package patch

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

// ResourcePatch contains json patches for kubernetes resources to be applied on a phase of the setup process.
// The patch is applied at the end of its Phase.
// For namespaced resources, the namespace of the setup is inferred.
type ResourcePatch struct {
	// Phase is a sequential step in the setup process.
	Phase Phase `json:"phase" yaml:"phase"`
	// ResourceReference uniquely identifies a kubernetes resource that should be patched.
	Resource ResourceReference `json:"resource" yaml:"resource"`
	// Patches contains a series of operations to be applied on the specified kubernetes resource.
	Patches []JsonPatch `json:"patches" yaml:"patches"`
}

func (rp *ResourcePatch) Validate() error {
	var errs []error

	if !existsPhase(rp.Phase) {
		errs = append(errs, fmt.Errorf("phase '%s' does not exist", rp.Phase))
	}

	if rp.Resource.Kind == "" {
		errs = append(errs, fmt.Errorf("resource kind must not be empty"))
	}

	if rp.Resource.Name == "" {
		errs = append(errs, fmt.Errorf("resource name must not be empty"))
	}

	if len(rp.Patches) == 0 {
		errs = append(errs, fmt.Errorf("no patches found which is a sign of a misconfiguration"))
	}

	for _, singlePatch := range rp.Patches {
		err := singlePatch.Validate()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func existsPhase(phase Phase) bool {
	switch phase {
	case ComponentPhase:
		fallthrough
	case DoguPhase:
		fallthrough
	case LoadbalancerPhase:
		return true
	default:
		return false
	}
}

// Phase is a sequential step in the setup process.
type Phase string

const (
	// ComponentPhase is the step where components will be installed.
	ComponentPhase Phase = "component"
	// DoguPhase is the step where dogus will be installed.
	DoguPhase Phase = "dogu"
	// LoadbalancerPhase is the step where the external loadbalancer for the Cloudogu EcoSystem will be created.
	LoadbalancerPhase Phase = "loadbalancer"
)

// ResourceReference uniquely identifies a kubernetes resource.
type ResourceReference struct {
	// ApiVersion contains the group and version of the resource.
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	// Kind contains the resource type.
	Kind string `json:"kind" yaml:"kind"`
	// Name contains the name of the resource.
	Name string `json:"name" yaml:"name"`
}

func (r ResourceReference) GroupVersionKind() schema.GroupVersionKind {
	parts := strings.Split(r.ApiVersion, "/")
	if len(parts) > 1 {
		return schema.GroupVersionKind{
			Group:   parts[0],
			Version: parts[1],
			Kind:    r.Kind,
		}
	}

	return schema.GroupVersionKind{
		Group:   "",
		Version: r.ApiVersion,
		Kind:    r.Kind,
	}
}

// JsonPatchOperation describes how a json object should be modified.
type JsonPatchOperation string

const (
	addOperation     JsonPatchOperation = "add"
	removeOperation  JsonPatchOperation = "remove"
	replaceOperation JsonPatchOperation = "replace"
)

// JsonPatch describes a single operation on a kubernetes resource.
type JsonPatch struct {
	// Operation describes how a json object should be modified.
	Operation JsonPatchOperation `yaml:"op" json:"op"`
	// Path contains a JSON pointer to the value that should be modified.
	// If keys contain '/' or '~', those characters have to be replaced with '~1' and '~0' respectively.
	Path string `yaml:"path" json:"path"`
	// Value contains the value that should be inserted at the specified path.
	Value any `yaml:"value,omitempty" json:"value,omitempty"`
}

// Validate checks the JsonPatch for errors.
func (j JsonPatch) Validate() error {
	if j.Value == nil && (j.Operation == addOperation || j.Operation == replaceOperation) {
		return fmt.Errorf("value must not be empty for operation '%s' on path '%s'", j.Operation, j.Path)
	}

	if j.Value != nil && j.Operation == removeOperation {
		return fmt.Errorf("the operation '%s' on path '%s' does not take a value but it was provided: '%s'", j.Operation, j.Path, j.Value)
	}

	return nil
}
