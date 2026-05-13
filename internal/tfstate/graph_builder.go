package tfstate

import "strings"

// BuildGraph constructs a dependency Graph by inspecting resource attributes
// for references to other resources in the state.
//
// It looks for string attribute values of the form "<type>.<name>.*" and adds
// an edge from the current resource to the referenced resource when both exist
// in the state.
func BuildGraph(s *State) (*Graph, []error) {
	g := NewGraph(s)
	var errs []error

	keys := s.Keys()
	// Build a quick lookup: "type.name" -> ResourceKey
	lookup := make(map[string]ResourceKey, len(keys))
	for _, k := range keys {
		lookup[k.Type+"."+k.Name] = k
	}

	for _, k := range keys {
		res, ok := s.Get(k)
		if !ok {
			continue
		}
		for _, v := range res.Attributes {
			str, ok := v.(string)
			if !ok {
				continue
			}
			ref := extractRef(str)
			if ref == "" {
				continue
			}
			if dep, found := lookup[ref]; found && dep != k {
				if err := g.AddEdge(k, dep); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return g, errs
}

// extractRef attempts to parse a Terraform-style reference like
// "aws_vpc.main.id" and return the "type.name" portion.
func extractRef(val string) string {
	parts := strings.SplitN(val, ".", 3)
	if len(parts) < 2 {
		return ""
	}
	// Heuristic: type segment contains an underscore (e.g. aws_vpc)
	if !strings.Contains(parts[0], "_") {
		return ""
	}
	return parts[0] + "." + parts[1]
}
