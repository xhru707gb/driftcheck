package tfstate

import (
	"fmt"
	"sort"
)

// DiffKind describes the type of difference between two states.
type DiffKind string

const (
	DiffAdded    DiffKind = "added"
	DiffRemoved  DiffKind = "removed"
	DiffModified DiffKind = "modified"
)

// ResourceDiff represents a single resource-level difference.
type ResourceDiff struct {
	Key      ResourceKey
	Kind     DiffKind
	OldAttrs map[string]interface{}
	NewAttrs map[string]interface{}
	Changes  []AttrChange
}

// AttrChange describes a changed attribute.
type AttrChange struct {
	Attribute string
	OldValue  interface{}
	NewValue  interface{}
}

// StateDiff holds all differences between two states.
type StateDiff struct {
	Diffs []ResourceDiff
}

// HasChanges returns true if any differences exist.
func (d *StateDiff) HasChanges() bool {
	return len(d.Diffs) > 0
}

// Summary returns a human-readable summary string.
func (d *StateDiff) Summary() string {
	added, removed, modified := 0, 0, 0
	for _, diff := range d.Diffs {
		switch diff.Kind {
		case DiffAdded:
			added++
		case DiffRemoved:
			removed++
		case DiffModified:
			modified++
		}
	}
	return fmt.Sprintf("added=%d removed=%d modified=%d", added, removed, modified)
}

// DiffStates compares two States and returns all differences.
func DiffStates(base, target *State) *StateDiff {
	result := &StateDiff{}

	baseKeys := toKeySet(base)
	targetKeys := toKeySet(target)

	// Find removed and modified
	for _, key := range sortedKeys(baseKeys) {
		baseRes, _ := base.Get(key)
		if _, exists := targetKeys[key]; !exists {
			result.Diffs = append(result.Diffs, ResourceDiff{
				Key:      key,
				Kind:     DiffRemoved,
				OldAttrs: baseRes.Attributes,
			})
			continue
		}
		targetRes, _ := target.Get(key)
		changes := attrDiff(baseRes.Attributes, targetRes.Attributes)
		if len(changes) > 0 {
			result.Diffs = append(result.Diffs, ResourceDiff{
				Key:      key,
				Kind:     DiffModified,
				OldAttrs: baseRes.Attributes,
				NewAttrs: targetRes.Attributes,
				Changes:  changes,
			})
		}
	}

	// Find added
	for _, key := range sortedKeys(targetKeys) {
		if _, exists := baseKeys[key]; !exists {
			targetRes, _ := target.Get(key)
			result.Diffs = append(result.Diffs, ResourceDiff{
				Key:      key,
				Kind:     DiffAdded,
				NewAttrs: targetRes.Attributes,
			})
		}
	}

	return result
}

func attrDiff(old, new map[string]interface{}) []AttrChange {
	var changes []AttrChange
	keys := make(map[string]struct{})
	for k := range old {
		keys[k] = struct{}{}
	}
	for k := range new {
		keys[k] = struct{}{}
	}
	for k := range keys {
		ov := old[k]
		nv := new[k]
		if fmt.Sprintf("%v", ov) != fmt.Sprintf("%v", nv) {
			changes = append(changes, AttrChange{Attribute: k, OldValue: ov, NewValue: nv})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Attribute < changes[j].Attribute })
	return changes
}

func toKeySet(s *State) map[ResourceKey]struct{} {
	set := make(map[ResourceKey]struct{})
	if s == nil {
		return set
	}
	for _, k := range s.Keys() {
		set[k] = struct{}{}
	}
	return set
}

func sortedKeys(m map[ResourceKey]struct{}) []ResourceKey {
	keys := make([]ResourceKey, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	return keys
}
