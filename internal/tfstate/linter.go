package tfstate

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning or error found in a state.
type LintIssue struct {
	Severity string // "warn" or "error"
	Resource string
	Message  string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Severity), i.Resource, i.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

// Lint inspects a State and returns a LintResult with any detected issues.
func Lint(s *State) *LintResult {
	result := &LintResult{}

	if s == nil {
		result.Issues = append(result.Issues, LintIssue{
			Severity: "error",
			Resource: "<state>",
			Message:  "state is nil",
		})
		return result
	}

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		if res.ID == "" {
			result.Issues = append(result.Issues, LintIssue{
				Severity: "error",
				Resource: key.String(),
				Message:  "resource has empty ID",
			})
		}

		if res.Type == "" {
			result.Issues = append(result.Issues, LintIssue{
				Severity: "error",
				Resource: key.String(),
				Message:  "resource has empty type",
			})
		}

		if len(res.Attributes) == 0 {
			result.Issues = append(result.Issues, LintIssue{
				Severity: "warn",
				Resource: key.String(),
				Message:  "resource has no attributes",
			})
		}
	}

	return result
}
