package validation

import (
	"fmt"
	"github.com/cloudogu/k8s-ces-setup/app/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewResourcePatchConfigurationValidator(t *testing.T) {
	t.Run("should return a valid object", func(t *testing.T) {
		actual := NewResourcePatchConfigurationValidator()
		require.NotNil(t, actual)
	})
}

func Test_resourcePatchValidator_Validate(t *testing.T) {
	validPatches := []patch.JsonPatch{{Operation: "add", Path: "/metadata/annotations", Value: map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}}
	validResourceConfigs := []patch.ResourcePatch{
		{Phase: patch.DoguPhase, Resource: patch.ResourceReference{ApiVersion: "v1", Kind: "Pod", Name: "my-pod"}, Patches: validPatches},
	}
	invalidResourceConfigs := []patch.ResourcePatch{
		{Phase: "boohoo", Resource: patch.ResourceReference{ApiVersion: "v1", Kind: "Pod", Name: "my-pod"}, Patches: validPatches},
	}

	type args struct {
		resourcePatchConfig []patch.ResourcePatch
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"validates", args{validResourceConfigs}, assert.NoError},
		{"invalidates", args{invalidResourceConfigs}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resourcePatchValidator{}
			tt.wantErr(t, r.Validate(tt.args.resourcePatchConfig), fmt.Sprintf("Validate(%v)", tt.args.resourcePatchConfig))
		})
	}
}
