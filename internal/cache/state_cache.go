package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a cached cloud resource state with metadata.
type Entry struct {
	ResourceID string                 `json:"resource_id"`
	Attributes map[string]interface{} `json:"attributes"`
	FetchedAt  time.Time              `json:"fetched_at"`
}

// StateCache persists fetched cloud resource state to disk to avoid redundant API calls.
type StateCache struct {
	dir string
	TTL time.Duration
}

// New creates a StateCache that stores entries under dir with the given TTL.
func New(dir string, ttl time.Duration) (*StateCache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &StateCache{dir: dir, TTL: ttl}, nil
}

func (c *StateCache) path(resourceID string) string {
	safe := filepath.Base(resourceID)
	return filepath.Join(c.dir, safe+".json")
}

// Get returns a cached Entry if it exists and has not expired.
func (c *StateCache) Get(resourceID string) (*Entry, bool) {
	data, err := os.ReadFile(c.path(resourceID))
	if err != nil {
		return nil, false
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, false
	}
	if time.Since(e.FetchedAt) > c.TTL {
		return nil, false
	}
	return &e, true
}

// Set writes an Entry to the cache.
func (c *StateCache) Set(e *Entry) error {
	e.FetchedAt = time.Now()
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(c.path(e.ResourceID), data, 0o644)
}

// Invalidate removes a cached entry for the given resource ID.
func (c *StateCache) Invalidate(resourceID string) error {
	err := os.Remove(c.path(resourceID))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
