package tfstate

import (
	"fmt"
	"io"
	"strings"
)

var planSymbols = map[PlanAction]string{
	PlanCreate:  "+",
	PlanUpdate:  "~",
	PlanDestroy: "-",
	PlanNoOp:    " ",
}

// WritePlanReport writes a human-readable plan report to the provided writer.
func WritePlanReport(w io.Writer, plan *Plan) error {
	if plan == nil {
		return fmt.Errorf("plan must not be nil")
	}

	if !plan.HasChanges() {
		fmt.Fprintln(w, "No changes. Infrastructure matches the desired state.")
		return nil
	}

	fmt.Fprintln(w, "Terraform will perform the following actions:")
	fmt.Fprintln(w, "")

	for _, entry := range plan.Entries {
		if entry.Action == PlanNoOp {
			continue
		}
		sym := planSymbols[entry.Action]
		fmt.Fprintf(w, "  %s %s.%s\n", sym, entry.Key.Type, entry.Key.Name)
		if entry.Action == PlanUpdate && len(entry.ChangedKeys) > 0 {
			fmt.Fprintf(w, "      ~ attributes: [%s]\n", strings.Join(entry.ChangedKeys, ", "))
		}
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, plan.Summary())
	return nil
}
