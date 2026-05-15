package tfstate

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ExportFormat defines the output format for state export.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "csv"
	ExportJSON ExportFormat = "json"
)

// ExportRow represents a single resource row in an export.
type ExportRow struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// Export writes the state resources to w in the given format.
func Export(s *State, format ExportFormat, w io.Writer) error {
	if s == nil {
		return fmt.Errorf("export: state is nil")
	}

	rows := buildRows(s)

	switch format {
	case ExportCSV:
		return exportCSV(rows, w)
	case ExportJSON:
		return exportJSON(rows, w)
	default:
		return fmt.Errorf("export: unsupported format %q", format)
	}
}

func buildRows(s *State) []ExportRow {
	keys := s.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	rows := make([]ExportRow, 0, len(keys))
	for _, k := range keys {
		res, ok := s.Get(k)
		if !ok {
			continue
		}
		attrs := make(map[string]string, len(res.Attributes))
		for ak, av := range res.Attributes {
			attrs[ak] = fmt.Sprintf("%v", av)
		}
		rows = append(rows, ExportRow{
			Type:       k.Type,
			Name:       k.Name,
			ID:         res.ID,
			Attributes: attrs,
		})
	}
	return rows
}

func exportCSV(rows []ExportRow, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"type", "name", "id"}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.Write([]string{r.Type, r.Name, r.ID}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportJSON(rows []ExportRow, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
