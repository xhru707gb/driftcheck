package tfstate

import "fmt"

// PlanAction represents the type of change proposed.
type PlanAction string

const (
	PlanCreate  PlanAction = "create"
	PlanUpdate  PlanAction = "update"
	PlanDestroy PlanAction = "destroy"
	PlanNoOp    PlanAction = "no-op"
)

// PlanEntry describes a proposed change for a single resource.
type PlanEntry struct {
	Key        ResourceKey
	Action     PlanAction
	ChangedKeys []string
}

// Plan holds the full set of proposed changes between two states.
type Plan struct {
	Entries []PlanEntry
}

// HasChanges returns true if any entry requires an action.
func (p *Plan) HasChanges() bool {
	for _, e := range p.Entries {
		if e.Action != PlanNoOp {
			return true
		}
	}
	return false
}

// Summary returns a human-readable one-liner for the plan.
func (p *Plan) Summary() string {
	var create, update, destroy int
	for _, e := range p.Entries {
		switch e.Action {
		case PlanCreate:
			create++
		case PlanUpdate:
			update++
		case PlanDestroy:
			destroy++
		}
	}
	return fmt.Sprintf("Plan: %d to create, %d to update, %d to destroy.", create, update, destroy)
}

// BuildPlan compares desired (Terraform) state against current (live) state
// and produces a Plan describing what changes would be applied.
func BuildPlan(desired, current *State) (*Plan, error) {
	if desired == nil {
		return nil, fmt.Errorf("desired state must not be nil")
	}
	if current == nil {
		return nil, fmt.Errorf("current state must not be nil")
	}

	plan := &Plan{}

	for _, key := range desired.Keys() {
		desiredRes, _ := desired.Get(key)
		currentRes, ok := current.Get(key)
		if !ok {
			plan.Entries = append(plan.Entries, PlanEntry{Key: key, Action: PlanCreate})
			continue
		}
		changed := diffAttrKeys(desiredRes.Attributes, currentRes.Attributes)
		if len(changed) > 0 {
			plan.Entries = append(plan.Entries, PlanEntry{Key: key, Action: PlanUpdate, ChangedKeys: changed})
		} else {
			plan.Entries = append(plan.Entries, PlanEntry{Key: key, Action: PlanNoOp})
		}
	}

	for _, key := range current.Keys() {
		if _, ok := desired.Get(key); !ok {
			plan.Entries = append(plan.Entries, PlanEntry{Key: key, Action: PlanDestroy})
		}
	}

	return plan, nil
}

func diffAttrKeys(desired, current map[string]interface{}) []string {
	var changed []string
	for k, dv := range desired {
		if cv, ok := current[k]; !ok || fmt.Sprintf("%v", cv) != fmt.Sprintf("%v", dv) {
			changed = append(changed, k)
		}
	}
	return changed
}
