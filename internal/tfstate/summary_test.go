package tfstate

import (
	"strings"
	"testing"
)

func buildSummaryState(t *testing.T) *State {
	t.Helper()
	s := NewState()

	s.Add(Resource{
		Type: "aws_instance",
		Name: "web",
		Instances: []Instance{{Attributes: map[string]interface{}{"id": "i-1"}}},
	})
	s.Add(Resource{
		Type: "aws_instance",
		Name: "worker",
		Instances: []Instance{
			{Attributes: map[string]interface{}{"id": "i-2"}},
			{Attributes: map[string]interface{}{"id": "i-3"}},
		},
	})
	s.Add(Resource{
		Type: "aws_s3_bucket",
		Name: "assets",
		Instances: []Instance{{Attributes: map[string]interface{}{"id": "b-1"}}},
	})
	return s
}

func TestSummarize_Totals(t *testing.T) {
	sum := Summarize(buildSummaryState(t))

	if sum.TotalResources != 3 {
		t.Errorf("expected 3 resources, got %d", sum.TotalResources)
	}
	if sum.TotalInstances != 4 {
		t.Errorf("expected 4 instances, got %d", sum.TotalInstances)
	}
}

func TestSummarize_ByType(t *testing.T) {
	sum := Summarize(buildSummaryState(t))

	if len(sum.ByType) != 2 {
		t.Fatalf("expected 2 types, got %d", len(sum.ByType))
	}

	ec2 := sum.ByType["aws_instance"]
	if ec2 == nil {
		t.Fatal("missing aws_instance type summary")
	}
	if ec2.Count != 2 {
		t.Errorf("expected 2 aws_instance resources, got %d", ec2.Count)
	}
	if ec2.Instances != 3 {
		t.Errorf("expected 3 aws_instance instances, got %d", ec2.Instances)
	}

	s3 := sum.ByType["aws_s3_bucket"]
	if s3 == nil {
		t.Fatal("missing aws_s3_bucket type summary")
	}
	if s3.Count != 1 || s3.Instances != 1 {
		t.Errorf("unexpected s3 counts: %+v", s3)
	}
}

func TestSummarize_EmptyState(t *testing.T) {
	sum := Summarize(NewState())

	if sum.TotalResources != 0 || sum.TotalInstances != 0 {
		t.Errorf("expected zeros for empty state, got %+v", sum)
	}
	if len(sum.ByType) != 0 {
		t.Errorf("expected empty ByType map")
	}
}

func TestSummarize_String(t *testing.T) {
	sum := Summarize(buildSummaryState(t))
	out := sum.String()

	if !strings.Contains(out, "Resources: 3") {
		t.Errorf("expected 'Resources: 3' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected 'aws_instance' in output")
	}
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Errorf("expected 'aws_s3_bucket' in output")
	}
}
