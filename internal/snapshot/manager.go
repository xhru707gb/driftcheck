package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles saving and loading named snapshots from a directory.
type Manager struct {
	dir string
}

// NewManager creates a Manager that stores snapshots under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot manager: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes snap to disk under the given name, timestamped.
func (m *Manager) Save(name string, snap *Snapshot) error {
	filename := m.filename(name)
	if err := snap.SaveToFile(filename); err != nil {
		return fmt.Errorf("snapshot manager: save %q: %w", name, err)
	}
	return nil
}

// Load reads the snapshot stored under name. Returns os.ErrNotExist if absent.
func (m *Manager) Load(name string) (*Snapshot, error) {
	filename := m.filename(name)
	snap, err := LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("snapshot manager: load %q: %w", name, err)
	}
	return snap, nil
}

// Exists reports whether a snapshot with the given name is stored on disk.
func (m *Manager) Exists(name string) bool {
	_, err := os.Stat(m.filename(name))
	return err == nil
}

// Delete removes the snapshot stored under name.
func (m *Manager) Delete(name string) error {
	filename := m.filename(name)
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot manager: delete %q: %w", name, err)
	}
	return nil
}

// CompareWithCurrent loads the stored snapshot named baseline and compares it
// against current, returning the diff result.
func (m *Manager) CompareWithCurrent(baseline string, current *Snapshot) ([]Diff, error) {
	old, err := m.Load(baseline)
	if err != nil {
		return nil, err
	}
	return Compare(old, current), nil
}

// ModTime returns the modification time of the stored snapshot file.
func (m *Manager) ModTime(name string) (time.Time, error) {
	info, err := os.Stat(m.filename(name))
	if err != nil {
		return time.Time{}, fmt.Errorf("snapshot manager: stat %q: %w", name, err)
	}
	return info.ModTime(), nil
}

func (m *Manager) filename(name string) string {
	return filepath.Join(m.dir, name+".snap.json")
}
