package tfstate

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteGroupReport writes a formatted summary of a GroupResult to w.
func WriteGroupReport(w io.Writer, result *GroupResult, by GroupBy) {
	if result == nil {
		fmt.Fprintln(w, "no group result available")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintf(tw, "Grouped by: %s\n", by)
	fmt.Fprintf(tw, "Total resources: %d\n\n", result.Total)
	fmt.Fprintf(tw, "%-30s\t%s\n", "GROUP", "COUNT")
	fmt.Fprintf(tw, "%-30s\t%s\n", "-----", "-----")

	for _, g := range result.Groups {
		label := g.Key
		if label == "" {
			label = "(unset)"
		}
		fmt.Fprintf(tw, "%-30s\t%d\n", label, len(g.Resources))
	}
}
