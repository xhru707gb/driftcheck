package tfstate

import "fmt"

// TagRule defines a required tag key and optional allowed values.
type TagRule struct {
	Key     string
	Values  []string // empty means any value is accepted
	Message string
}

// TagViolation describes a missing or invalid tag on a resource.
type TagViolation struct {
	Resource ResourceKey
	Rule     TagRule
	Actual   string // empty if tag is missing
}

func (v TagViolation) String() string {
	if v.Actual == "" {
		return fmt.Sprintf("%s: missing required tag %q", v.Resource, v.Rule.Key)
	}
	return fmt.Sprintf("%s: tag %q has disallowed value %q", v.Resource, v.Rule.Key, v.Actual)
}

// TagReport holds all violations found during tag enforcement.
type TagReport struct {
	Violations []TagViolation
}

func (r *TagReport) HasViolations() bool {
	return len(r.Violations) > 0
}

// EnforceTags checks all resources in state against the provided rules.
func EnforceTags(s *State, rules []TagRule) *TagReport {
	report := &TagReport{}
	if s == nil {
		return report
	}
	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}
		attrs := res.Attributes
		for _, rule := range rules {
			tagKey := "tags." + rule.Key
			val, exists := attrs[tagKey]
			if !exists {
				report.Violations = append(report.Violations, TagViolation{
					Resource: key,
					Rule:     rule,
				})
				continue
			}
			if len(rule.Values) > 0 && !contains(rule.Values, val) {
				report.Violations = append(report.Violations, TagViolation{
					Resource: key,
					Rule:     rule,
					Actual:   val,
				})
			}
		}
	}
	return report
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
