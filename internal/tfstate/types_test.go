package tfstate

import (
	"testing"
)

func buildInventoryState() *State {
	s := NewState()
	s.Add(Resource{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-001"}})
	s.Add(Resource{Type: "aws_instance", Name: "api", Attributes: map[string]interface{}{"id": "i-002"}})
	s.Add(Resource{Type: "aws_s3_bucket", Name: "assets", Attributes: map[string]interface{}{"id": "bucket-1"}})
	s.Add(Resource{Type: "aws_security_group", Name: "default", Attributes: map[string]interface{}{"id": "sg-001"}})
	return s
}

func TestBuildTypeInventory_Total(t *testing.T) {
	inv := BuildTypeInventory(buildInventoryState())
	if got := inv.Total(); got != 4 {
		t.Errorf("Total() = %d, want 4", got)
	}
}

func TestBuildTypeInventory_Count(t *testing.T) {
	inv := BuildTypeInventory(buildInventoryState())
	if got := inv.Count("aws_instance"); got != 2 {
		t.Errorf("Count(aws_instance) = %d, want 2", got)
	}
	if got := inv.Count("aws_s3_bucket"); got != 1 {
		t.Errorf("Count(aws_s3_bucket) = %d, want 1", got)
	}
	if got := inv.Count("aws_lambda_function"); got != 0 {
		t.Errorf("Count(aws_lambda_function) = %d, want 0", got)
	}
}

func TestBuildTypeInventory_Types(t *testing.T) {
	inv := BuildTypeInventory(buildInventoryState())
	types := inv.Types()
	if len(types) != 3 {
		t.Fatalf("Types() len = %d, want 3", len(types))
	}
	// Verify sorted order
	expected := []string{"aws_instance", "aws_s3_bucket", "aws_security_group"}
	for i, rt := range expected {
		if types[i] != rt {
			t.Errorf("Types()[%d] = %q, want %q", i, types[i], rt)
		}
	}
}

func TestBuildTypeInventory_Has(t *testing.T) {
	inv := BuildTypeInventory(buildInventoryState())
	if !inv.Has("aws_instance") {
		t.Error("Has(aws_instance) = false, want true")
	}
	if inv.Has("aws_lambda_function") {
		t.Error("Has(aws_lambda_function) = true, want false")
	}
}

func TestBuildTypeInventory_EmptyState(t *testing.T) {
	inv := BuildTypeInventory(NewState())
	if inv.Total() != 0 {
		t.Errorf("Total() on empty state = %d, want 0", inv.Total())
	}
	if len(inv.Types()) != 0 {
		t.Errorf("Types() on empty state = %v, want []", inv.Types())
	}
}
