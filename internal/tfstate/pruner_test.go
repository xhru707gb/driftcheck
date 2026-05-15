package tfstate

import (
	"bytes"
	"testing"
)

func buildPrunerState() *State {
	s := NewState()
	s.Add(Resource{Type: "aws_instance", Name: "web", ID: "i-1", Attributes: map[string]string{"ami": "ami-abc"}})
	s.Add(Resource{Type: "aws_instance", Name: "worker", ID: "i-2", Attributes: map[string]string{"ami": "ami-def"}})
	s.Add(Resource{Type: "aws_s3_bucket", Name: "assets", ID: "bkt-1", Attributes: map[string]string{"region": "us-east-1"}})
	s.Add(Resource{Type: "aws_sg", Name: "tmp_firewall", ID: "sg-1", Attributes: map[string]string{}})
	s.Add(Resource{Type: "aws_iam_role", Name: "orphan", ID: "r-1", Attributes: map[string]string{}})
	return s
}

func TestPrune_NilState(t *testing.T) {
	out, res := Prune(nil, PruneOptions{})
	if out == nil {
		t.Fatal("expected non-nil state")
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
}

func TestPrune_NoOptions(t *testing.T) {
	s := buildPrunerState()
	out, res := Prune(s, PruneOptions{})
	if len(res.Removed) != 0 {
		t.Errorf("expected nothing removed, got %d", len(res.Removed))
	}
	if res.Kept != len(s.Keys()) {
		t.Errorf("kept mismatch: want %d got %d", len(s.Keys()), res.Kept)
	}
	_ = out
}

func TestPrune_ByType(t *testing.T) {
	s := buildPrunerState()
	out, res := Prune(s, PruneOptions{RemoveTypes: []string{"aws_instance"}})
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	for _, key := range out.Keys() {
		r, _ := out.Get(key)
		if r.Type == "aws_instance" {
			t.Errorf("aws_instance should have been pruned: %s", key)
		}
	}
}

func TestPrune_ByPrefix(t *testing.T) {
	s := buildPrunerState()
	_, res := Prune(s, PruneOptions{RemoveByPrefix: "tmp_"})
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(res.Removed))
	}
}

func TestPrune_Orphans(t *testing.T) {
	s := buildPrunerState()
	_, res := Prune(s, PruneOptions{RemoveOrphans: true})
	// tmp_firewall (aws_sg) and orphan (aws_iam_role) have empty attributes
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 orphans removed, got %d", len(res.Removed))
	}
}

func TestPrune_Combined(t *testing.T) {
	s := buildPrunerState()
	_, res := Prune(s, PruneOptions{
		RemoveTypes:    []string{"aws_s3_bucket"},
		RemoveOrphans:  true,
	})
	// s3 bucket + 2 orphans = 3
	if len(res.Removed) != 3 {
		t.Errorf("expected 3 removed, got %d", len(res.Removed))
	}
}

func TestWritePruneReport_NilResult(t *testing.T) {
	var buf bytes.Buffer
	WritePruneReport(&buf, nil)
	if buf.Len() == 0 {
		t.Error("expected non-empty output for nil result")
	}
}

func TestWritePruneReport_WithRemovals(t *testing.T) {
	s := buildPrunerState()
	_, res := Prune(s, PruneOptions{RemoveTypes: []string{"aws_instance"}})
	var buf bytes.Buffer
	WritePruneReport(&buf, &res)
	out := buf.String()
	if out == "" {
		t.Error("expected report output")
	}
}
