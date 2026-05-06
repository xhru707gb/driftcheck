package cloud

import (
	"context"
	"fmt"

	"github.com/example/driftcheck/internal/tfstate"
)

// DriftItem describes a single attribute difference between Terraform state and live cloud state.
type DriftItem struct {
	ResourceType string
	ResourceID   string
	Attribute    string
	Expected     interface{}
	Actual       interface{}
}

// ReconcileResult holds all drift items found during reconciliation.
type ReconcileResult struct {
	Drifts  []DriftItem
	Missing []string // resource IDs present in state but not found in cloud
}

// Reconciler compares tfstate resources against live cloud state.
type Reconciler struct {
	fetcher Fetcher
}

// NewReconciler creates a Reconciler with the given Fetcher.
func NewReconciler(f Fetcher) *Reconciler {
	return &Reconciler{fetcher: f}
}

// Reconcile iterates over all resources in the state and checks for drift.
func (r *Reconciler) Reconcile(ctx context.Context, state *tfstate.State) (*ReconcileResult, error) {
	result := &ReconcileResult{}

	for _, res := range state.Resources {
		for _, inst := range res.Instances {
			id, ok := inst.Attributes["id"]
			if !ok {
				continue
			}
			resourceID := fmt.Sprintf("%v", id)

			live, err := r.fetcher.Fetch(ctx, res.Type, resourceID)
			if err != nil {
				result.Missing = append(result.Missing, resourceID)
				continue
			}

			for key, expected := range inst.Attributes {
				if key == "id" {
					continue
				}
				actual, exists := live.Attributes[key]
				if !exists {
					continue
				}
				if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", actual) {
					result.Drifts = append(result.Drifts, DriftItem{
						ResourceType: res.Type,
						ResourceID:   resourceID,
						Attribute:    key,
						Expected:     expected,
						Actual:       actual,
					})
				}
			}
		}
	}
	return result, nil
}
