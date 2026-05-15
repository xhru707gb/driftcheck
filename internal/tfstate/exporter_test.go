package tfstate

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildExporterState() *State {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		ID:         "i-abc123",
		Attributes: map[string]interface{}{"ami": "ami-123", "instance_type": "t2.micro"},
	})
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "data",
		ID:         "my-bucket",
		Attributes: map[string]interface{}{"region": "us-east-1"},
	})
	return s
}

func TestExport_NilState(t *testing.T) {
	var buf bytes.Buffer
	err := Export(nil, ExportJSON, &buf)
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	s := buildExporterState()
	var buf bytes.Buffer
	err := Export(s, ExportFormat("xml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_JSON(t *testing.T) {
	s := buildExporterState()
	var buf bytes.Buffer
	if err := Export(s, ExportJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rows []ExportRow
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Type != "aws_instance" {
		t.Errorf("expected first row type aws_instance, got %s", rows[0].Type)
	}
	if rows[0].ID != "i-abc123" {
		t.Errorf("expected ID i-abc123, got %s", rows[0].ID)
	}
}

func TestExport_CSV(t *testing.T) {
	s := buildExporterState()
	var buf bytes.Buffer
	if err := Export(s, ExportCSV, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 data rows
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header+2), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "type,name,id") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExport_EmptyState(t *testing.T) {
	s := NewState()
	var buf bytes.Buffer
	if err := Export(s, ExportJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []ExportRow
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty state, got %d", len(rows))
	}
}
