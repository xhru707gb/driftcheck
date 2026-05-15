package tfstate

import (
	"testing"
)

func buildLabelerState() *State {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		ID:         "i-abc123",
		Attributes: map[string]interface{}{"region": "us-east-1", "env": "prod"},
	})
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "data",
		ID:         "my-bucket",
		Attributes: map[string]interface{}{},
	})
	return s
}

func TestApplyLabels_NilState(t *testing.T) {
	_, err := ApplyLabels(nil, []LabelRule{})
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestApplyLabels_NoRules(t *testing.T) {
	s := buildLabelerState()
	report, err := ApplyLabels(s, []LabelRule{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Total != 2 {
		t.Errorf("expected Total=2, got %d", report.Total)
	}
	if report.Labeled != 0 {
		t.Errorf("expected Labeled=0, got %d", report.Labeled)
	}
}

func TestApplyLabels_FromAttr(t *testing.T) {
	s := buildLabelerState()
	rules := []LabelRule{
		{Key: "computed_region", FromAttr: "region", Prefix: "region:"},
	}
	report, err := ApplyLabels(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Labeled != 1 {
		t.Errorf("expected Labeled=1, got %d", report.Labeled)
	}
	res, _ := s.Get(ResourceKey{Type: "aws_instance", Name: "web"})
	if res.Attributes["computed_region"] != "region:us-east-1" {
		t.Errorf("unexpected label value: %v", res.Attributes["computed_region"])
	}
}

func TestApplyLabels_MissingAttr_Skipped(t *testing.T) {
	s := buildLabelerState()
	rules := []LabelRule{
		{Key: "computed_env", FromAttr: "env"},
	}
	report, err := ApplyLabels(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// aws_s3_bucket/data has no "env" attr, should be skipped
	for _, r := range report.Results {
		if r.ResourceKey == "aws_s3_bucket.data" {
			if len(r.Skipped) != 1 || r.Skipped[0] != "computed_env" {
				t.Errorf("expected skipped label for s3 bucket, got %+v", r)
			}
		}
	}
}

func TestApplyLabels_EmptyRuleKey_Error(t *testing.T) {
	s := buildLabelerState()
	rules := []LabelRule{
		{Key: "", FromAttr: "region"},
	}
	report, err := ApplyLabels(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range report.Results {
		if len(r.Errors) == 0 {
			t.Errorf("expected error for empty rule key in result %s", r.ResourceKey)
		}
	}
}

func TestApplyLabels_FromID(t *testing.T) {
	s := buildLabelerState()
	rules := []LabelRule{
		{Key: "id_label", Prefix: "id:", Suffix: "-v1"},
	}
	_, err := ApplyLabels(s, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, _ := s.Get(ResourceKey{Type: "aws_instance", Name: "web"})
	expected := "id:i-abc123-v1"
	if res.Attributes["id_label"] != expected {
		t.Errorf("expected %q, got %v", expected, res.Attributes["id_label"])
	}
}
