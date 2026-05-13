package tfstate

import (
	"fmt"
	"testing"
)

func TestResourceKey_String(t *testing.T) {
	k := ResourceKey{Type: "aws_instance", Name: "web"}
	want := "aws_instance.web"
	if got := k.String(); got != want {
		t.Errorf("ResourceKey.String() = %q; want %q", got, want)
	}
}

func TestNewState_Empty(t *testing.T) {
	s := NewState()
	if s.Len() != 0 {
		t.Errorf("expected empty state, got %d resources", s.Len())
	}
}

func TestState_AddAndGet(t *testing.T) {
	s := NewState()
	r := &Resource{
		Key:        ResourceKey{Type: "aws_s3_bucket", Name: "assets"},
		Provider:   "aws",
		Attributes: map[string]interface{}{"bucket": "my-assets"},
	}
	s.Add(r)

	if s.Len() != 1 {
		t.Fatalf("expected 1 resource, got %d", s.Len())
	}

	got := s.Get(r.Key)
	if got == nil {
		t.Fatal("expected to find resource, got nil")
	}
	if got.Attributes["bucket"] != "my-assets" {
		t.Errorf("unexpected attribute value: %v", got.Attributes["bucket"])
	}
}

func TestState_GetMissing(t *testing.T) {
	s := NewState()
	k := ResourceKey{Type: "aws_instance", Name: "missing"}
	if got := s.Get(k); got != nil {
		t.Errorf("expected nil for missing key, got %v", got)
	}
}

func TestState_Keys(t *testing.T) {
	s := NewState()
	for i := 0; i < 3; i++ {
		s.Add(&Resource{
			Key:        ResourceKey{Type: "aws_instance", Name: fmt.Sprintf("web%d", i)},
			Attributes: map[string]interface{}{},
		})
	}

	keys := s.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestState_AddOverwrite(t *testing.T) {
	s := NewState()
	k := ResourceKey{Type: "aws_instance", Name: "web"}

	s.Add(&Resource{Key: k, Attributes: map[string]interface{}{"ami": "ami-old"}})
	s.Add(&Resource{Key: k, Attributes: map[string]interface{}{"ami": "ami-new"}})

	if s.Len() != 1 {
		t.Fatalf("expected 1 resource after overwrite, got %d", s.Len())
	}
	if s.Get(k).Attributes["ami"] != "ami-new" {
		t.Errorf("expected overwritten value ami-new")
	}
}

func TestState_Remove(t *testing.T) {
	s := NewState()
	k := ResourceKey{Type: "aws_instance", Name: "web"}
	s.Add(&Resource{Key: k, Attributes: map[string]interface{}{"ami": "ami-123"}})

	if s.Len() != 1 {
		t.Fatalf("expected 1 resource before remove, got %d", s.Len())
	}

	s.Remove(k)

	if s.Len() != 0 {
		t.Errorf("expected 0 resources after remove, got %d", s.Len())
	}
	if s.Get(k) != nil {
		t.Errorf("expected nil after remove, got resource")
	}
}
