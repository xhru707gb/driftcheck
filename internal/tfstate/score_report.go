package tfstate

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteScoreReport writes a formatted drift score report to w.
func WriteScoreReport(w io.Writer, score *DriftScore) error {
	if score == nil {
		return fmt.Errorf("score_report: score must not be nil")
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "Drift Score Report")
	fmt.Fprintln(tw, "==================")
	fmt.Fprintf(tw, "Severity:\t%s\n", score.Severity())
	fmt.Fprintf(tw, "Total Score:\t%.1f\n", score.Total)
	fmt.Fprintln(tw, "------------------")
	fmt.Fprintf(tw, "Added:\t%d\n", score.Added)
	fmt.Fprintf(tw, "Removed:\t%d\n", score.Removed)
	fmt.Fprintf(tw, "Modified:\t%d\n", score.Modified)

	return tw.Flush()
}
