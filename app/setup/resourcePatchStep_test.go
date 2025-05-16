package setup

import (
	"github.com/cloudogu/k8s-ces-setup/v4/app/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var testPatches = []patch.ResourcePatch{}

func Test_resourcePatchStep_PerformSetupStep(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		// given
		mockPatcher := newMockResourcePatcher(t)
		mockPatcher.EXPECT().Patch(testCtx, patch.DoguPhase, testPatches).Return(nil)
		sut := &resourcePatchStep{phase: patch.DoguPhase, patcher: mockPatcher, patches: testPatches}

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
	t.Run("should return an error", func(t *testing.T) {
		// given
		mockPatcher := newMockResourcePatcher(t)
		mockPatcher.EXPECT().Patch(testCtx, patch.DoguPhase, testPatches).Return(assert.AnError)
		sut := &resourcePatchStep{phase: patch.DoguPhase, patcher: mockPatcher, patches: testPatches}

		// when
		err := sut.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
