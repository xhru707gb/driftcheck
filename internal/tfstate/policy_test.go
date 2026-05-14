package tfstate_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/tfstate"
)

func buildPolicyState() *tfstate.State {
	s := tfstate.NewState()
	s.Add(tfstate.Resource{
		Type: "aws_instance",
		Name: "web",
		ID:   "i-abc123",
		Attributes: map[string]interface{}{
			"instance_type": "t3.micro",
			"tags":          map[string]interface{}{"env": "prod"},
		},
	})
	s.Add(tfstate.Resource{
		Type: "aws_s3_bucket",
		Name: "data",
		ID:   "my-bucket",
		Attributes: map[string]interface{}{
			"versioning": false,
		},
	})
	return s
}

func TestEnforcePolicy_NilState(t *testing.T) {
	_, err := tfstate.EnforcePolicy(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestEnforcePolicy_NoRules(t *testing.T) {
	s := buildPolicyState()
	report, err := tfstate.EnforcePolicy(s, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasViolations() {
		t.Fatal("expected no violations with no rules")
	}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	s := buildPolicyState()
	rules := []tfstate.PolicyRule{
		{
			Name:        "has-id",
			Description: "resource must have an ID",
			Check:       func(r tfstate.Resource) bool { return r.ID != "" },
		},
	}
	report, err := tfstate.EnforcePolicy(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasViolations() {
		t.Fatalf("expected no violations, got %d", len(report.Violations))
	}
	if report.Checked != 2 {
		t.Fatalf("expected 2 resources checked, got %d", report.Checked)
	}
}

func TestEnforcePolicy_WithViolation(t *testing.T) {
	s := buildPolicyState()
	rules := []tfstate.PolicyRule{
		{
			Name:        "versioning-enabled",
			Description: "S3 buckets must have versioning enabled",
			Check: func(r tfstate.Resource) bool {
				if r.Type != "aws_s3_bucket" {
					return true
				}
				v, ok := r.Attributes["versioning"]
				if !ok {
					return false
				}
				b, _ := v.(bool)
				return b
			},
		},
	}
	report, err := tfstate.EnforcePolicy(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasViolations() {
		t.Fatal("expected violations")
	}
	if len(report.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(report.Violations))
	}
	if report.Violations[0].RuleName != "versioning-enabled" {
		t.Errorf("unexpected rule name: %s", report.Violations[0].RuleName)
	}
}

func TestEnforcePolicy_MultipleViolations(t *testing.T) {
	s := buildPolicyState()
	rules := []tfstate.PolicyRule{
		{
			Name:        "always-fail",
			Description: "always fails",
			Check:       func(r tfstate.Resource) bool { return false },
		},
	}
	report, err := tfstate.EnforcePolicy(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(report.Violations))
	}
}
