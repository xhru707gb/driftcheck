package tfstate

import (
	"fmt"
	"io"
	"strings"
)

// WriteTagReport writes a human-readable tag enforcement report to w.
func WriteTagReport(w io.Writer, report *TagReport) error {
	if report == nil || !report.HasViolations() {
		_, err := fmt.Fprintln(w, "Tag enforcement: OK — no violations found.")
		return err
	}

	_, err := fmt.Fprintf(w, "Tag enforcement: %d violation(s) found\n", len(report.Violations))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, strings.Repeat("-", 50))
	if err != nil {
		return err
	}

	for _, v := range report.Violations {
		_, err = fmt.Fprintf(w, "  [!] %s\n", v.String())
		if err != nil {
			return err
		}
		if v.Rule.Message != "" {
			_, err = fmt.Fprintf(w, "      hint: %s\n", v.Rule.Message)
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintln(w, strings.Repeat("-", 50))
	return err
}
