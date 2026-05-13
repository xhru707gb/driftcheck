package tfstate_test

import (
	"testing"

	"github.com/your-org/driftcheck/internal/tfstate"
)

func buildDiffState(resources map[string]map[string]interface{}) *tfstate.State {
	s := tfstate.NewState()
	for name, attrs := range resources {
		s.Add(tfstate.Resource{
			Type:       "aws_instance",
			Name:       name,
			ID:         name + "-id",
			Attributes: attrs,
		})
	}
	return s
}

func TestDiffStates_NoDiff(t *testing.T) {
	attrs := map[string]interface{}{"ami": "ami-123", "instance_type": "t2.micro"}
	base := buildDiffState(map[string]map[string]interface{}{"web": attrs})
	target := buildDiffState(map[string]map[string]interface{}{"web": attrs})

	diff := tfstate.DiffStates(base, target)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got: %s", diff.Summary())
	}
}

func TestDiffStates_Added(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{})
	target := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-123"},
	})

	diff := tfstate.DiffStates(base, target)
	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(diff.Diffs) != 1 || diff.Diffs[0].Kind != tfstate.DiffAdded {
		t.Errorf("expected DiffAdded, got %+v", diff.Diffs)
	}
}

func TestDiffStates_Removed(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-123"},
	})
	target := buildDiffState(map[string]map[string]interface{}{})

	diff := tfstate.DiffStates(base, target)
	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	if diff.Diffs[0].Kind != tfstate.DiffRemoved {
		t.Errorf("expected DiffRemoved, got %v", diff.Diffs[0].Kind)
	}
}

func TestDiffStates_Modified(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-old", "instance_type": "t2.micro"},
	})
	target := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-new", "instance_type": "t2.micro"},
	})

	diff := tfstate.DiffStates(base, target)
	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	d := diff.Diffs[0]
	if d.Kind != tfstate.DiffModified {
		t.Errorf("expected DiffModified, got %v", d.Kind)
	}
	if len(d.Changes) != 1 || d.Changes[0].Attribute != "ami" {
		t.Errorf("unexpected changes: %+v", d.Changes)
	}
}

func TestDiffStates_Summary(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{
		"old": {"x": "1"},
		"mod": {"y": "a"},
	})
	target := buildDiffState(map[string]map[string]interface{}{
		"new": {"z": "2"},
		"mod": {"y": "b"},
	})

	diff := tfstate.DiffStates(base, target)
	summary := diff.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if summary != "added=1 removed=1 modified=1" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestDiffStates_NilBase(t *testing.T) {
	target := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-123"},
	})
	diff := tfstate.DiffStates(nil, target)
	if len(diff.Diffs) != 1 || diff.Diffs[0].Kind != tfstate.DiffAdded {
		t.Errorf("expected 1 added diff, got %+v", diff.Diffs)
	}
}
