package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot holds a point-in-time capture of live cloud resource attributes.
type Snapshot struct {
	CapturedAt time.Time                       `json:"captured_at"`
	Resources  map[string]map[string]string    `json:"resources"`
}

// New creates an empty Snapshot with the current timestamp.
func New() *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Resources:  make(map[string]map[string]string),
	}
}

// Add stores the attributes for a resource key.
func (s *Snapshot) Add(key string, attrs map[string]string) {
	copy := make(map[string]string, len(attrs))
	for k, v := range attrs {
		copy[k] = v
	}
	s.Resources[key] = copy
}

// Get returns the attributes for a resource key and whether it was found.
func (s *Snapshot) Get(key string) (map[string]string, bool) {
	attrs, ok := s.Resources[key]
	return attrs, ok
}

// SaveToFile serialises the snapshot as JSON to the given path.
func (s *Snapshot) SaveToFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadFromFile deserialises a snapshot from a JSON file.
func LoadFromFile(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}
