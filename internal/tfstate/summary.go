package tfstate

import "fmt"

// TypeSummary holds aggregated counts for a single resource type.
type TypeSummary struct {
	Type      string
	Count     int
	Instances int
}

// StateSummary provides an overview of all resources in a State.
type StateSummary struct {
	TotalResources  int
	TotalInstances  int
	ByType          map[string]*TypeSummary
	TypeOrder       []string
}

// Summarize builds a StateSummary from the given State.
func Summarize(s *State) *StateSummary {
	sum := &StateSummary{
		ByType: make(map[string]*TypeSummary),
	}

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		sum.TotalResources++
		sum.TotalInstances += len(res.Instances)

		ts, exists := sum.ByType[res.Type]
		if !exists {
			ts = &TypeSummary{Type: res.Type}
			sum.ByType[res.Type] = ts
			sum.TypeOrder = append(sum.TypeOrder, res.Type)
		}
		ts.Count++
		ts.Instances += len(res.Instances)
	}

	return sum
}

// String returns a human-readable summary.
func (s *StateSummary) String() string {
	out := fmt.Sprintf("Resources: %d  Instances: %d\n", s.TotalResources, s.TotalInstances)
	for _, t := range s.TypeOrder {
		ts := s.ByType[t]
		out += fmt.Sprintf("  %-40s resources: %d  instances: %d\n", ts.Type, ts.Count, ts.Instances)
	}
	return out
}
