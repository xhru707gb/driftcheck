package tfstate

import "sort"

// Group represents a collection of resources sharing a common key.
type Group struct {
	Key       string
	Resources []Resource
}

// GroupBy defines the field to group resources by.
type GroupBy string

const (
	GroupByType   GroupBy = "type"
	GroupByModule GroupBy = "module"
	GroupByRegion GroupBy = "region"
)

// GroupResult holds all groups produced by GroupResources.
type GroupResult struct {
	Groups []Group
	Total  int
}

// GroupResources partitions all resources in s by the given GroupBy field.
// Resources that do not have the requested attribute fall into the "" bucket.
func GroupResources(s *State, by GroupBy) (*GroupResult, error) {
	if s == nil {
		return &GroupResult{}, nil
	}

	buckets := make(map[string][]Resource)

	for _, key := range s.Keys() {
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		var bucket string
		switch by {
		case GroupByType:
			bucket = key.Type
		case GroupByModule:
			if v, ok := res.Attributes["module"]; ok {
				bucket = v
			}
		case GroupByRegion:
			if v, ok := res.Attributes["region"]; ok {
				bucket = v
			}
		}

		buckets[bucket] = append(buckets[bucket], res)
	}

	result := &GroupResult{Total: 0}
	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		group := Group{Key: k, Resources: buckets[k]}
		result.Groups = append(result.Groups, group)
		result.Total += len(group.Resources)
	}

	return result, nil
}
