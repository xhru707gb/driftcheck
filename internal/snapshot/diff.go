package snapshot

// DiffKind describes the kind of difference found between two snapshots.
type DiffKind string

const (
	DiffAdded    DiffKind = "added"
	DiffRemoved  DiffKind = "removed"
	DiffModified DiffKind = "modified"
)

// AttributeChange holds the before/after values for a single attribute.
type AttributeChange struct {
	Attribute string
	OldValue  string
	NewValue  string
}

// ResourceDiff describes all changes for a single resource between two snapshots.
type ResourceDiff struct {
	Key        string
	Kind       DiffKind
	Attributes []AttributeChange
}

// Compare returns the differences between an older and a newer snapshot.
func Compare(old, new *Snapshot) []ResourceDiff {
	var diffs []ResourceDiff

	// Resources present in old — check for removals and modifications.
	for key, oldAttrs := range old.Resources {
		newAttrs, exists := new.Resources[key]
		if !exists {
			diffs = append(diffs, ResourceDiff{Key: key, Kind: DiffRemoved})
			continue
		}
		if changes := attrChanges(oldAttrs, newAttrs); len(changes) > 0 {
			diffs = append(diffs, ResourceDiff{Key: key, Kind: DiffModified, Attributes: changes})
		}
	}

	// Resources present only in new — additions.
	for key := range new.Resources {
		if _, exists := old.Resources[key]; !exists {
			diffs = append(diffs, ResourceDiff{Key: key, Kind: DiffAdded})
		}
	}

	return diffs
}

func attrChanges(old, new map[string]string) []AttributeChange {
	var changes []AttributeChange
	for k, oldVal := range old {
		if newVal, ok := new[k]; !ok || newVal != oldVal {
			changes = append(changes, AttributeChange{Attribute: k, OldValue: oldVal, NewValue: new[k]})
		}
	}
	for k, newVal := range new {
		if _, ok := old[k]; !ok {
			changes = append(changes, AttributeChange{Attribute: k, OldValue: "", NewValue: newVal})
		}
	}
	return changes
}
