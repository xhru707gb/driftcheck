package tfstate

import (
	"testing"
)

func buildGraphState() *State {
	s := NewState()
	s.Add(Resource{Key: ResourceKey{Type: "aws_vpc", Name: "main"}, ID: "vpc-1", Attributes: map[string]interface{}{}})
	s.Add(Resource{Key: ResourceKey{Type: "aws_subnet", Name: "pub"}, ID: "sub-1", Attributes: map[string]interface{}{}})
	s.Add(Resource{Key: ResourceKey{Type: "aws_instance", Name: "web"}, ID: "i-1", Attributes: map[string]interface{}{}})
	return s
}

func TestNewGraph_NodeCount(t *testing.T) {
	g := NewGraph(buildGraphState())
	if g.NodeCount() != 3 {
		t.Fatalf("expected 3 nodes, got %d", g.NodeCount())
	}
}

func TestAddEdge_Valid(t *testing.T) {
	g := NewGraph(buildGraphState())
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	sub := ResourceKey{Type: "aws_subnet", Name: "pub"}
	if err := g.AddEdge(sub, vpc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.EdgeCount() != 1 {
		t.Fatalf("expected 1 edge, got %d", g.EdgeCount())
	}
}

func TestAddEdge_UnknownResource(t *testing.T) {
	g := NewGraph(buildGraphState())
	unknown := ResourceKey{Type: "aws_s3_bucket", Name: "logs"}
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	if err := g.AddEdge(unknown, vpc); err == nil {
		t.Fatal("expected error for unknown from-resource")
	}
}

func TestDependencies(t *testing.T) {
	g := NewGraph(buildGraphState())
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	sub := ResourceKey{Type: "aws_subnet", Name: "pub"}
	inst := ResourceKey{Type: "aws_instance", Name: "web"}
	_ = g.AddEdge(sub, vpc)
	_ = g.AddEdge(inst, sub)

	deps := g.Dependencies(inst)
	if len(deps) != 1 || deps[0] != sub {
		t.Fatalf("unexpected dependencies: %v", deps)
	}
}

func TestDependents(t *testing.T) {
	g := NewGraph(buildGraphState())
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	sub := ResourceKey{Type: "aws_subnet", Name: "pub"}
	_ = g.AddEdge(sub, vpc)

	dependents := g.Dependents(vpc)
	if len(dependents) != 1 || dependents[0] != sub {
		t.Fatalf("unexpected dependents: %v", dependents)
	}
}

func TestHasCycle_False(t *testing.T) {
	g := NewGraph(buildGraphState())
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	sub := ResourceKey{Type: "aws_subnet", Name: "pub"}
	_ = g.AddEdge(sub, vpc)
	if g.HasCycle() {
		t.Fatal("expected no cycle")
	}
}

func TestHasCycle_True(t *testing.T) {
	g := NewGraph(buildGraphState())
	vpc := ResourceKey{Type: "aws_vpc", Name: "main"}
	sub := ResourceKey{Type: "aws_subnet", Name: "pub"}
	_ = g.AddEdge(sub, vpc)
	_ = g.AddEdge(vpc, sub) // creates cycle
	if !g.HasCycle() {
		t.Fatal("expected cycle to be detected")
	}
}
