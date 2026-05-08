package tfstate

import "strings"

// FilterOptions controls which resources are included in a filtered State.
type FilterOptions struct {
	Types       []string // include only these resource types
	ExcludeTypes []string // exclude these resource types
	NamePrefix  string   // include only resources whose name starts with this prefix
}

// Filter returns a new State containing only the resources that match opts.
func Filter(s *State, opts FilterOptions) *State {
	includeTypes := toSet(opts.Types)
	excludeTypes := toSet(opts.ExcludeTypes)

	out := NewState()
	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		// type allow-list
		if len(includeTypes) > 0 {
			if _, allowed := includeTypes[res.Type]; !allowed {
				continue
			}
		}

		// type deny-list
		if _, excluded := excludeTypes[res.Type]; excluded {
			continue
		}

		// name prefix
		if opts.NamePrefix != "" && !strings.HasPrefix(res.Name, opts.NamePrefix) {
			continue
		}

		out.Add(res)
	}
	return out
}

// toSet converts a string slice to a set (map[string]struct{}).
func toSet(items []string) map[string]struct{} {
	m := make(map[string]struct{}, len(items))
	for _, v := range items {
		m[v] = struct{}{}
	}
	return m
}
