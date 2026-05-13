package tfstate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/driftcheck/internal/tfstate"
)

func TestWriteDiffReport_NoDiff(t *testing.T) {
	s := buildDiffState(map[string]map[string]interface{}{"web": {"ami": "ami-1"}})
	diff := tfstate.DiffStates(s, s)

	var buf bytes.Buffer
	if err := tfstate.WriteDiffReport(&buf, diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No state differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteDiffReport_Added(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{})
	target := buildDiffState(map[string]map[string]interface{}{
		"web": {"ami": "ami-123", "type": "t2.micro"},
	})
	diff := tfstate.DiffStates(base, target)

	var buf bytes.Buffer
	if err := tfstate.WriteDiffReport(&buf, diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got: %s", out)
	}
	if !strings.Contains(out, "Summary:") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestWriteDiffReport_Removed(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{
		"db": {"engine": "mysql"},
	})
	target := buildDiffState(map[string]map[string]interface{}{})
	diff := tfstate.DiffStates(base, target)

	var buf bytes.Buffer
	if err := tfstate.WriteDiffReport(&buf, diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "removed") {
		t.Errorf("expected 'removed' in output, got: %s", buf.String())
	}
}

func TestWriteDiffReport_Modified(t *testing.T) {
	base := buildDiffState(map[string]map[string]interface{}{
		"svc": {"port": "80", "proto": "http"},
	})
	target := buildDiffState(map[string]map[string]interface{}{
		"svc": {"port": "443", "proto": "https"},
	})
	diff := tfstate.DiffStates(base, target)

	var buf bytes.Buffer
	if err := tfstate.WriteDiffReport(&buf, diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "modified") {
		t.Errorf("expected 'modified' in output, got: %s", out)
	}
	if !strings.Contains(out, "->") {
		t.Errorf("expected attribute change arrow in output, got: %s", out)
	}
}
