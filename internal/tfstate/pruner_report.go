package tfstate

import (
	"fmt"
	"io"
)

// WritePruneReport writes a human-readable summary of a PruneResult to w.
func WritePruneReport(w io.Writer, r *PruneResult) {
	if r == nil {
		fmt.Fprintln(w, "no prune result")
		return
	}

	if len(r.Removed) == 0 {
		fmt.Fprintln(w, "pruner: nothing removed")
		return
	}

	fmt.Fprintf(w, "pruner: removed %d resource(s), kept %d\n", len(r.Removed), r.Kept)
	for _, key := range r.Removed {
		fmt.Fprintf(w, "  - %s\n", key)
	}
}
