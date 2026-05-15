package tfstate

import "strings"

// PruneOptions controls which resources are removed from a State.
type PruneOptions struct {
	// RemoveTypes removes all resources whose type is in this list.
	RemoveTypes []string
	// RemoveByPrefix removes resources whose name starts with the given prefix.
	RemoveByPrefix string
	// RemoveOrphans removes resources that have no attributes.
	RemoveOrphans bool
}

// PruneResult summarises what was removed.
type PruneResult struct {
	Removed []ResourceKey
	Kept    int
}

// Prune removes resources from state according to the supplied options and
// returns a new *State together with a PruneResult describing the changes.
func Prune(s *State, opts PruneOptions) (*State, PruneResult) {
	if s == nil {
		return NewState(), PruneResult{}
	}

	typeSet := toTypeSet(opts.RemoveTypes)
	out := NewState()
	result := PruneResult{}

	for _, key := range s.Keys() {
		res, _ := s.Get(key)

		if _, drop := typeSet[res.Type]; drop {
			result.Removed = append(result.Removed, key)
			continue
		}

		if opts.RemoveByPrefix != "" && strings.HasPrefix(res.Name, opts.RemoveByPrefix) {
			result.Removed = append(result.Removed, key)
			continue
		}

		if opts.RemoveOrphans && len(res.Attributes) == 0 {
			result.Removed = append(result.Removed, key)
			continue
		}

		out.Add(res)
		result.Kept++
	}

	return out, result
}

func toTypeSet(types []string) map[string]struct{} {
	m := make(map[string]struct{}, len(types))
	for _, t := range types {
		m[t] = struct{}{}
	}
	return m
}
