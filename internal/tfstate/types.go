package tfstate

import "sort"

// TypeInventory holds a count of resources grouped by their type.
type TypeInventory struct {
	counts map[string]int
}

// BuildTypeInventory scans the given State and returns a TypeInventory
// mapping each resource type to the number of instances present.
func BuildTypeInventory(s *State) *TypeInventory {
	inv := &TypeInventory{counts: make(map[string]int)}
	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}
		inv.counts[res.Type]++
	}
	return inv
}

// Count returns the number of resources for the given type.
func (t *TypeInventory) Count(resourceType string) int {
	return t.counts[resourceType]
}

// Types returns a sorted slice of all resource types present.
func (t *TypeInventory) Types() []string {
	types := make([]string, 0, len(t.counts))
	for rt := range t.counts {
		types = append(types, rt)
	}
	sort.Strings(types)
	return types
}

// Total returns the total number of resources across all types.
func (t *TypeInventory) Total() int {
	n := 0
	for _, c := range t.counts {
		n += c
	}
	return n
}

// Has returns true if the inventory contains at least one resource of the given type.
func (t *TypeInventory) Has(resourceType string) bool {
	return t.counts[resourceType] > 0
}
