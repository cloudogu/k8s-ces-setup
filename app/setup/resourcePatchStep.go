package setup

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/v4/app/patch"
)

type resourcePatcher interface {
	// Patch applies a configured patch to a Kubernetes resource.
	Patch(ctx context.Context, phase patch.Phase, patches []patch.ResourcePatch) error
}

// NewResourcePatchStep creates a new setup step which patches arbitrary Kubernetes resources according the given setup phase.
func NewResourcePatchStep(phase patch.Phase, patcher resourcePatcher, patches []patch.ResourcePatch) *resourcePatchStep {
	return &resourcePatchStep{phase: phase, patcher: patcher, patches: patches}
}

type resourcePatchStep struct {
	phase   patch.Phase
	patcher resourcePatcher
	patches []patch.ResourcePatch
}

// GetStepDescription returns the textual description of the resource patch step.
func (r *resourcePatchStep) GetStepDescription() string {
	return fmt.Sprintf("Patching kubernetes resources in phase %s", r.phase)
}

// PerformSetupStep executes the resource patch setup step.
func (r *resourcePatchStep) PerformSetupStep(ctx context.Context) error {
	err := r.patcher.Patch(ctx, r.phase, r.patches)
	if err != nil {
		return fmt.Errorf("failed to patch resources in phase %s: %w", r.phase, err)
	}

	return nil
}
