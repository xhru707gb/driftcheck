package tfstate_test

import (
	"testing"

	"github.com/your-org/driftcheck/internal/tfstate"
)

const sampleState = `{
  "version": 4,
  "resources": [
    {
      "type": "aws_instance",
      "name": "web",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "ami": "ami-0c55b159cbfafe1f0",
            "instance_type": "t2.micro"
          }
        }
      ]
    }
  ]
}`

func TestParse_ValidState(t *testing.T) {
	state, err := tfstate.Parse([]byte(sampleState))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	r := state.Resources[0]
	if r.Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %s", r.Type)
	}
	if r.Attributes["instance_type"] != "t2.micro" {
		t.Errorf("unexpected instance_type: %v", r.Attributes["instance_type"])
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := tfstate.Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestResourceMap(t *testing.T) {
	state, err := tfstate.Parse([]byte(sampleState))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := state.ResourceMap()
	if _, ok := m["aws_instance.web"]; !ok {
		t.Error("expected key aws_instance.web in resource map")
	}
}

func TestParse_NoInstances(t *testing.T) {
	raw := `{"version":4,"resources":[{"type":"aws_s3_bucket","name":"data","provider":"p","instances":[]}]}`
	state, err := tfstate.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources[0].Attributes) != 0 {
		t.Error("expected empty attributes for resource with no instances")
	}
}
