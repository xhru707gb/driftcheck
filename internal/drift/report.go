package drift

import (
	"fmt"
	"io"
	"strings"
)

// ReportFormat controls the output style of a drift report.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatSummary ReportFormat = "summary"
)

// Report writes drift findings to w in the requested format.
func Report(w io.Writer, findings []Finding, format ReportFormat) error {
	switch format {
	case FormatSummary:
		return writeSummary(w, findings)
	default:
		return writeText(w, findings)
	}
}

func writeText(w io.Writer, findings []Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "✅ No drift detected.")
		return err
	}
	_, err := fmt.Fprintf(w, "⚠️  Drift detected — %d finding(s):\n", len(findings))
	if err != nil {
		return err
	}
	for i, f := range findings {
		_, err = fmt.Fprintf(w, "  %d. %s\n", i+1, f.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func writeSummary(w io.Writer, findings []Finding) error {
	counts := map[DriftKind]int{}
	for _, f := range findings {
		counts[f.Kind]++
	}
	parts := []string{}
	for _, k := range []DriftKind{KindMissing, KindModified, KindExtra} {
		if n := counts[k]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", k, n))
		}
	}
	if len(parts) == 0 {
		_, err := fmt.Fprintln(w, "no drift")
		return err
	}
	_, err := fmt.Fprintln(w, strings.Join(parts, " "))
	return err
}
