package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

func TestWritePlanReport_NilPlan(t *testing.T) {
	var buf bytes.Buffer
	err := WritePlanReport(&buf, nil)
	if err == nil {
		t.Fatal("expected error for nil plan")
	}
}

func TestWritePlanReport_NoChanges(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Key: ResourceKey{Type: "aws_s3_bucket", Name: "main"}, Action: PlanNoOp},
		},
	}
	var buf bytes.Buffer
	if err := WritePlanReport(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' in output, got: %s", buf.String())
	}
}

func TestWritePlanReport_WithCreate(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Key: ResourceKey{Type: "aws_instance", Name: "web"}, Action: PlanCreate},
		},
	}
	var buf bytes.Buffer
	if err := WritePlanReport(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ aws_instance.web") {
		t.Errorf("expected create symbol in output, got: %s", out)
	}
	if !strings.Contains(out, "1 to create") {
		t.Errorf("expected summary in output, got: %s", out)
	}
}

func TestWritePlanReport_WithUpdate(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Key: ResourceKey{Type: "aws_s3_bucket", Name: "logs"}, Action: PlanUpdate, ChangedKeys: []string{"acl", "tags"}},
		},
	}
	var buf bytes.Buffer
	if err := WritePlanReport(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ aws_s3_bucket.logs") {
		t.Errorf("expected update symbol in output, got: %s", out)
	}
	if !strings.Contains(out, "acl") {
		t.Errorf("expected changed key 'acl' in output, got: %s", out)
	}
}

func TestWritePlanReport_WithDestroy(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Key: ResourceKey{Type: "aws_vpc", Name: "legacy"}, Action: PlanDestroy},
		},
	}
	var buf bytes.Buffer
	if err := WritePlanReport(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- aws_vpc.legacy") {
		t.Errorf("expected destroy symbol in output, got: %s", out)
	}
}
