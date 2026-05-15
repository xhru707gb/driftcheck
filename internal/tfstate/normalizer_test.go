package tfstate

import (
	"testing"
)

func buildNormalizerState() *State {
	s := NewState()
	s.Add(Resource{
		Type:       "AWS_Instance",
		Name:       "web",
		ID:         "i-001",
		Attributes: map[string]interface{}{"ami": "  ami-123  ", "instance_type": "t2.micro"},
	})
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "data",
		ID:         "bkt-001",
		Attributes: map[string]interface{}{"bucket": "my-bucket", "region": " us-east-1 "},
	})
	return s
}

func TestNormalize_NilState(t *testing.T) {
	_, err := Normalize(nil)
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestNormalize_TrimsWhitespace(t *testing.T) {
	s := buildNormalizerState()
	res, err := Normalize(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Normalized == 0 {
		t.Fatal("expected at least one normalization")
	}

	web, ok := s.Get(ResourceKey{Type: "aws_instance", Name: "web"})
	if !ok {
		// key was stored under original type before normalise; look up new key
		t.Fatal("resource 'web' not found after normalization")
	}
	if v, _ := web.Attributes["ami"].(string); v != "ami-123" {
		t.Errorf("expected trimmed ami, got %q", v)
	}
}

func TestNormalize_LowercasesType(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "AWS_Instance",
		Name:       "srv",
		ID:         "i-002",
		Attributes: map[string]interface{}{},
	})

	_, err := Normalize(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := s.Get(ResourceKey{Type: "aws_instance", Name: "srv"})
	if !ok {
		t.Error("expected resource stored under lower-cased type key")
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_vpc",
		Name:       "main",
		ID:         "vpc-001",
		Attributes: map[string]interface{}{"cidr": "10.0.0.0/16"},
	})

	res, err := Normalize(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Normalized != 0 {
		t.Errorf("expected 0 normalizations, got %d", res.Normalized)
	}
}

func TestNormalizeResult_String_NoChanges(t *testing.T) {
	r := &NormalizeResult{}
	if r.String() != "no normalization changes" {
		t.Errorf("unexpected string: %s", r.String())
	}
}

func TestNormalizeResult_String_WithChanges(t *testing.T) {
	r := &NormalizeResult{Normalized: 1, Changes: []string{"web (aws_instance)"}}
	if r.String() == "no normalization changes" {
		t.Error("expected non-empty summary")
	}
}
