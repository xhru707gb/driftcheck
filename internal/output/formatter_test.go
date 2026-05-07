package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/output"
)

func makeReport() *drift.Report {
	return &drift.Report{
		Added:    1,
		Modified: 1,
		Deleted:  0,
		Changes: []drift.Change{
			{Kind: drift.Added, ResourceKey: "aws_instance.web"},
			{
				Kind:        drift.Modified,
				ResourceKey: "aws_s3_bucket.data",
				Attributes: map[string]drift.AttributeDiff{
					"acl": {Want: "private", Got: "public-read"},
				},
			},
		},
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []output.Format{output.FormatText, output.FormatJSON, output.FormatTable} {
		fmt, err := output.New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
		if fmt == nil {
			t.Errorf("New(%q) returned nil formatter", f)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := output.New("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestTextFormatter_NoDrift(t *testing.T) {
	f, _ := output.New(output.FormatText)
	var buf bytes.Buffer
	if err := f.Write(&buf, &drift.Report{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestTextFormatter_WithDrift(t *testing.T) {
	f, _ := output.New(output.FormatText)
	var buf bytes.Buffer
	if err := f.Write(&buf, makeReport()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Error("expected resource key in text output")
	}
	if !strings.Contains(out, "acl") {
		t.Error("expected attribute name in text output")
	}
}

func TestJSONFormatter_WithDrift(t *testing.T) {
	f, _ := output.New(output.FormatJSON)
	var buf bytes.Buffer
	if err := f.Write(&buf, makeReport()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift_detected": true`) {
		t.Errorf("expected drift_detected true in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"total_changes": 2`) {
		t.Errorf("expected total_changes 2 in JSON, got: %s", out)
	}
}
