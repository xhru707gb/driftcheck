package tfstate

import "fmt"

// DriftScore represents a weighted severity score for a set of drift findings.
type DriftScore struct {
	Total    float64
	Added    int
	Removed  int
	Modified int
}

// ScoreWeights controls how much each drift type contributes to the total score.
type ScoreWeights struct {
	Added    float64
	Removed  float64
	Modified float64
}

// DefaultWeights provides sensible defaults: removals are most severe.
var DefaultWeights = ScoreWeights{
	Added:    1.0,
	Removed:  3.0,
	Modified: 2.0,
}

// Score computes a drift severity score from two states using the provided weights.
// It uses DiffStates internally to determine what changed.
func Score(baseline, current *State, w ScoreWeights) (*DriftScore, error) {
	if baseline == nil || current == nil {
		return nil, fmt.Errorf("scorer: both states must be non-nil")
	}

	result, err := DiffStates(baseline, current)
	if err != nil {
		return nil, fmt.Errorf("scorer: diff failed: %w", err)
	}

	s := &DriftScore{}
	for _, d := range result {
		switch d.ChangeType {
		case ChangeAdded:
			s.Added++
		case ChangeRemoved:
			s.Removed++
		case ChangeModified:
			s.Modified++
		}
	}

	s.Total = float64(s.Added)*w.Added +
		float64(s.Removed)*w.Removed +
		float64(s.Modified)*w.Modified

	return s, nil
}

// Severity returns a human-readable severity label based on the total score.
func (s *DriftScore) Severity() string {
	switch {
	case s.Total == 0:
		return "none"
	case s.Total < 5:
		return "low"
	case s.Total < 15:
		return "medium"
	default:
		return "high"
	}
}

// String returns a compact summary of the score.
func (s *DriftScore) String() string {
	return fmt.Sprintf("score=%.1f severity=%s (added=%d removed=%d modified=%d)",
		s.Total, s.Severity(), s.Added, s.Removed, s.Modified)
}
