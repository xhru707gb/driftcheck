package tfstate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftcheck/internal/tfstate"
)

func buildBaselineState(t *testing.T) *tfstate.State {
	t.Helper()
	s := tfstate.NewState()
	s.TFVersion = "1.5.0"
	s.Add(tfstate.Resource{
		Type: "aws_instance", Name: "web",
		Attributes: map[string]interface{}{"instance_type": "t3.micro", "ami": "ami-123"},
	})
	s.Add(tfstate.Resource{
		Type: "aws_s3_bucket", Name: "assets",
		Attributes: map[string]interface{}{"bucket": "my-assets"},
	})
	return s
}

func TestNewBaseline_NilState(t *testing.T) {
	_, err := tfstate.NewBaseline(nil)
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestNewBaseline_ResourceCount(t *testing.T) {
	s := buildBaselineState(t)
	b, err := tfstate.NewBaseline(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(b.Resources); got != 2 {
		t.Errorf("want 2 resources, got %d", got)
	}
}

func TestNewBaseline_AttributesCopied(t *testing.T) {
	s := buildBaselineState(t)
	b, _ := tfstate.NewBaseline(s)
	key := "aws_instance.web"
	br, ok := b.Resources[key]
	if !ok {
		t.Fatalf("expected key %q in baseline", key)
	}
	if br.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("unexpected attribute value: %v", br.Attributes["instance_type"])
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	s := buildBaselineState(t)
	b, _ := tfstate.NewBaseline(s)

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := tfstate.SaveBaseline(b, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := tfstate.LoadBaseline(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Resources) != len(b.Resources) {
		t.Errorf("resource count mismatch: want %d got %d", len(b.Resources), len(loaded.Resources))
	}
	if loaded.TFVersion != b.TFVersion {
		t.Errorf("TFVersion mismatch: want %q got %q", b.TFVersion, loaded.TFVersion)
	}
}

func TestLoadBaseline_NotFound(t *testing.T) {
	_, err := tfstate.LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveBaseline_BadPath(t *testing.T) {
	s := buildBaselineState(t)
	b, _ := tfstate.NewBaseline(s)
	err := tfstate.SaveBaseline(b, filepath.Join(os.DevNull, "sub", "file.json"))
	if err == nil {
		t.Fatal("expected error for unwritable path")
	}
}
