package tfstate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference state used for drift comparison.
type Baseline struct {
	CreatedAt  time.Time         `json:"created_at"`
	TFVersion  string            `json:"terraform_version"`
	Resources  map[string]BaselineResource `json:"resources"`
}

// BaselineResource holds the expected attribute values for a single resource.
type BaselineResource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
}

// NewBaseline creates a Baseline snapshot from the given State.
func NewBaseline(s *State) (*Baseline, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot create baseline from nil state")
	}
	b := &Baseline{
		CreatedAt: time.Now().UTC(),
		TFVersion: s.TFVersion,
		Resources: make(map[string]BaselineResource, len(s.resources)),
	}
	for k, r := range s.resources {
		attrs := make(map[string]interface{}, len(r.Attributes))
		for ak, av := range r.Attributes {
			attrs[ak] = av
		}
		b.Resources[k.String()] = BaselineResource{
			Type:       r.Type,
			Name:       r.Name,
			Attributes: attrs,
		}
	}
	return b, nil
}

// SaveBaseline serialises a Baseline to a JSON file at path.
func SaveBaseline(b *Baseline, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create baseline file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a Baseline from a JSON file at path.
func LoadBaseline(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open baseline file: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode baseline: %w", err)
	}
	return &b, nil
}
