package patch

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
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
		{"add okay", fields{addOperation, "/spec/thing", `a.value: 15`}, assert.NoError},
		{"add misses value", fields{addOperation, "/spec/thing", nil}, assert.Error},
		{"rm okay", fields{removeOperation, "/spec/thing", nil}, assert.NoError},
		{"rm unexpected value", fields{removeOperation, "/spec/thing", `a.value: 15`}, assert.Error},
		{"repl okay", fields{replaceOperation, "/spec/thing", `a.value: 15`}, assert.NoError},
		{"repl misses value", fields{replaceOperation, "/spec/thing", nil}, assert.Error},
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
	validPatches := []JsonPatch{{Operation: addOperation, Path: "/spec/thething", Value: "Annotation: 0"}}
	invalidPatches := []JsonPatch{{Operation: addOperation, Path: "/spec/thething"}}

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
		{"resource reference name is empty", fields{DoguPhase, ResourceReference{"v1", "Pod", ""}, validPatches}, assert.Error},
		{"kind is empty", fields{DoguPhase, ResourceReference{"v1", "", "ignore"}, validPatches}, assert.Error},
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
