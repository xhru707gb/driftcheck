package cloud

import (
	"fmt"
	"io"
	"strings"
)

// WriteReport writes a human-readable reconciliation report to w.
func WriteReport(w io.Writer, result *ReconcileResult) {
	if len(result.Drifts) == 0 && len(result.Missing) == 0 {
		fmt.Fprintln(w, "✓ No drift detected. Infrastructure matches Terraform state.")
		return
	}

	if len(result.Missing) > 0 {
		fmt.Fprintf(w, "\n[MISSING] %d resource(s) not found in cloud:\n", len(result.Missing))
		for _, id := range result.Missing {
			fmt.Fprintf(w, "  - %s\n", id)
		}
	}

	if len(result.Drifts) > 0 {
		fmt.Fprintf(w, "\n[DRIFT] %d attribute(s) differ:\n", len(result.Drifts))
		for _, d := range result.Drifts {
			fmt.Fprintf(w, "  %s / %s\n", d.ResourceType, d.ResourceID)
			fmt.Fprintf(w, "    attribute : %s\n", d.Attribute)
			fmt.Fprintf(w, "    expected  : %v\n", d.Expected)
			fmt.Fprintf(w, "    actual    : %v\n", d.Actual)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "Summary: %d drift(s), %d missing\n",
		len(result.Drifts), len(result.Missing))
}
