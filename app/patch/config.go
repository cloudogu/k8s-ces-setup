package patch

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

// ResourcePatch contains json patches for kubernetes resources to be applied on a phase of the setup process.
// The patch is applied at the end of its Phase.
// For namespaced resources, the namespace of the setup is inferred.
type ResourcePatch struct {
	// Phase is a sequential step in the setup process.
	Phase Phase `yaml:"phase"`
	// ResourceReference uniquely identifies a kubernetes resource that should be patched.
	Resource ResourceReference `yaml:"resource"`
	// Patches contains a series of operations to be applied on the specified kubernetes resource.
	Patches []JsonPatch `yaml:"patches"`
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
	ApiVersion string `yaml:"apiVersion"`
	// Kind contains the resource type.
	Kind string `yaml:"kind"`
	// Name contains the name of the resource.
	Name string `yaml:"name"`
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
	Value interface{} `yaml:"value,omitempty" json:"value,omitempty"`
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
