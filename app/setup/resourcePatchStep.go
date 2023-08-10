package setup

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/patch"
)

type resourcePatcher interface {
	Patch(phase patch.Phase, patches []patch.ResourcePatch) error
}

func NewResourcePatchStep(phase patch.Phase, patcher resourcePatcher, patches []patch.ResourcePatch) *resourcePatchStep {
	return &resourcePatchStep{phase: phase, patcher: patcher, patches: patches}
}

type resourcePatchStep struct {
	phase   patch.Phase
	patcher resourcePatcher
	patches []patch.ResourcePatch
}

func (r *resourcePatchStep) GetStepDescription() string {
	return fmt.Sprintf("Patching kubernetes resources in phase %s", r.phase)
}

func (r *resourcePatchStep) PerformSetupStep() error {
	err := r.patcher.Patch(r.phase, r.patches)
	if err != nil {
		return fmt.Errorf("failed to patch resources in phase %s: %w", r.phase, err)
	}

	return nil
}
