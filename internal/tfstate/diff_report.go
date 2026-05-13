package tfstate

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteDiffReport writes a human-readable diff report to w.
func WriteDiffReport(w io.Writer, diff *StateDiff) error {
	if !diff.HasChanges() {
		_, err := fmt.Fprintln(w, "No state differences detected.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintf(tw, "%-10s\t%-40s\t%s\n", "CHANGE", "RESOURCE", "DETAILS")
	fmt.Fprintf(tw, "%-10s\t%-40s\t%s\n", "------", "--------", "-------")

	for _, d := range diff.Diffs {
		symbol := diffSymbol(d.Kind)
		details := diffDetails(d)
		fmt.Fprintf(tw, "%-10s\t%-40s\t%s\n", symbol, d.Key.String(), details)
	}

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, "Summary: %s\n", diff.Summary())
	return nil
}

func diffSymbol(kind DiffKind) string {
	switch kind {
	case DiffAdded:
		return "[+] added"
	case DiffRemoved:
		return "[-] removed"
	case DiffModified:
		return "[~] modified"
	default:
		return "[?] unknown"
	}
}

func diffDetails(d ResourceDiff) string {
	switch d.Kind {
	case DiffAdded:
		return fmt.Sprintf("%d attribute(s)", len(d.NewAttrs))
	case DiffRemoved:
		return fmt.Sprintf("%d attribute(s)", len(d.OldAttrs))
	case DiffModified:
		if len(d.Changes) == 0 {
			return "no attribute changes"
		}
		return fmt.Sprintf("%s: %v -> %v (+%d more)",
			d.Changes[0].Attribute,
			d.Changes[0].OldValue,
			d.Changes[0].NewValue,
			len(d.Changes)-1,
		)
	default:
		return ""
	}
}
