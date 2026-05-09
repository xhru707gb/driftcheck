package tfstate

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError holds a list of validation issues found in a State.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("state validation failed with %d issue(s):\n  - %s",
		len(e.Issues), strings.Join(e.Issues, "\n  - "))
}

// Validate checks a State for common structural problems and returns a
// *ValidationError if any issues are found, or nil when the state is clean.
func Validate(s *State) error {
	if s == nil {
		return errors.New("state must not be nil")
	}

	var issues []string

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		if res.Type == "" {
			issues = append(issues, fmt.Sprintf("resource %q has an empty type", key))
		}

		if res.Name == "" {
			issues = append(issues, fmt.Sprintf("resource %q has an empty name", key))
		}

		if len(res.Attributes) == 0 {
			issues = append(issues, fmt.Sprintf("resource %q has no attributes", key))
		}

		if _, hasID := res.Attributes["id"]; !hasID {
			issues = append(issues, fmt.Sprintf("resource %q is missing required attribute \"id\"", key))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
