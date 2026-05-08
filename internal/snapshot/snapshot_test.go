package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftcheck/internal/snapshot"
)

func TestNew_EmptySnapshot(t *testing.T) {
	s := snapshot.New()
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(s.Resources) != 0 {
		t.Errorf("expected empty resources, got %d", len(s.Resources))
	}
}

func TestAdd_And_Get(t *testing.T) {
	s := snapshot.New()
	attrs := map[string]string{"instance_type": "t3.micro", "region": "us-east-1"}
	s.Add("aws_instance.web", attrs)

	got, ok := s.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected resource to be found")
	}
	if got["instance_type"] != "t3.micro" {
		t.Errorf("unexpected instance_type: %s", got["instance_type"])
	}
}

func TestGet_Missing(t *testing.T) {
	s := snapshot.New()
	_, ok := s.Get("aws_instance.missing")
	if ok {
		t.Error("expected miss, got hit")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := snapshot.New()
	s.Add("aws_s3_bucket.logs", map[string]string{"acl": "private"})

	if err := s.SaveToFile(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := snapshot.LoadFromFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	attrs, ok := loaded.Get("aws_s3_bucket.logs")
	if !ok {
		t.Fatal("resource not found after round-trip")
	}
	if attrs["acl"] != "private" {
		t.Errorf("expected acl=private, got %s", attrs["acl"])
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := snapshot.LoadFromFile("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := snapshot.LoadFromFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
