package tfstate

import (
	"testing"
)

func buildAnnotatorState() *State {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "web",
		ID:         "i-abc123",
		Attributes: map[string]interface{}{"instance_type": "t2.micro", "tags": "env=prod"},
	})
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "data",
		ID:         "my-bucket",
		Attributes: map[string]interface{}{"bucket": "my-bucket"},
	})
	s.Add(Resource{
		Type:       "aws_vpc",
		Name:       "main",
		ID:         "",
		Attributes: map[string]interface{}{},
	})
	return s
}

func TestAnnotate_NilState(t *testing.T) {
	as := Annotate(nil)
	if as == nil {
		t.Fatal("expected non-nil AnnotatedState")
	}
	if len(as.Annotations) != 0 {
		t.Errorf("expected 0 annotations for nil state, got %d", len(as.Annotations))
	}
}

func TestAnnotate_NoIssues(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_instance",
		Name:       "clean",
		ID:         "i-clean",
		Attributes: map[string]interface{}{"instance_type": "t3.small", "tags": "env=staging"},
	})
	as := Annotate(s)
	if len(as.Annotations) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(as.Annotations))
	}
}

func TestAnnotate_MissingID(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_vpc",
		Name:       "main",
		ID:         "",
		Attributes: map[string]interface{}{"cidr": "10.0.0.0/16"},
	})
	as := Annotate(s)
	errs := as.ByLevel(AnnotationError)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error annotation, got %d", len(errs))
	}
	if errs[0].Message != "resource has no ID" {
		t.Errorf("unexpected message: %s", errs[0].Message)
	}
}

func TestAnnotate_MissingTags(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_s3_bucket",
		Name:       "logs",
		ID:         "logs-bucket",
		Attributes: map[string]interface{}{"bucket": "logs-bucket"},
	})
	as := Annotate(s)
	infos := as.ByLevel(AnnotationInfo)
	if len(infos) != 1 {
		t.Fatalf("expected 1 info annotation for missing tags, got %d", len(infos))
	}
}

func TestAnnotate_EmptyAttributes(t *testing.T) {
	s := NewState()
	s.Add(Resource{
		Type:       "aws_subnet",
		Name:       "private",
		ID:         "subnet-001",
		Attributes: map[string]interface{}{},
	})
	as := Annotate(s)
	warns := as.ByLevel(AnnotationWarning)
	if len(warns) != 1 {
		t.Fatalf("expected 1 warning for empty attributes, got %d", len(warns))
	}
}

func TestAnnotate_ByLevel_MultipleAnnotations(t *testing.T) {
	s := buildAnnotatorState()
	as := Annotate(s)

	if len(as.Annotations) == 0 {
		t.Fatal("expected some annotations")
	}

	errs := as.ByLevel(AnnotationError)
	for _, a := range errs {
		if a.Level != AnnotationError {
			t.Errorf("expected error level, got %s", a.Level)
		}
	}
}
