package tfstate

import (
	"strings"
	"testing"
)

func buildValidState(t *testing.T) *State {
	t.Helper()
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		Attributes: map[string]interface{}{"id": "i-123", "ami": "ami-abc"},
	})
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "data",
		Attributes: map[string]interface{}{"id": "my-bucket", "region": "us-east-1"},
	})
	return s
}

func TestValidate_NilState(t *testing.T) {
	if err := Validate(nil); err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestValidate_ValidState(t *testing.T) {
	s := buildValidState(t)
	if err := Validate(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyState(t *testing.T) {
	s := NewState()
	if err := Validate(s); err != nil {
		t.Fatalf("empty state should be valid, got: %v", err)
	}
}

func TestValidate_MissingID(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		Attributes: map[string]interface{}{"ami": "ami-abc"},
	})
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for missing id")
	}
	if !strings.Contains(err.Error(), "missing required attribute \"id\"") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_EmptyAttributes(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		Attributes: map[string]interface{}{},
	})
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error for empty attributes")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) == 0 {
		t.Error("expected at least one issue")
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "",
		Name:       "",
		Attributes: map[string]interface{}{},
	})
	err := Validate(s)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve := err.(*ValidationError)
	if len(ve.Issues) < 3 {
		t.Errorf("expected at least 3 issues, got %d: %v", len(ve.Issues), ve.Issues)
	}
}
