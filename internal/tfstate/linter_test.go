package tfstate

import (
	"testing"
)

func buildLintState() *State {
	s := NewState()
	s.Add(Resource{
		Type: "aws_instance",
		Name: "web",
		ID:   "i-abc123",
		Attributes: map[string]interface{}{
			"ami":           "ami-0abcdef",
			"instance_type": "t3.micro",
		},
	})
	return s
}

func TestLint_ValidState(t *testing.T) {
	s := buildLintState()
	result := Lint(s)
	if result.HasIssues() {
		for _, issue := range result.Issues {
			t.Errorf("unexpected issue: %s", issue)
		}
	}
}

func TestLint_NilState(t *testing.T) {
	result := Lint(nil)
	if !result.HasErrors() {
		t.Fatal("expected error for nil state")
	}
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestLint_EmptyID(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "logs",
		ID:         "",
		Attributes: map[string]interface{}{"bucket": "my-logs"},
	})
	result := Lint(s)
	if !result.HasErrors() {
		t.Fatal("expected error for empty ID")
	}
}

func TestLint_NoAttributes(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "empty",
		ID:         "i-000",
		Attributes: map[string]interface{}{},
	})
	result := Lint(s)
	if result.HasErrors() {
		t.Fatal("expected no errors, only warnings")
	}
	if !result.HasIssues() {
		t.Fatal("expected a warning for no attributes")
	}
	if result.Issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", result.Issues[0].Severity)
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Severity: "error", Resource: "aws_instance.web", Message: "empty ID"}
	s := issue.String()
	if s != "[ERROR] aws_instance.web: empty ID" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	s := NewState()
	s.Add(Resource{Type: "", Name: "bad", ID: "", Attributes: map[string]interface{}{}})
	result := Lint(s)
	if len(result.Issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d", len(result.Issues))
	}
}
