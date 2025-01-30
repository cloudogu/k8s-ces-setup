package patch

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceReference_GroupVersionKind(t *testing.T) {
	t.Run("should return GVK for a core API", func(t *testing.T) {
		// given

		sut := ResourceReference{
			ApiVersion: "v1",
			Kind:       "Pod",
			Name:       "my-pod",
		}

		// when
		actual := sut.GroupVersionKind()

		// then
		assert.Equal(t, schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}, actual)
	})

	t.Run("should return GVK for a grouped API", func(t *testing.T) {
		// given

		sut := ResourceReference{
			ApiVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       "my-deploy",
		}

		// when
		actual := sut.GroupVersionKind()

		// then
		assert.Equal(t, schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    "Deployment",
		}, actual)
	})
}

func TestJsonPatch_unmarshalling(t *testing.T) {
	tests := []struct {
		jsonRepresentation   []byte
		objectRepresentation JsonPatch
	}{
		{
			jsonRepresentation:   []byte(`{}`),
			objectRepresentation: JsonPatch{},
		},
		{
			jsonRepresentation: []byte(`{"op": "add", "path": "/metadata/labels", "value": {"foo": "bar"}}`),
			objectRepresentation: JsonPatch{
				Operation: "add",
				Path:      "/metadata/labels",
				Value:     map[string]interface{}{"foo": "bar"},
			},
		},
		{
			jsonRepresentation: []byte(`{"op": "add", "path": "/spec/loadBalancerIP", "value": "203.0.113.55"}`),
			objectRepresentation: JsonPatch{
				Operation: "add",
				Path:      "/spec/loadBalancerIP",
				Value:     "203.0.113.55",
			},
		},
		{
			jsonRepresentation: []byte(`{"op": "remove", "path": "/spec/loadBalancerIP"}`),
			objectRepresentation: JsonPatch{
				Operation: "remove",
				Path:      "/spec/loadBalancerIP",
			},
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.jsonRepresentation), func(t *testing.T) {
			var actual JsonPatch
			err := json.Unmarshal(tt.jsonRepresentation, &actual)
			require.NoError(t, err)
			assert.Equal(t, tt.objectRepresentation, actual)
		})
	}
}

func TestJsonPatch_marshalling(t *testing.T) {
	tests := []struct {
		jsonRepresentation   []byte
		objectRepresentation JsonPatch
	}{
		{
			jsonRepresentation:   []byte(`{"op":"","path":""}`),
			objectRepresentation: JsonPatch{},
		},
		{
			jsonRepresentation: []byte(`{"op":"add","path":"/metadata/labels","value":{"foo":"bar"}}`),
			objectRepresentation: JsonPatch{
				Operation: "add",
				Path:      "/metadata/labels",
				Value:     map[string]interface{}{"foo": "bar"},
			},
		},
		{
			jsonRepresentation: []byte(`{"op":"add","path":"/spec/loadBalancerIP","value":"203.0.113.55"}`),
			objectRepresentation: JsonPatch{
				Operation: "add",
				Path:      "/spec/loadBalancerIP",
				Value:     "203.0.113.55",
			},
		},
		{
			jsonRepresentation: []byte(`{"op":"remove","path":"/spec/loadBalancerIP"}`),
			objectRepresentation: JsonPatch{
				Operation: "remove",
				Path:      "/spec/loadBalancerIP",
			},
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.jsonRepresentation), func(t *testing.T) {
			actual, err := json.Marshal(tt.objectRepresentation)
			require.NoError(t, err)
			assert.Equal(t, tt.jsonRepresentation, actual)
		})
	}
}

func TestJsonPatch_Validate(t *testing.T) {
	type fields struct {
		Operation JsonPatchOperation
		Path      string
		Value     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"add okay", fields{addOperation, "/metadata/annotations", map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}, assert.NoError},
		{"add misses value", fields{addOperation, "/metadata/annotations", nil}, assert.Error},
		{"rm okay", fields{removeOperation, "/metadata/annotations", nil}, assert.NoError},
		{"rm unexpected value", fields{removeOperation, "/metadata/annotations", map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}, assert.Error},
		{"repl okay", fields{replaceOperation, "/metadata/annotations", map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}, assert.NoError},
		{"repl misses value", fields{replaceOperation, "/metadata/annotations", nil}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JsonPatch{
				Operation: tt.fields.Operation,
				Path:      tt.fields.Path,
				Value:     tt.fields.Value,
			}
			tt.wantErr(t, j.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}

func Test_existsPhase(t *testing.T) {
	assert.True(t, existsPhase(DoguPhase))
	assert.True(t, existsPhase(ComponentPhase))
	assert.True(t, existsPhase(LoadbalancerPhase))
	assert.False(t, existsPhase("notexisting"))
}

func TestResourcePatch_Validate(t *testing.T) {
	validPatches := []JsonPatch{{Operation: addOperation, Path: "/metadata/annotations", Value: map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}}
	invalidPatches := []JsonPatch{{Operation: addOperation, Path: "/metadata/annotations"}}

	type fields struct {
		Phase    Phase
		Resource ResourceReference
		Patches  []JsonPatch
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"validates", fields{DoguPhase, ResourceReference{"v1", "Pod", "my-pod"}, validPatches}, assert.NoError},
		{"invalid phase", fields{"typohere", ResourceReference{"v1", "Pod", "my-pod"}, validPatches}, assert.Error},
		{"invalid patch", fields{DoguPhase, ResourceReference{"v1", "Pod", "my-pod"}, invalidPatches}, assert.Error},
		{"empty resource reference name", fields{DoguPhase, ResourceReference{"v1", "Pod", ""}, validPatches}, assert.Error},
		{"empty kind", fields{DoguPhase, ResourceReference{"v1", "", "ignore"}, validPatches}, assert.Error},
		{"nil patch slice", fields{DoguPhase, ResourceReference{"v1", "Pod", "ignore"}, nil}, assert.Error},
		{"empty patch slice", fields{DoguPhase, ResourceReference{"v1", "Pod", "ignore"}, []JsonPatch{}}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResourcePatch{
				Phase:    tt.fields.Phase,
				Resource: tt.fields.Resource,
				Patches:  tt.fields.Patches,
			}
			tt.wantErr(t, rp.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}
