package tfstate

import "fmt"

// Edge represents a dependency between two resources.
type Edge struct {
	From ResourceKey
	To   ResourceKey
}

// Graph holds a directed dependency graph of Terraform resources.
type Graph struct {
	nodes map[ResourceKey]bool
	edges []Edge
	adj   map[ResourceKey][]ResourceKey
}

// NewGraph constructs a dependency graph from a State.
func NewGraph(s *State) *Graph {
	g := &Graph{
		nodes: make(map[ResourceKey]bool),
		adj:   make(map[ResourceKey][]ResourceKey),
	}
	for _, k := range s.Keys() {
		g.nodes[k] = true
	}
	return g
}

// AddEdge records a dependency: `from` depends on `to`.
func (g *Graph) AddEdge(from, to ResourceKey) error {
	if !g.nodes[from] {
		return fmt.Errorf("unknown resource: %s", from)
	}
	if !g.nodes[to] {
		return fmt.Errorf("unknown resource: %s", to)
	}
	g.edges = append(g.edges, Edge{From: from, To: to})
	g.adj[from] = append(g.adj[from], to)
	return nil
}

// Dependents returns all resources that directly depend on the given key.
func (g *Graph) Dependents(key ResourceKey) []ResourceKey {
	var result []ResourceKey
	for _, e := range g.edges {
		if e.To == key {
			result = append(result, e.From)
		}
	}
	return result
}

// Dependencies returns all direct dependencies of the given key.
func (g *Graph) Dependencies(key ResourceKey) []ResourceKey {
	return g.adj[key]
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	return len(g.edges)
}

// HasCycle reports whether the graph contains a cycle using DFS.
func (g *Graph) HasCycle() bool {
	visited := make(map[ResourceKey]bool)
	onStack := make(map[ResourceKey]bool)
	for k := range g.nodes {
		if !visited[k] {
			if g.dfs(k, visited, onStack) {
				return true
			}
		}
	}
	return false
}

func (g *Graph) dfs(k ResourceKey, visited, onStack map[ResourceKey]bool) bool {
	visited[k] = true
	onStack[k] = true
	for _, dep := range g.adj[k] {
		if !visited[dep] && g.dfs(dep, visited, onStack) {
			return true
		} else if onStack[dep] {
			return true
		}
	}
	onStack[k] = false
	return false
}
