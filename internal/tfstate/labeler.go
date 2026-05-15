package tfstate

import (
	"fmt"
	"strings"
)

// LabelRule defines a rule for applying a computed label to a resource.
type LabelRule struct {
	Key      string
	Prefix   string
	Suffix   string
	FromAttr string // derive label value from this attribute
}

// LabelResult holds the labeling outcome for a single resource.
type LabelResult struct {
	ResourceKey string
	Applied     []string
	Skipped     []string
	Errors      []string
}

// LabelReport aggregates results from ApplyLabels.
type LabelReport struct {
	Results []LabelResult
	Total   int
	Labeled int
}

// ApplyLabels applies a set of LabelRules to all resources in the state.
// It returns a LabelReport describing what was applied, skipped, or errored.
func ApplyLabels(s *State, rules []LabelRule) (*LabelReport, error) {
	if s == nil {
		return nil, fmt.Errorf("ApplyLabels: state must not be nil")
	}

	report := &LabelReport{}

	for _, key := range s.Keys() {
		res, _ := s.Get(key)
		result := LabelResult{ResourceKey: key.String()}

		for _, rule := range rules {
			if rule.Key == "" {
				result.Errors = append(result.Errors, "rule has empty key")
				continue
			}

			var val string
			if rule.FromAttr != "" {
				attrVal, ok := res.Attributes[rule.FromAttr]
				if !ok {
					result.Skipped = append(result.Skipped, rule.Key)
					continue
				}
				val = fmt.Sprintf("%v", attrVal)
			} else {
				val = res.ID
			}

			computed := strings.TrimSpace(rule.Prefix + val + rule.Suffix)
			if res.Attributes == nil {
				res.Attributes = make(map[string]interface{})
			}
			res.Attributes[rule.Key] = computed
			result.Applied = append(result.Applied, rule.Key)
		}

		if len(result.Applied) > 0 {
			report.Labeled++
		}
		report.Results = append(report.Results, result)
		report.Total++
	}

	return report, nil
}
