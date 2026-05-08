package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestState() *State {
	s := NewState()
	s.Add(Resource{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-1"}})
	s.Add(Resource{Type: "aws_instance", Name: "api", Attributes: map[string]interface{}{"id": "i-2"}})
	s.Add(Resource{Type: "aws_s3_bucket", Name: "assets", Attributes: map[string]interface{}{"id": "b-1"}})
	s.Add(Resource{Type: "aws_security_group", Name: "web_sg", Attributes: map[string]interface{}{"id": "sg-1"}})
	return s
}

func TestFilter_NoOptions(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{})
	assert.Equal(t, len(s.Keys()), len(out.Keys()))
}

func TestFilter_ByType(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{Types: []string{"aws_instance"}})
	require.Equal(t, 2, len(out.Keys()))
	for _, k := range out.Keys() {
		res, _ := out.Get(k)
		assert.Equal(t, "aws_instance", res.Type)
	}
}

func TestFilter_ExcludeType(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{ExcludeTypes: []string{"aws_s3_bucket"}})
	assert.Equal(t, 3, len(out.Keys()))
	for _, k := range out.Keys() {
		res, _ := out.Get(k)
		assert.NotEqual(t, "aws_s3_bucket", res.Type)
	}
}

func TestFilter_ByNamePrefix(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{NamePrefix: "web"})
	require.Equal(t, 2, len(out.Keys()))
}

func TestFilter_TypeAndPrefix(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{
		Types:      []string{"aws_instance"},
		NamePrefix: "api",
	})
	require.Equal(t, 1, len(out.Keys()))
	res, ok := out.Get(ResourceKey{Type: "aws_instance", Name: "api"})
	require.True(t, ok)
	assert.Equal(t, "api", res.Name)
}

func TestFilter_ExcludeOverridesInclude(t *testing.T) {
	s := buildTestState()
	out := Filter(s, FilterOptions{
		Types:        []string{"aws_instance"},
		ExcludeTypes: []string{"aws_instance"},
	})
	assert.Equal(t, 0, len(out.Keys()))
}
