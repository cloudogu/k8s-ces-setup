package patch

import (
	"encoding/json"
	"errors"
	"fmt"
)

type resourcePatcher struct {
	applier jsonPatchApplier
}

func NewResourcePatcher(applier applier) *resourcePatcher {
	return &resourcePatcher{applier: applier}
}

func filterPatchesByPhase(phase Phase, patches []ResourcePatch) []ResourcePatch {
	var filtered []ResourcePatch
	for _, patch := range patches {
		if patch.Phase == phase {
			filtered = append(filtered, patch)
		}
	}
	return filtered
}

func (r *resourcePatcher) Patch(phase Phase, patches []ResourcePatch) error {
	var errs []error
	for _, patch := range filterPatchesByPhase(phase, patches) {
		err := r.patchSingle(patch)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *resourcePatcher) patchSingle(patch ResourcePatch) error {
	jsonPatch, err := json.Marshal(patch.Patches)
	if err != nil {
		return fmt.Errorf("failed to marshal json patch: %w", err)
	}

	return r.applier.Patch(jsonPatch, patch.Resource.GroupVersionKind(), patch.Resource.Name)
}
