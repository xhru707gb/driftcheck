package tfstate

import "fmt"

// BaselineDiffKind categorises a change relative to a baseline.
type BaselineDiffKind string

const (
	BaselineAdded    BaselineDiffKind = "added"
	BaselineRemoved  BaselineDiffKind = "removed"
	BaselineModified BaselineDiffKind = "modified"
)

// BaselineDiffEntry describes a single divergence from the baseline.
type BaselineDiffEntry struct {
	Kind         BaselineDiffKind
	ResourceKey  string
	Attribute    string // empty for added/removed
	BaselineVal  interface{}
	CurrentVal   interface{}
}

func (e BaselineDiffEntry) String() string {
	switch e.Kind {
	case BaselineAdded:
		return fmt.Sprintf("+ %s (new resource)", e.ResourceKey)
	case BaselineRemoved:
		return fmt.Sprintf("- %s (removed resource)", e.ResourceKey)
	case BaselineModified:
		return fmt.Sprintf("~ %s .%s: %v -> %v", e.ResourceKey, e.Attribute, e.BaselineVal, e.CurrentVal)
	}
	return ""
}

// CompareToBaseline returns the list of differences between a live State and a
// previously captured Baseline.
func CompareToBaseline(b *Baseline, current *State) ([]BaselineDiffEntry, error) {
	if b == nil {
		return nil, fmt.Errorf("baseline must not be nil")
	}
	if current == nil {
		return nil, fmt.Errorf("current state must not be nil")
	}

	var diffs []BaselineDiffEntry

	// Check for removed or modified resources.
	for key, br := range b.Resources {
		cr, ok := current.GetByKey(key)
		if !ok {
			diffs = append(diffs, BaselineDiffEntry{Kind: BaselineRemoved, ResourceKey: key})
			continue
		}
		for ak, bv := range br.Attributes {
			cv, exists := cr.Attributes[ak]
			if !exists || fmt.Sprintf("%v", cv) != fmt.Sprintf("%v", bv) {
				diffs = append(diffs, BaselineDiffEntry{
					Kind:        BaselineModified,
					ResourceKey: key,
					Attribute:   ak,
					BaselineVal: bv,
					CurrentVal:  cv,
				})
			}
		}
	}

	// Check for added resources.
	for _, rk := range current.Keys() {
		if _, ok := b.Resources[rk.String()]; !ok {
			diffs = append(diffs, BaselineDiffEntry{Kind: BaselineAdded, ResourceKey: rk.String()})
		}
	}

	return diffs, nil
}
