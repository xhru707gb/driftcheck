package tfstate

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single Terraform-managed resource.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// State holds parsed Terraform state data.
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// rawState mirrors the on-disk terraform.tfstate JSON structure.
type rawState struct {
	Version   int          `json:"version"`
	Resources []rawResource `json:"resources"`
}

type rawResource struct {
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	Instances []rawInstance `json:"instances"`
}

type rawInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ParseFile reads and parses a terraform.tfstate file from disk.
func ParseFile(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}
	return Parse(data)
}

// Parse decodes raw JSON bytes into a State.
func Parse(data []byte) (*State, error) {
	var raw rawState
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unmarshalling state: %w", err)
	}

	state := &State{Version: raw.Version}
	for _, r := range raw.Resources {
		attrs := map[string]interface{}{}
		if len(r.Instances) > 0 {
			attrs = r.Instances[0].Attributes
		}
		state.Resources = append(state.Resources, Resource{
			Type:       r.Type,
			Name:       r.Name,
			Provider:   r.Provider,
			Attributes: attrs,
		})
	}
	return state, nil
}

// ResourceMap returns resources indexed by "type.name" for quick lookup.
func (s *State) ResourceMap() map[string]Resource {
	m := make(map[string]Resource, len(s.Resources))
	for _, r := range s.Resources {
		key := r.Type + "." + r.Name
		m[key] = r
	}
	return m
}
