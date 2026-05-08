package drift

import (
	"fmt"

	"github.com/your-org/driftcheck/internal/tfstate"
)

// DriftKind classifies the type of configuration drift.
type DriftKind string

const (
	KindMissing  DriftKind = "MISSING"   // resource exists in state but not in live cloud
	KindModified DriftKind = "MODIFIED"  // attribute value differs
	KindExtra    DriftKind = "EXTRA"     // resource exists in live cloud but not in state
)

// Finding describes a single drift item.
type Finding struct {
	Kind         DriftKind
	ResourceKey  string
	Attribute    string
	Expected     interface{}
	Actual       interface{}
}

func (f Finding) String() string {
	switch f.Kind {
	case KindMissing:
		return fmt.Sprintf("[%s] %s — resource not found in live state", f.Kind, f.ResourceKey)
	case KindExtra:
		return fmt.Sprintf("[%s] %s — resource exists in live state but not in Terraform", f.Kind, f.ResourceKey)
	default:
		return fmt.Sprintf("[%s] %s.%s: expected=%v actual=%v", f.Kind, f.ResourceKey, f.Attribute, f.Expected, f.Actual)
	}
}

// Detector compares Terraform state against live cloud attributes.
type Detector struct{}

// New returns a new Detector.
func New() *Detector { return &Detector{} }

// Compare checks each resource from the Terraform state against liveAttrs.
// liveAttrs is keyed by "type.name" and maps attribute name → live value.
func (d *Detector) Compare(state *tfstate.State, liveAttrs map[string]map[string]interface{}) []Finding {
	var findings []Finding
	planned := state.ResourceMap()

	for key, resource := range planned {
		live, ok := liveAttrs[key]
		if !ok {
			findings = append(findings, Finding{Kind: KindMissing, ResourceKey: key})
			continue
		}
		for attr, expected := range resource.Attributes {
			actual, exists := live[attr]
			if !exists || fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
				findings = append(findings, Finding{
					Kind:        KindModified,
					ResourceKey: key,
					Attribute:   attr,
					Expected:    expected,
					Actual:      actual,
				})
			}
		}
	}

	for key := range liveAttrs {
		if _, ok := planned[key]; !ok {
			findings = append(findings, Finding{Kind: KindExtra, ResourceKey: key})
		}
	}
	return findings
}

// HasDrift returns true if any findings of the given kinds are present.
// If no kinds are specified, it returns true for any finding.
func HasDrift(findings []Finding, kinds ...DriftKind) bool {
	if len(kinds) == 0 {
		return len(findings) > 0
	}
	kindSet := make(map[DriftKind]struct{}, len(kinds))
	for _, k := range kinds {
		kindSet[k] = struct{}{}
	}
	for _, f := range findings {
		if _, ok := kindSet[f.Kind]; ok {
			return true
		}
	}
	return false
}
