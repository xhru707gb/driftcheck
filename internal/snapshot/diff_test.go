package snapshot_test

import (
	"testing"

	"github.com/example/driftcheck/internal/snapshot"
)

func makeSnap(resources map[string]map[string]string) *snapshot.Snapshot {
	s := snapshot.New()
	for k, attrs := range resources {
		s.Add(k, attrs)
	}
	return s
}

func TestCompare_NoDiff(t *testing.T) {
	old := makeSnap(map[string]map[string]string{
		"aws_instance.web": {"type": "t3.micro"},
	})
	new := makeSnap(map[string]map[string]string{
		"aws_instance.web": {"type": "t3.micro"},
	})
	diffs := snapshot.Compare(old, new)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompare_ModifiedAttribute(t *testing.T) {
	old := makeSnap(map[string]map[string]string{
		"aws_instance.web": {"type": "t3.micro"},
	})
	new := makeSnap(map[string]map[string]string{
		"aws_instance.web": {"type": "t3.large"},
	})
	diffs := snapshot.Compare(old, new)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != snapshot.DiffModified {
		t.Errorf("expected modified, got %s", diffs[0].Kind)
	}
	if len(diffs[0].Attributes) != 1 {
		t.Errorf("expected 1 attribute change, got %d", len(diffs[0].Attributes))
	}
}

func TestCompare_AddedResource(t *testing.T) {
	old := makeSnap(map[string]map[string]string{})
	new := makeSnap(map[string]map[string]string{
		"aws_s3_bucket.data": {"acl": "private"},
	})
	diffs := snapshot.Compare(old, new)
	if len(diffs) != 1 || diffs[0].Kind != snapshot.DiffAdded {
		t.Errorf("expected 1 added diff, got %+v", diffs)
	}
}

func TestCompare_RemovedResource(t *testing.T) {
	old := makeSnap(map[string]map[string]string{
		"aws_s3_bucket.data": {"acl": "private"},
	})
	new := makeSnap(map[string]map[string]string{})
	diffs := snapshot.Compare(old, new)
	if len(diffs) != 1 || diffs[0].Kind != snapshot.DiffRemoved {
		t.Errorf("expected 1 removed diff, got %+v", diffs)
	}
}

func TestCompare_MultipleChanges(t *testing.T) {
	old := makeSnap(map[string]map[string]string{
		"aws_instance.a": {"type": "t3.micro"},
		"aws_instance.b": {"type": "t3.small"},
	})
	new := makeSnap(map[string]map[string]string{
		"aws_instance.a": {"type": "t3.large"},
		"aws_instance.c": {"type": "t3.nano"},
	})
	diffs := snapshot.Compare(old, new)
	if len(diffs) != 3 {
		t.Errorf("expected 3 diffs (1 modified, 1 removed, 1 added), got %d", len(diffs))
	}
}
