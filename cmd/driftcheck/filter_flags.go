package main

import (
	"flag"
	"strings"

	"github.com/example/driftcheck/internal/tfstate"
)

// multiFlag is a flag.Value that accumulates comma-separated or repeated values.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			*m = append(*m, part)
		}
	}
	return nil
}

// filterFlags holds the CLI flags that map to tfstate.FilterOptions.
type filterFlags struct {
	types        multiFlag
	excludeTypes multiFlag
	namePrefix   string
}

// register binds the filter flags to the provided FlagSet.
func (f *filterFlags) register(fs *flag.FlagSet) {
	fs.Var(&f.types, "type", "include only this resource type (repeatable, comma-separated)")
	fs.Var(&f.excludeTypes, "exclude-type", "exclude this resource type (repeatable, comma-separated)")
	fs.StringVar(&f.namePrefix, "name-prefix", "", "include only resources whose name starts with this prefix")
}

// toFilterOptions converts the parsed flags into a tfstate.FilterOptions value.
func (f *filterFlags) toFilterOptions() tfstate.FilterOptions {
	return tfstate.FilterOptions{
		Types:        []string(f.types),
		ExcludeTypes: []string(f.excludeTypes),
		NamePrefix:   f.namePrefix,
	}
}
