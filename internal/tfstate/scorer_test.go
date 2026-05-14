package tfstate

import (
	"testing"
)

func buildScorerState(resources map[string]map[string]interface{}) *State {
	s := NewState()
	for id, attrs := range resources {
		s.Add(Resource{
			Type:       "aws_instance",
			Name:       id,
			ID:         id,
			Attributes: attrs,
		})
	}
	return s
}

func TestScore_NoDrift(t *testing.T) {
	attrs := map[string]interface{}{"ami": "ami-123"}
	a := buildScorerState(map[string]map[string]interface{}{"res1": attrs})
	b := buildScorerState(map[string]map[string]interface{}{"res1": attrs})

	score, err := Score(a, b, DefaultWeights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.Total != 0 {
		t.Errorf("expected score 0, got %.1f", score.Total)
	}
	if score.Severity() != "none" {
		t.Errorf("expected severity none, got %s", score.Severity())
	}
}

func TestScore_AddedResource(t *testing.T) {
	attrs := map[string]interface{}{"ami": "ami-123"}
	a := buildScorerState(map[string]map[string]interface{}{})
	b := buildScorerState(map[string]map[string]interface{}{"res1": attrs})

	score, err := Score(a, b, DefaultWeights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.Added != 1 {
		t.Errorf("expected 1 added, got %d", score.Added)
	}
	if score.Total != DefaultWeights.Added {
		t.Errorf("expected total %.1f, got %.1f", DefaultWeights.Added, score.Total)
	}
}

func TestScore_RemovedResource(t *testing.T) {
	attrs := map[string]interface{}{"ami": "ami-123"}
	a := buildScorerState(map[string]map[string]interface{}{"res1": attrs})
	b := buildScorerState(map[string]map[string]interface{}{})

	score, err := Score(a, b, DefaultWeights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", score.Removed)
	}
	if score.Severity() != "low" {
		t.Errorf("expected severity low, got %s", score.Severity())
	}
}

func TestScore_NilState(t *testing.T) {
	_, err := Score(nil, NewState(), DefaultWeights)
	if err == nil {
		t.Error("expected error for nil state")
	}
}

func TestScore_HighSeverity(t *testing.T) {
	a := buildScorerState(map[string]map[string]interface{}{
		"r1": {"x": "1"}, "r2": {"x": "2"}, "r3": {"x": "3"},
		"r4": {"x": "4"}, "r5": {"x": "5"},
	})
	b := buildScorerState(map[string]map[string]interface{}{})

	score, err := Score(a, b, DefaultWeights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.Severity() != "high" {
		t.Errorf("expected high severity, got %s", score.Severity())
	}
}

func TestDriftScore_String(t *testing.T) {
	s := &DriftScore{Total: 4.0, Added: 2, Removed: 0, Modified: 1}
	out := s.String()
	if out == "" {
		t.Error("expected non-empty string")
	}
}
