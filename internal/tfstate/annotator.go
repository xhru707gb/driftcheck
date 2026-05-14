package tfstate

import "fmt"

// AnnotationLevel represents the severity of an annotation.
type AnnotationLevel string

const (
	AnnotationInfo    AnnotationLevel = "info"
	AnnotationWarning AnnotationLevel = "warning"
	AnnotationError   AnnotationLevel = "error"
)

// Annotation holds a message attached to a specific resource.
type Annotation struct {
	Key     ResourceKey
	Level   AnnotationLevel
	Message string
}

// AnnotatedState pairs a State with annotations produced during analysis.
type AnnotatedState struct {
	State       *State
	Annotations []Annotation
}

// Annotate inspects each resource in state and attaches annotations for
// common issues such as missing tags, empty attributes, or deprecated types.
func Annotate(s *State) *AnnotatedState {
	as := &AnnotatedState{State: s}
	if s == nil {
		return as
	}

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		if res.ID == "" {
			as.add(key, AnnotationError, "resource has no ID")
		}

		if len(res.Attributes) == 0 {
			as.add(key, AnnotationWarning, "resource has no attributes")
		}

		if _, hasTag := res.Attributes["tags"]; !hasTag {
			if isTaggable(key.Type) {
				as.add(key, AnnotationInfo, fmt.Sprintf("resource type %q has no tags attribute", key.Type))
			}
		}
	}

	return as
}

// ByLevel returns all annotations matching the given level.
func (as *AnnotatedState) ByLevel(level AnnotationLevel) []Annotation {
	var out []Annotation
	for _, a := range as.Annotations {
		if a.Level == level {
			out = append(out, a)
		}
	}
	return out
}

func (as *AnnotatedState) add(key ResourceKey, level AnnotationLevel, msg string) {
	as.Annotations = append(as.Annotations, Annotation{Key: key, Level: level, Message: msg})
}

func isTaggable(resourceType string) bool {
	taggable := map[string]bool{
		"aws_instance":        true,
		"aws_s3_bucket":       true,
		"aws_security_group":  true,
		"aws_vpc":             true,
		"aws_subnet":          true,
		"aws_lambda_function": true,
	}
	return taggable[resourceType]
}
