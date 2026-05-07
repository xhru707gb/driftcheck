package output

import (
	"encoding/json"
	"io"

	"github.com/driftcheck/internal/drift"
)

// JSONFormatter writes a machine-readable JSON drift report.
type JSONFormatter struct{}

type jsonReport struct {
	DriftDetected bool          `json:"drift_detected"`
	TotalChanges  int           `json:"total_changes"`
	Added         int           `json:"added"`
	Modified      int           `json:"modified"`
	Deleted       int           `json:"deleted"`
	Changes       []jsonChange  `json:"changes"`
}

type jsonChange struct {
	Kind        string            `json:"kind"`
	ResourceKey string            `json:"resource_key"`
	Attributes  map[string]jsonDiff `json:"attributes,omitempty"`
}

type jsonDiff struct {
	Want string `json:"want"`
	Got  string `json:"got"`
}

func (f *JSONFormatter) Write(w io.Writer, report *drift.Report) error {
	out := jsonReport{
		DriftDetected: len(report.Changes) > 0,
		TotalChanges:  len(report.Changes),
		Added:         report.Added,
		Modified:      report.Modified,
		Deleted:       report.Deleted,
		Changes:       make([]jsonChange, 0, len(report.Changes)),
	}

	for _, c := range report.Changes {
		jc := jsonChange{
			Kind:        string(c.Kind),
			ResourceKey: c.ResourceKey,
		}
		if len(c.Attributes) > 0 {
			jc.Attributes = make(map[string]jsonDiff, len(c.Attributes))
			for attr, diff := range c.Attributes {
				jc.Attributes[attr] = jsonDiff{Want: diff.Want, Got: diff.Got}
			}
		}
		out.Changes = append(out.Changes, jc)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
