package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronview/internal/snapshot"
)

// SnapshotViewOptions controls rendering of snapshot comparison output.
type SnapshotViewOptions struct {
	Width    int
	MaxRows  int
}

// DefaultSnapshotViewOptions returns sensible defaults.
func DefaultSnapshotViewOptions() SnapshotViewOptions {
	return SnapshotViewOptions{Width: 80, MaxRows: 20}
}

// RenderSnapshotDiff writes a human-readable diff table between two snapshots.
func RenderSnapshotDiff(w io.Writer, deltas []snapshot.Delta, opts SnapshotViewOptions) {
	if len(deltas) == 0 {
		fmt.Fprintln(w, "No differences between snapshots.")
		return
	}

	header := fmt.Sprintf("%-30s %8s %8s %10s %6s",
		"JOB", "OLD %", "NEW %", "DELTA", "NOTE")
	sep := strings.Repeat("-", opts.Width)

	fmt.Fprintln(w, header)
	fmt.Fprintln(w, sep)

	limit := len(deltas)
	if opts.MaxRows > 0 && opts.MaxRows < limit {
		limit = opts.MaxRows
	}

	for _, d := range deltas[:limit] {
		note := ""
		switch {
		case d.OnlyInNew:
			note = "[new]"
		case d.OnlyInOld:
			note = "[gone]"
		case d.RateDelta > 0.05:
			note = "▲ worse"
		case d.RateDelta < -0.05:
			note = "▼ better"
		}
		fmt.Fprintf(w, "%-30s %7.1f%% %7.1f%% %+9.1f%% %6s\n",
			truncate(d.JobName, 30),
			d.OldRate*100,
			d.NewRate*100,
			d.RateDelta*100,
			note,
		)
	}

	if len(deltas) > limit {
		fmt.Fprintf(w, "... and %d more\n", len(deltas)-limit)
	}
}
