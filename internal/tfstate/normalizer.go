package tfstate

import (
	"fmt"
	"strings"
)

// NormalizeResult holds the outcome of normalizing a State.
type NormalizeResult struct {
	Normalized int
	Changes    []string
}

func (r *NormalizeResult) String() string {
	if r.Normalized == 0 {
		return "no normalization changes"
	}
	return fmt.Sprintf("%d resource(s) normalized:\n  - %s",
		r.Normalized, strings.Join(r.Changes, "\n  - "))
}

// Normalize cleans up a State in-place by trimming whitespace from attribute
// values and lower-casing resource type names. It returns a NormalizeResult
// describing every change made.
func Normalize(s *State) (*NormalizeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("normalize: state must not be nil")
	}

	result := &NormalizeResult{}

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		changed := false

		// Normalise type to lower-case.
		normType := strings.ToLower(res.Type)
		if normType != res.Type {
			res.Type = normType
			changed = true
		}

		// Trim whitespace from every attribute value.
		for k, v := range res.Attributes {
			if sv, ok := v.(string); ok {
				trimmed := strings.TrimSpace(sv)
				if trimmed != sv {
					res.Attributes[k] = trimmed
					changed = true
				}
			}
		}

		if changed {
			s.Add(res)
			result.Normalized++
			result.Changes = append(result.Changes,
				fmt.Sprintf("%s (%s)", res.Name, res.Type))
		}
	}

	return result, nil
}
