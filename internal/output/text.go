package output

import (
	"fmt"
	"io"

	"github.com/driftcheck/internal/drift"
)

// TextFormatter writes a human-readable plain-text drift report.
type TextFormatter struct{}

func (f *TextFormatter) Write(w io.Writer, report *drift.Report) error {
	if len(report.Changes) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}

	fmt.Fprintf(w, "Drift detected: %d change(s)\n\n", len(report.Changes))

	for _, c := range report.Changes {
		fmt.Fprintf(w, "  [%s] %s\n", changeSymbol(c.Kind), c.ResourceKey)
		if c.Kind == drift.Modified {
			for attr, diff := range c.Attributes {
				fmt.Fprintf(w, "      ~ %s: %q => %q\n", attr, diff.Got, diff.Want)
			}
		}
	}

	fmt.Fprintf(w, "\nSummary: %d added, %d modified, %d deleted\n",
		report.Added, report.Modified, report.Deleted)
	return nil
}

func changeSymbol(k drift.ChangeKind) string {
	switch k {
	case drift.Added:
		return "+"
	case drift.Deleted:
		return "-"
	case drift.Modified:
		return "~"
	default:
		return "?"
	}
}
