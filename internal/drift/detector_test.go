package drift_test

import (
	"testing"

	"github.com/your-org/driftcheck/internal/drift"
	"github.com/your-org/driftcheck/internal/tfstate"
)

func makeState(resources []tfstate.Resource) *tfstate.State {
	return &tfstate.State{Version: 4, Resources: resources}
}

func TestDetector_NoChanges(t *testing.T) {
	state := makeState([]tfstate.Resource{
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"instance_type": "t2.micro"}},
	})
	live := map[string]map[string]interface{}{
		"aws_instance.web": {"instance_type": "t2.micro"},
	}
	findings := drift.New().Compare(state, live)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestDetector_MissingResource(t *testing.T) {
	state := makeState([]tfstate.Resource{
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{}},
	})
	findings := drift.New().Compare(state, map[string]map[string]interface{}{})
	if len(findings) != 1 || findings[0].Kind != drift.KindMissing {
		t.Errorf("expected 1 MISSING finding, got %v", findings)
	}
}

func TestDetector_ModifiedAttribute(t *testing.T) {
	state := makeState([]tfstate.Resource{
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"instance_type": "t2.micro"}},
	})
	live := map[string]map[string]interface{}{
		"aws_instance.web": {"instance_type": "t3.small"},
	}
	findings := drift.New().Compare(state, live)
	if len(findings) != 1 || findings[0].Kind != drift.KindModified {
		t.Errorf("expected 1 MODIFIED finding, got %v", findings)
	}
	if findings[0].Attribute != "instance_type" {
		t.Errorf("unexpected attribute: %s", findings[0].Attribute)
	}
}

func TestDetector_ExtraResource(t *testing.T) {
	state := makeState(nil)
	live := map[string]map[string]interface{}{
		"aws_s3_bucket.logs": {"bucket": "my-logs"},
	}
	findings := drift.New().Compare(state, live)
	if len(findings) != 1 || findings[0].Kind != drift.KindExtra {
		t.Errorf("expected 1 EXTRA finding, got %v", findings)
	}
}

func TestFindingString(t *testing.T) {
	f := drift.Finding{Kind: drift.KindModified, ResourceKey: "aws_instance.web", Attribute: "ami", Expected: "ami-old", Actual: "ami-new"}
	if f.String() == "" {
		t.Error("expected non-empty string representation")
	}
}
