package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

func buildGrouperState() *State {
	s := NewState()
	s.Add(Resource{ID: "1", Type: "aws_instance", Name: "web", Attributes: map[string]string{"region": "us-east-1", "module": "app"}})
	s.Add(Resource{ID: "2", Type: "aws_instance", Name: "db", Attributes: map[string]string{"region": "us-west-2", "module": "data"}})
	s.Add(Resource{ID: "3", Type: "aws_s3_bucket", Name: "assets", Attributes: map[string]string{"region": "us-east-1", "module": "app"}})
	s.Add(Resource{ID: "4", Type: "aws_s3_bucket", Name: "logs", Attributes: map[string]string{"region": "eu-west-1"}})
	return s
}

func TestGroupResources_NilState(t *testing.T) {
	result, err := GroupResources(nil, GroupByType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("expected 0 total, got %d", result.Total)
	}
}

func TestGroupResources_ByType(t *testing.T) {
	s := buildGrouperState()
	result, err := GroupResources(s, GroupByType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 4 {
		t.Errorf("expected total 4, got %d", result.Total)
	}
	if len(result.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(result.Groups))
	}
}

func TestGroupResources_ByRegion(t *testing.T) {
	s := buildGrouperState()
	result, _ := GroupResources(s, GroupByRegion)
	regions := make(map[string]int)
	for _, g := range result.Groups {
		regions[g.Key] = len(g.Resources)
	}
	if regions["us-east-1"] != 2 {
		t.Errorf("expected 2 resources in us-east-1, got %d", regions["us-east-1"])
	}
	if regions["eu-west-1"] != 1 {
		t.Errorf("expected 1 resource in eu-west-1, got %d", regions["eu-west-1"])
	}
}

func TestGroupResources_ByModule_UnsetFallback(t *testing.T) {
	s := buildGrouperState()
	result, _ := GroupResources(s, GroupByModule)
	for _, g := range result.Groups {
		if g.Key == "" && len(g.Resources) != 1 {
			t.Errorf("expected 1 resource with unset module, got %d", len(g.Resources))
		}
	}
}

func TestGroupResources_Sorted(t *testing.T) {
	s := buildGrouperState()
	result, _ := GroupResources(s, GroupByType)
	if result.Groups[0].Key > result.Groups[1].Key {
		t.Errorf("groups not sorted: %q > %q", result.Groups[0].Key, result.Groups[1].Key)
	}
}

func TestWriteGroupReport_Output(t *testing.T) {
	s := buildGrouperState()
	result, _ := GroupResources(s, GroupByType)
	var buf bytes.Buffer
	WriteGroupReport(&buf, result, GroupByType)
	out := buf.String()
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected aws_instance in output")
	}
	if !strings.Contains(out, "Total resources: 4") {
		t.Errorf("expected total line in output")
	}
}

func TestWriteGroupReport_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteGroupReport(&buf, nil, GroupByType)
	if !strings.Contains(buf.String(), "no group result") {
		t.Errorf("expected fallback message for nil result")
	}
}
