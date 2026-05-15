package tfstate

import (
	"fmt"
	"io"
	"strings"
)

// WriteWatchlistReport writes a human-readable summary of the watchlist to w.
func WriteWatchlistReport(w io.Writer, wl *Watchlist) {
	if wl == nil {
		fmt.Fprintln(w, "watchlist: nil")
		return
	}
	entries := wl.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "Watchlist: monitoring ALL resources")
		return
	}

	fmt.Fprintf(w, "Watchlist: %d resource(s) monitored\n", len(entries))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, e := range entries {
		attrs := "<all attributes>"
		if len(e.Attributes) > 0 {
			attrs = strings.Join(e.Attributes, ", ")
		}
		fmt.Fprintf(w, "  %-20s  attrs: %s\n",
			e.ResourceType+"."+e.ResourceName, attrs)
	}
}
