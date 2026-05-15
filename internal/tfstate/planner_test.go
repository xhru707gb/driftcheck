package tfstate

import (
	"testing"
)

func buildPlannerState(resources []Resource) *State {
	s := NewState()
	for _, r := range resources {
		s.Add(r)
	}
	return s
}

func TestBuildPlan_NilDesired(t *testing.T) {
	_, err := BuildPlan(nil, NewState())
	if err == nil {
		t.Fatal("expected error for nil desired state")
	}
}

func TestBuildPlan_NilCurrent(t *testing.T) {
	_, err := BuildPlan(NewState(), nil)
	if err == nil {
		t.Fatal("expected error for nil current state")
	}
}

func TestBuildPlan_NoChanges(t *testing.T) {
	res := Resource{Type: "aws_s3_bucket", Name: "main", ID: "b1", Attributes: map[string]interface{}{"region": "us-east-1"}}
	desired := buildPlannerState([]Resource{res})
	current := buildPlannerState([]Resource{res})
	plan, err := BuildPlan(desired, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestBuildPlan_Create(t *testing.T) {
	res := Resource{Type: "aws_instance", Name: "web", ID: "i-1", Attributes: map[string]interface{}{}}
	desired := buildPlannerState([]Resource{res})
	current := buildPlannerState([]Resource{})
	plan, err := BuildPlan(desired, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.Entries) != 1 || plan.Entries[0].Action != PlanCreate {
		t.Errorf("expected one create entry, got %+v", plan.Entries)
	}
}

func TestBuildPlan_Destroy(t *testing.T) {
	res := Resource{Type: "aws_instance", Name: "old", ID: "i-2", Attributes: map[string]interface{}{}}
	desired := buildPlannerState([]Resource{})
	current := buildPlannerState([]Resource{res})
	plan, err := BuildPlan(desired, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.Entries) != 1 || plan.Entries[0].Action != PlanDestroy {
		t.Errorf("expected one destroy entry, got %+v", plan.Entries)
	}
}

func TestBuildPlan_Update(t *testing.T) {
	desiredRes := Resource{Type: "aws_s3_bucket", Name: "logs", ID: "b2", Attributes: map[string]interface{}{"acl": "private"}}
	currentRes := Resource{Type: "aws_s3_bucket", Name: "logs", ID: "b2", Attributes: map[string]interface{}{"acl": "public"}}
	desired := buildPlannerState([]Resource{desiredRes})
	current := buildPlannerState([]Resource{currentRes})
	plan, err := BuildPlan(desired, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.Entries) != 1 || plan.Entries[0].Action != PlanUpdate {
		t.Errorf("expected one update entry, got %+v", plan.Entries)
	}
	if len(plan.Entries[0].ChangedKeys) == 0 {
		t.Error("expected changed keys to be populated")
	}
}

func TestPlan_Summary(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Action: PlanCreate},
			{Action: PlanUpdate},
			{Action: PlanDestroy},
			{Action: PlanNoOp},
		},
	}
	s := plan.Summary()
	expected := "Plan: 1 to create, 1 to update, 1 to destroy."
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}
