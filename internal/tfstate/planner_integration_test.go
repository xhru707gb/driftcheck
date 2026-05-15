package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

// TestBuildPlanAndReport exercises BuildPlan + WritePlanReport end-to-end.
func TestBuildPlanAndReport_FullScenario(t *testing.T) {
	desired := buildPlannerState([]Resource{
		{Type: "aws_instance", Name: "web", ID: "i-1", Attributes: map[string]interface{}{"ami": "ami-123"}},
		{Type: "aws_s3_bucket", Name: "logs", ID: "b-1", Attributes: map[string]interface{}{"acl": "private"}},
	})
	current := buildPlannerState([]Resource{
		{Type: "aws_s3_bucket", Name: "logs", ID: "b-1", Attributes: map[string]interface{}{"acl": "public-read"}},
		{Type: "aws_vpc", Name: "old", ID: "vpc-1", Attributes: map[string]interface{}{}},
	})

	plan, err := BuildPlan(desired, current)
	if err != nil {
		t.Fatalf("BuildPlan error: %v", err)
	}

	if !plan.HasChanges() {
		t.Fatal("expected plan to have changes")
	}

	actions := map[PlanAction]int{}
	for _, e := range plan.Entries {
		actions[e.Action]++
	}
	if actions[PlanCreate] != 1 {
		t.Errorf("expected 1 create, got %d", actions[PlanCreate])
	}
	if actions[PlanUpdate] != 1 {
		t.Errorf("expected 1 update, got %d", actions[PlanUpdate])
	}
	if actions[PlanDestroy] != 1 {
		t.Errorf("expected 1 destroy, got %d", actions[PlanDestroy])
	}

	var buf bytes.Buffer
	if err := WritePlanReport(&buf, plan); err != nil {
		t.Fatalf("WritePlanReport error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 to create") {
		t.Errorf("summary missing create count: %s", out)
	}
	if !strings.Contains(out, "1 to update") {
		t.Errorf("summary missing update count: %s", out)
	}
	if !strings.Contains(out, "1 to destroy") {
		t.Errorf("summary missing destroy count: %s", out)
	}
}
