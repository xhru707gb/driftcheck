package snapshot

import (
	"os"
	"testing"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func buildSnap(t *testing.T, resources map[string]map[string]string) *Snapshot {
	t.Helper()
	s := New()
	for id, attrs := range resources {
		s.Add(id, attrs)
	}
	return s
}

func TestManager_SaveAndLoad(t *testing.T) {
	m := newTestManager(t)
	snap := buildSnap(t, map[string]map[string]string{
		"aws_instance.web": {"instance_type": "t3.micro"},
	})
	if err := m.Save("baseline", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := m.Load("baseline")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	attrs, ok := loaded.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected resource not found in loaded snapshot")
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("got instance_type=%q, want t3.micro", attrs["instance_type"])
	}
}

func TestManager_Exists(t *testing.T) {
	m := newTestManager(t)
	if m.Exists("nope") {
		t.Fatal("expected Exists to return false for missing snapshot")
	}
	snap := buildSnap(t, nil)
	_ = m.Save("present", snap)
	if !m.Exists("present") {
		t.Fatal("expected Exists to return true after Save")
	}
}

func TestManager_Delete(t *testing.T) {
	m := newTestManager(t)
	snap := buildSnap(t, nil)
	_ = m.Save("todelete", snap)
	if err := m.Delete("todelete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if m.Exists("todelete") {
		t.Fatal("expected snapshot to be gone after Delete")
	}
	// Deleting non-existent should not error
	if err := m.Delete("ghost"); err != nil {
		t.Fatalf("Delete non-existent: %v", err)
	}
}

func TestManager_Load_NotFound(t *testing.T) {
	m := newTestManager(t)
	_, err := m.Load("missing")
	if err == nil {
		t.Fatal("expected error loading missing snapshot")
	}
}

func TestManager_CompareWithCurrent(t *testing.T) {
	m := newTestManager(t)
	old := buildSnap(t, map[string]map[string]string{
		"aws_instance.web": {"instance_type": "t3.micro"},
	})
	_ = m.Save("baseline", old)

	current := buildSnap(t, map[string]map[string]string{
		"aws_instance.web": {"instance_type": "t3.large"},
	})
	diffs, err := m.CompareWithCurrent("baseline", current)
	if err != nil {
		t.Fatalf("CompareWithCurrent: %v", err)
	}
	if len(diffs) == 0 {
		t.Fatal("expected diffs, got none")
	}
}

func TestManager_ModTime(t *testing.T) {
	m := newTestManager(t)
	snap := buildSnap(t, nil)
	_ = m.Save("ts", snap)
	mod, err := m.ModTime("ts")
	if err != nil {
		t.Fatalf("ModTime: %v", err)
	}
	if mod.IsZero() {
		t.Fatal("expected non-zero ModTime")
	}
	_, err = m.ModTime("absent")
	if !os.IsNotExist(err) && err == nil {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}
