package tfstate_test

import (
	"testing"

	"github.com/your-org/driftcheck/internal/tfstate"
)

func buildTestState(t *testing.T) *tfstate.State {
	t.Helper()
	s := tfstate.NewState()
	s.Add(tfstate.Resource{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-001"}})
	s.Add(tfstate.Resource{Type: "aws_instance", Name: "worker", Attributes: map[string]interface{}{"id": "i-002"}})
	s.Add(tfstate.Resource{Type: "aws_s3_bucket", Name: "assets", Attributes: map[string]interface{}{"id": "my-bucket"}})
	s.Add(tfstate.Resource{Type: "aws_security_group", Name: "web_sg", Attributes: map[string]interface{}{"id": "sg-001"}})
	return s
}

func TestFilter_NoOptions(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{})
	if len(got.Keys()) != 4 {
		t.Fatalf("expected 4 resources, got %d", len(got.Keys()))
	}
}

func TestFilter_ByType(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{ResourceTypes: []string{"aws_instance"}})
	if len(got.Keys()) != 2 {
		t.Fatalf("expected 2 aws_instance resources, got %d", len(got.Keys()))
	}
}

func TestFilter_ExcludeType(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{ExcludeTypes: []string{"aws_s3_bucket", "aws_security_group"}})
	if len(got.Keys()) != 2 {
		t.Fatalf("expected 2 resources after exclusion, got %d", len(got.Keys()))
	}
}

func TestFilter_ByNamePrefix(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{NamePrefix: "web"})
	if len(got.Keys()) != 2 {
		t.Fatalf("expected 2 resources with prefix 'web', got %d", len(got.Keys()))
	}
}

func TestFilter_TypeAndPrefix(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{
		ResourceTypes: []string{"aws_instance"},
		NamePrefix:    "web",
	})
	if len(got.Keys()) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got.Keys()))
	}
}

func TestFilter_ExcludeTakesPrecedence(t *testing.T) {
	s := buildTestState(t)
	got := s.Filter(tfstate.FilterOptions{
		ResourceTypes: []string{"aws_instance"},
		ExcludeTypes:  []string{"aws_instance"},
	})
	if len(got.Keys()) != 0 {
		t.Fatalf("expected 0 resources when type is both included and excluded, got %d", len(got.Keys()))
	}
}
