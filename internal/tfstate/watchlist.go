package tfstate

import (
	"fmt"
	"strings"
)

// WatchEntry represents a single resource being watched for drift.
type WatchEntry struct {
	ResourceType string
	ResourceName string
	Attributes   []string // specific attributes to watch; empty means all
}

// Watchlist holds a set of resources to monitor for drift.
type Watchlist struct {
	entries []WatchEntry
}

// NewWatchlist creates an empty Watchlist.
func NewWatchlist() *Watchlist {
	return &Watchlist{}
}

// Add appends a WatchEntry to the watchlist.
func (w *Watchlist) Add(entry WatchEntry) {
	w.entries = append(w.entries, entry)
}

// Entries returns all watch entries.
func (w *Watchlist) Entries() []WatchEntry {
	return w.entries
}

// Matches reports whether the given resource key is covered by this watchlist.
// If the watchlist is empty, all resources match.
func (w *Watchlist) Matches(key ResourceKey) bool {
	if len(w.entries) == 0 {
		return true
	}
	for _, e := range w.entries {
		if e.ResourceType == key.Type && e.ResourceName == key.Name {
			return true
		}
	}
	return false
}

// WatchedAttributes returns the attributes to watch for a resource key.
// Returns nil (meaning all) when the entry has no attribute restrictions.
func (w *Watchlist) WatchedAttributes(key ResourceKey) []string {
	for _, e := range w.entries {
		if e.ResourceType == key.Type && e.ResourceName == key.Name {
			return e.Attributes
		}
	}
	return nil
}

// String returns a human-readable summary of the watchlist.
func (w *Watchlist) String() string {
	if len(w.entries) == 0 {
		return "watchlist: (all resources)"
	}
	lines := make([]string, 0, len(w.entries))
	for _, e := range w.entries {
		attrs := "all attributes"
		if len(e.Attributes) > 0 {
			attrs = strings.Join(e.Attributes, ", ")
		}
		lines = append(lines, fmt.Sprintf("  %s.%s [%s]", e.ResourceType, e.ResourceName, attrs))
	}
	return "watchlist:\n" + strings.Join(lines, "\n")
}
