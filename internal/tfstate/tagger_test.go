package tfstate

import (
	"testing"
)

func buildTaggerState() *State {
	s := NewState()
	s.Add(Resource{
		Type: "aws_instance",
		Name: "web",
		ID:   "i-001",
		Attributes: map[string]string{
			"tags.env":  "prod",
			"tags.team": "platform",
		},
	})
	s.Add(Resource{
		Type: "aws_s3_bucket",
		Name: "assets",
		ID:   "bucket-1",
		Attributes: map[string]string{
			"tags.env": "staging",
		},
	})
	return s
}

func TestEnforceTags_NilState(t *testing.T) {
	rules := []TagRule{{Key: "env"}}
	report := EnforceTags(nil, rules)
	if report.HasViolations() {
		t.Error("expected no violations for nil state")
	}
}

func TestEnforceTags_NoViolations(t *testing.T) {
	s := buildTaggerState()
	rules := []TagRule{{Key: "env"}}
	report := EnforceTags(s, rules)
	if report.HasViolations() {
		t.Errorf("expected no violations, got %d", len(report.Violations))
	}
}

func TestEnforceTags_MissingTag(t *testing.T) {
	s := buildTaggerState()
	// aws_s3_bucket/assets has no tags.team
	rules := []TagRule{{Key: "team"}}
	report := EnforceTags(s, rules)
	if !report.HasViolations() {
		t.Fatal("expected violations")
	}
	if len(report.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(report.Violations))
	}
	if report.Violations[0].Resource.Name != "assets" {
		t.Errorf("unexpected resource %s", report.Violations[0].Resource.Name)
	}
}

func TestEnforceTags_DisallowedValue(t *testing.T) {
	s := buildTaggerState()
	// aws_s3_bucket/assets has tags.env=staging which is not in allowed list
	rules := []TagRule{{Key: "env", Values: []string{"prod", "dev"}}}
	report := EnforceTags(s, rules)
	if !report.HasViolations() {
		t.Fatal("expected violations")
	}
	v := report.Violations[0]
	if v.Actual != "staging" {
		t.Errorf("expected actual=staging, got %q", v.Actual)
	}
}

func TestEnforceTags_MultipleRules(t *testing.T) {
	s := buildTaggerState()
	rules := []TagRule{
		{Key: "env", Values: []string{"prod"}},
		{Key: "team"},
	}
	report := EnforceTags(s, rules)
	// aws_s3_bucket/assets: env=staging (disallowed) + missing team = 2 violations
	if len(report.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(report.Violations))
	}
}

func TestTagViolation_String_Missing(t *testing.T) {
	v := TagViolation{
		Resource: ResourceKey{Type: "aws_instance", Name: "web"},
		Rule:     TagRule{Key: "owner"},
	}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
