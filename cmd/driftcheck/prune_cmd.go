package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/driftcheck/internal/tfstate"
)

// runPrune implements the `driftcheck prune` sub-command.
func runPrune(args []string) error {
	fs := flag.NewFlagSet("prune", flag.ContinueOnError)

	statePath := fs.String("state", "terraform.tfstate", "path to terraform state file")
	removeTypes := fs.String("remove-types", "", "comma-separated resource types to remove")
	removePrefix := fs.String("remove-prefix", "", "remove resources whose name starts with this prefix")
	removeOrphans := fs.Bool("remove-orphans", false, "remove resources with no attributes")
	dryRun := fs.Bool("dry-run", false, "print what would be removed without writing changes")

	if err := fs.Parse(args); err != nil {
		return err
	}

	s, err := tfstate.ParseFile(*statePath)
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	opts := tfstate.PruneOptions{
		RemoveByPrefix: *removePrefix,
		RemoveOrphans:  *removeOrphans,
	}
	if *removeTypes != "" {
		for _, t := range strings.Split(*removeTypes, ",") {
			opts.RemoveTypes = append(opts.RemoveTypes, strings.TrimSpace(t))
		}
	}

	_, result := tfstate.Prune(s, opts)
	tfstate.WritePruneReport(os.Stdout, &result)

	if *dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
		return nil
	}

	if len(result.Removed) == 0 {
		return nil
	}

	fmt.Fprintf(os.Stdout, "pruned %d resource(s) from state\n", len(result.Removed))
	return nil
}
