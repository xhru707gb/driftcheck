package tfstate

import (
	"fmt"
	"io"
)

// WritePolicyReport writes a human-readable policy report to w.
func WritePolicyReport(w io.Writer, report *PolicyReport) {
	if report == nil {
		fmt.Fprintln(w, "no policy report available")
		return
	}

	fmt.Fprintf(w, "Policy check: %d resource(s) evaluated\n", report.Checked)

	if !report.HasViolations() {
		fmt.Fprintln(w, "✓ All resources comply with policy rules.")
		return
	}

	fmt.Fprintf(w, "✗ %d violation(s) found:\n\n", len(report.Violations))
	for _, v := range report.Violations {
		fmt.Fprintf(w, "  • %s\n", v.String())
	}
}
