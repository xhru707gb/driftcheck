package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

// TestPruneAndReport_FullScenario exercises Prune followed by WritePruneReport
// end-to-end to ensure the two components integrate correctly.
func TestPruneAndReport_FullScenario(t *testing.T) {
	s := NewState()
	s.Add(Resource{Type: "aws_instance", Name: "api", ID: "i-10", Attributes: map[string]string{"ami": "ami-1"}})
	s.Add(Resource{Type: "aws_instance", Name: "batch", ID: "i-11", Attributes: map[string]string{"ami": "ami-2"}})
	s.Add(Resource{Type: "aws_s3_bucket", Name: "logs", ID: "bkt-2", Attributes: map[string]string{}})
	s.Add(Resource{Type: "aws_rds_instance", Name: "db", ID: "db-1", Attributes: map[string]string{"engine": "postgres"}})

	opts := PruneOptions{
		RemoveTypes:   []string{"aws_instance"},
		RemoveOrphans: true,
	}

	out, res := Prune(s, opts)

	// aws_instance x2 + aws_s3_bucket (orphan) = 3 removed
	if len(res.Removed) != 3 {
		t.Fatalf("expected 3 removed, got %d", len(res.Removed))
	}
	if res.Kept != 1 {
		t.Fatalf("expected 1 kept, got %d", res.Kept)
	}

	keys := out.Keys()
	if len(keys) != 1 {
		t.Fatalf("output state should have 1 resource, got %d", len(keys))
	}

	remaining, _ := out.Get(keys[0])
	if remaining.Type != "aws_rds_instance" {
		t.Errorf("expected aws_rds_instance to survive, got %s", remaining.Type)
	}

	var buf bytes.Buffer
	WritePruneReport(&buf, &res)
	report := buf.String()

	if !strings.Contains(report, "removed 3") {
		t.Errorf("report should mention 3 removals, got: %s", report)
	}
	if !strings.Contains(report, "kept 1") {
		t.Errorf("report should mention 1 kept, got: %s", report)
	}
}
