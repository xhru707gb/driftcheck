package tfstate

import "strings"

// FilterOptions controls which resources are included in drift checks.
type FilterOptions struct {
	// ResourceTypes limits checking to specific resource types (e.g. "aws_instance").
	// If empty, all types are included.
	ResourceTypes []string

	// NamePrefix limits checking to resources whose name starts with the given prefix.
	// If empty, all names are included.
	NamePrefix string

	// ExcludeTypes lists resource types to skip entirely.
	ExcludeTypes []string
}

// Filter returns a new State containing only the resources that match the
// given FilterOptions.
func (s *State) Filter(opts FilterOptions) *State {
	includeTypes := toSet(opts.ResourceTypes)
	excludeTypes := toSet(opts.ExcludeTypes)

	filtered := NewState()

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		// Apply type exclusion first.
		if len(excludeTypes) > 0 {
			if _, excluded := excludeTypes[key.Type]; excluded {
				continue
			}
		}

		// Apply type inclusion filter.
		if len(includeTypes) > 0 {
			if _, included := includeTypes[key.Type]; !included {
				continue
			}
		}

		// Apply name prefix filter.
		if opts.NamePrefix != "" && !strings.HasPrefix(key.Name, opts.NamePrefix) {
			continue
		}

		filtered.Add(res)
	}

	return filtered
}

func toSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}
