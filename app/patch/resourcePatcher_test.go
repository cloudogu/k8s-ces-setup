package patch

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gkvLoadbalancer = ResourceReference{
	ApiVersion: "v1",
	Kind:       "Service",
	Name:       "ces-loadbalancer",
}

var testCtx = context.Background()

func Test_resourcePatcher_Patch(t *testing.T) {
	t.Run("should patch one resource in loadbalancer phase and ignore other phases", func(t *testing.T) {
		validPatches := []JsonPatch{{Operation: addOperation, Path: "/spec/thething", Value: map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}}
		// given
		patches := []ResourcePatch{{
			Phase:    LoadbalancerPhase,
			Resource: gkvLoadbalancer,
			Patches:  validPatches,
		}, {
			Phase: DoguPhase,
			Resource: ResourceReference{
				ApiVersion: "must not be processed",
				Kind:       "must not be processed",
				Name:       "must not be processed",
			},
			Patches: validPatches, // valid to make sure that it will be ignored
		}}
		mockApplier := newMockJsonPatchApplier(t)
		patchesBytes := marshalJson(t, validPatches)
		mockApplier.EXPECT().Patch(testCtx, patchesBytes, gkvLoadbalancer.GroupVersionKind(), gkvLoadbalancer.Name).Return(nil)

		sut := resourcePatcher{applier: mockApplier}

		// when
		err := sut.Patch(testCtx, LoadbalancerPhase, patches)

		// then
		require.NoError(t, err)
	})
	t.Run("should fail because patch failed to apply", func(t *testing.T) {
		validPatches := []JsonPatch{{Operation: addOperation, Path: "/spec/thething", Value: map[string]interface{}{"service.beta.kubernetes.io/azure-load-balancer-internal": "true"}}}
		invalidJsonMap := map[bool]string{true: "😅"}
		invalidPatches := []JsonPatch{{Operation: addOperation, Path: "/spec/thething", Value: map[string]interface{}{"invalid": invalidJsonMap}}}
		// given
		patches := []ResourcePatch{{
			Phase:    LoadbalancerPhase,
			Resource: gkvLoadbalancer,
			Patches:  invalidPatches,
		}, {
			Phase:    DoguPhase,
			Resource: ResourceReference{},
			Patches:  validPatches, // valid to make sure that it will be ignored
		}}
		mockApplier := newMockJsonPatchApplier(t)

		sut := resourcePatcher{applier: mockApplier}

		// when
		err := sut.Patch(testCtx, LoadbalancerPhase, patches)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to marshal json patch: json: unsupported type: map[bool]string")
	})
}

func marshalJson(t *testing.T, patches []JsonPatch) []byte {
	t.Helper()

	bytes, err := json.Marshal(patches)
	require.NoError(t, err)

	return bytes
}

func TestNewResourcePatcher(t *testing.T) {
	assert.NotNil(t, NewResourcePatcher(nil))
}
