package patch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type resourcePatcher struct {
	applier jsonPatchApplier
}

// NewResourcePatcher creates a new resource patcher.
func NewResourcePatcher(applier jsonPatchApplier) *resourcePatcher {
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

// Patch applies a configured patch in JSON patch format to a Kubernetes resource.
func (r *resourcePatcher) Patch(ctx context.Context, phase Phase, patches []ResourcePatch) error {
	var errs []error
	for _, patch := range filterPatchesByPhase(phase, patches) {
		err := r.patchSingle(ctx, patch)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *resourcePatcher) patchSingle(ctx context.Context, patch ResourcePatch) error {
	jsonPatch, err := json.Marshal(patch.Patches)
	if err != nil {
		return fmt.Errorf("failed to marshal json patch: %w", err)
	}

	return r.applier.Patch(ctx, jsonPatch, patch.Resource.GroupVersionKind(), patch.Resource.Name)
}
