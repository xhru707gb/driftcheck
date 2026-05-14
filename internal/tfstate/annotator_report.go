package tfstate

import (
	"fmt"
	"io"
	"sort"
)

// WriteAnnotationReport writes a human-readable annotation report to w.
// Annotations are grouped by severity level (error > warning > info) and
// sorted alphabetically by resource key within each group.
func WriteAnnotationReport(w io.Writer, as *AnnotatedState) {
	if as == nil || len(as.Annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return
	}

	// Group by level for ordered output: error > warning > info
	levels := []AnnotationLevel{AnnotationError, AnnotationWarning, AnnotationInfo}
	levelLabel := map[AnnotationLevel]string{
		AnnotationError:   "ERROR",
		AnnotationWarning: "WARN ",
		AnnotationInfo:    "INFO ",
	}

	total := len(as.Annotations)
	fmt.Fprintf(w, "Annotation Report (%d total)\n", total)
	fmt.Fprintln(w, "================================")

	for _, level := range levels {
		group := as.ByLevel(level)
		if len(group) == 0 {
			continue
		}

		// Sort within group by resource key for deterministic output.
		sort.Slice(group, func(i, j int) bool {
			return group[i].Key.String() < group[j].Key.String()
		})

		for _, a := range group {
			fmt.Fprintf(w, "[%s] %s — %s\n", levelLabel[level], a.Key.String(), a.Message)
		}
	}

	fmt.Fprintln(w, "================================")
	errorCount := len(as.ByLevel(AnnotationError))
	warnCount := len(as.ByLevel(AnnotationWarning))
	infoCount := len(as.ByLevel(AnnotationInfo))
	fmt.Fprintf(w, "Summary: %d error(s), %d warning(s), %d info(s)\n", errorCount, warnCount, infoCount)
}

// HasErrors reports whether the annotated state contains any error-level annotations.
// This is useful for callers that need to determine whether to treat the result
// as a failure without re-examining the full annotation list.
func HasErrors(as *AnnotatedState) bool {
	if as == nil {
		return false
	}
	return len(as.ByLevel(AnnotationError)) > 0
}
