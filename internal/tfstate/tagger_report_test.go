package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteTagReport_NoViolations(t *testing.T) {
	report := &TagReport{}
	var buf bytes.Buffer
	if err := WriteTagReport(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK message, got: %s", buf.String())
	}
}

func TestWriteTagReport_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTagReport(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK message, got: %s", buf.String())
	}
}

func TestWriteTagReport_WithViolations(t *testing.T) {
	report := &TagReport{
		Violations: []TagViolation{
			{
				Resource: ResourceKey{Type: "aws_instance", Name: "web"},
				Rule:     TagRule{Key: "env", Message: "set env to prod or dev"},
			},
			{
				Resource: ResourceKey{Type: "aws_s3_bucket", Name: "data"},
				Rule:     TagRule{Key: "team"},
				Actual:   "unknown",
			},
		},
	}
	var buf bytes.Buffer
	if err := WriteTagReport(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 violation") {
		t.Errorf("expected violation count, got: %s", out)
	}
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected resource name in output")
	}
	if !strings.Contains(out, "set env to prod or dev") {
		t.Errorf("expected hint message in output")
	}
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected actual value in output")
	}
}
