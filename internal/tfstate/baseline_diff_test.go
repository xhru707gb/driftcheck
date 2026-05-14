package tfstate_test

import (
	"testing"

	"github.com/example/driftcheck/internal/tfstate"
)

func buildDiffBaseline(t *testing.T) *tfstate.Baseline {
	t.Helper()
	s := buildBaselineState(t)
	b, err := tfstate.NewBaseline(s)
	if err != nil {
		t.Fatalf("NewBaseline: %v", err)
	}
	return b
}

func TestCompareToBaseline_NoDiff(t *testing.T) {
	b := buildDiffBaseline(t)
	current := buildBaselineState(t)
	diffs, err := tfstate.CompareToBaseline(b, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d: %v", len(diffs), diffs)
	}
}

func TestCompareToBaseline_ModifiedAttribute(t *testing.T) {
	b := buildDiffBaseline(t)
	current := buildBaselineState(t)
	r, _ := current.Get(tfstate.ResourceKey{Type: "aws_instance", Name: "web"})
	r.Attributes["instance_type"] = "t3.large"
	current.Add(r)

	diffs, err := tfstate.CompareToBaseline(b, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var found bool
	for _, d := range diffs {
		if d.Kind == tfstate.BaselineModified && d.Attribute == "instance_type" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected modified diff for instance_type, got: %v", diffs)
	}
}

func TestCompareToBaseline_AddedResource(t *testing.T) {
	b := buildDiffBaseline(t)
	current := buildBaselineState(t)
	current.Add(tfstate.Resource{
		Type: "aws_vpc", Name: "main",
		Attributes: map[string]interface{}{"cidr_block": "10.0.0.0/16"},
	})

	diffs, _ := tfstate.CompareToBaseline(b, current)
	for _, d := range diffs {
		if d.Kind == tfstate.BaselineAdded && d.ResourceKey == "aws_vpc.main" {
			return
		}
	}
	t.Errorf("expected added diff for aws_vpc.main, got: %v", diffs)
}

func TestCompareToBaseline_RemovedResource(t *testing.T) {
	b := buildDiffBaseline(t)
	current := tfstate.NewState()
	current.Add(tfstate.Resource{
		Type: "aws_instance", Name: "web",
		Attributes: map[string]interface{}{"instance_type": "t3.micro", "ami": "ami-123"},
	})
	// aws_s3_bucket.assets is absent

	diffs, _ := tfstate.CompareToBaseline(b, current)
	for _, d := range diffs {
		if d.Kind == tfstate.BaselineRemoved && d.ResourceKey == "aws_s3_bucket.assets" {
			return
		}
	}
	t.Errorf("expected removed diff for aws_s3_bucket.assets, got: %v", diffs)
}

func TestCompareToBaseline_NilInputs(t *testing.T) {
	s := buildBaselineState(t)
	b, _ := tfstate.NewBaseline(s)

	if _, err := tfstate.CompareToBaseline(nil, s); err == nil {
		t.Error("expected error for nil baseline")
	}
	if _, err := tfstate.CompareToBaseline(b, nil); err == nil {
		t.Error("expected error for nil state")
	}
}
