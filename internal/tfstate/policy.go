package tfstate

import "fmt"

// PolicyRule defines a single compliance rule applied to resources.
type PolicyRule struct {
	Name        string
	Description string
	Check       func(r Resource) bool
}

// PolicyViolation records a rule failure for a specific resource.
type PolicyViolation struct {
	ResourceKey string
	RuleName    string
	Description string
}

func (v PolicyViolation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.ResourceKey, v.RuleName, v.Description)
}

// PolicyReport holds the results of a policy evaluation.
type PolicyReport struct {
	Violations []PolicyViolation
	Checked    int
}

func (r *PolicyReport) HasViolations() bool {
	return len(r.Violations) > 0
}

// EnforcePolicy evaluates all provided rules against every resource in state.
func EnforcePolicy(s *State, rules []PolicyRule) (*PolicyReport, error) {
	if s == nil {
		return nil, fmt.Errorf("state must not be nil")
	}
	if len(rules) == 0 {
		return &PolicyReport{}, nil
	}

	report := &PolicyReport{}
	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}
		report.Checked++
		for _, rule := range rules {
			if !rule.Check(res) {
				report.Violations = append(report.Violations, PolicyViolation{
					ResourceKey: key.String(),
					RuleName:    rule.Name,
					Description: rule.Description,
				})
			}
		}
	}
	return report, nil
}
