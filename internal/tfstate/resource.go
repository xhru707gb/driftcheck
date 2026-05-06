package tfstate

import "fmt"

// ResourceType represents the Terraform resource type (e.g. "aws_instance").
type ResourceType string

// ResourceKey uniquely identifies a resource by type and name.
type ResourceKey struct {
	Type ResourceType
	Name string
}

// String returns a human-readable representation of the key.
func (k ResourceKey) String() string {
	return fmt.Sprintf("%s.%s", k.Type, k.Name)
}

// Resource holds the parsed state of a single Terraform-managed resource.
type Resource struct {
	Key        ResourceKey
	Provider   string
	Attributes map[string]interface{}
}

// State is a collection of resources keyed by their ResourceKey.
type State struct {
	Resources map[ResourceKey]*Resource
}

// NewState initialises an empty State.
func NewState() *State {
	return &State{Resources: make(map[ResourceKey]*Resource)}
}

// Add inserts or replaces a resource in the state.
func (s *State) Add(r *Resource) {
	s.Resources[r.Key] = r
}

// Get retrieves a resource by key; returns nil if not found.
func (s *State) Get(k ResourceKey) *Resource {
	return s.Resources[k]
}

// Keys returns all resource keys present in the state.
func (s *State) Keys() []ResourceKey {
	keys := make([]ResourceKey, 0, len(s.Resources))
	for k := range s.Resources {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of resources in the state.
func (s *State) Len() int {
	return len(s.Resources)
}
