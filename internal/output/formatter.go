package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftcheck/internal/drift"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

// Formatter writes drift reports in a specific format.
type Formatter interface {
	Write(w io.Writer, report *drift.Report) error
}

// New returns a Formatter for the given format string.
// Returns an error if the format is not recognized.
func New(format Format) (Formatter, error) {
	switch format {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatTable:
		return &TableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown output format %q: must be one of [%s]",
			format, strings.Join([]string{string(FormatText), string(FormatJSON), string(FormatTable)}, ", "))
	}
}
