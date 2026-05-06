package cloud_test

import (
	"context"
	"testing"

	"github.com/example/driftcheck/internal/cloud"
	"github.com/example/driftcheck/internal/tfstate"
)

func makeState(resType, id string, attrs map[string]interface{}) *tfstate.State {
	return &tfstate.State{
		Resources: []tfstate.Resource{
			{
				Type: resType,
				Instances: []tfstate.Instance{
					{Attributes: attrs},
				},
			},
		},
	}
}

func TestReconciler_NoDrift(t *testing.T) {
	f := &mockFetcher{
		resources: map[string]*cloud.LiveResource{
			"aws_instance/i-abc": {
				Type:       "aws_instance",
				ID:         "i-abc",
				Attributes: cloud.ResourceAttributes{"instance_type": "t3.micro"},
			},
		},
	}
	state := makeState("aws_instance", "i-abc", map[string]interface{}{
		"id":            "i-abc",
		"instance_type": "t3.micro",
	})
	r := cloud.NewReconciler(f)
	res, err := r.Reconcile(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Drifts) != 0 {
		t.Errorf("expected no drifts, got %d", len(res.Drifts))
	}
}

func TestReconciler_DriftDetected(t *testing.T) {
	f := &mockFetcher{
		resources: map[string]*cloud.LiveResource{
			"aws_instance/i-abc": {
				Type:       "aws_instance",
				ID:         "i-abc",
				Attributes: cloud.ResourceAttributes{"instance_type": "t3.large"},
			},
		},
	}
	state := makeState("aws_instance", "i-abc", map[string]interface{}{
		"id":            "i-abc",
		"instance_type": "t3.micro",
	})
	r := cloud.NewReconciler(f)
	res, err := r.Reconcile(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Drifts) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(res.Drifts))
	}
	if res.Drifts[0].Attribute != "instance_type" {
		t.Errorf("unexpected attribute: %s", res.Drifts[0].Attribute)
	}
}

func TestReconciler_MissingResource(t *testing.T) {
	f := &mockFetcher{resources: map[string]*cloud.LiveResource{}}
	state := makeState("aws_instance", "i-gone", map[string]interface{}{"id": "i-gone"})
	r := cloud.NewReconciler(f)
	res, err := r.Reconcile(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "i-gone" {
		t.Errorf("expected missing i-gone, got %v", res.Missing)
	}
}
