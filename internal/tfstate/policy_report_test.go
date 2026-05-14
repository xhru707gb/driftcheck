package tfstate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/tfstate"
)

func TestWritePolicyReport_NilReport(t *testing.T) {
	var buf bytes.Buffer
	tfstate.WritePolicyReport(&buf, nil)
	if !strings.Contains(buf.String(), "no policy report") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestWritePolicyReport_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	report := &tfstate.PolicyReport{Checked: 3}
	tfstate.WritePolicyReport(&buf, report)
	out := buf.String()
	if !strings.Contains(out, "3 resource(s)") {
		t.Errorf("expected resource count in output: %s", out)
	}
	if !strings.Contains(out, "✓") {
		t.Errorf("expected pass symbol in output: %s", out)
	}
}

func TestWritePolicyReport_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	report := &tfstate.PolicyReport{
		Checked: 2,
		Violations: []tfstate.PolicyViolation{
			{
				ResourceKey: "aws_s3_bucket.data",
				RuleName:    "versioning-enabled",
				Description: "S3 buckets must have versioning enabled",
			},
		},
	}
	tfstate.WritePolicyReport(&buf, report)
	out := buf.String()
	if !strings.Contains(out, "✗") {
		t.Errorf("expected fail symbol in output: %s", out)
	}
	if !strings.Contains(out, "versioning-enabled") {
		t.Errorf("expected rule name in output: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket.data") {
		t.Errorf("expected resource key in output: %s", out)
	}
}
