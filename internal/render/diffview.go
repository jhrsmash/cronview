package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronview/internal/model"
)

// DiffViewOptions controls how the job stat diff is rendered.
type DiffViewOptions struct {
	Width     int
	ShowDelta bool
}

// DefaultDiffViewOptions returns sensible defaults for diff rendering.
func DefaultDiffViewOptions() DiffViewOptions {
	return DiffViewOptions{
		Width:     80,
		ShowDelta: true,
	}
}

// RenderDiff writes a side-by-side comparison of two snapshots of job stats.
// It highlights jobs whose failure rate has changed between the two snapshots.
func RenderDiff(w io.Writer, before, after []model.AggregateStats, opts DiffViewOptions) {
	beforeMap := make(map[string]model.AggregateStats, len(before))
	for _, s := range before {
		beforeMap[s.JobName] = s
	}

	fmt.Fprintf(w, "%-30s %10s %10s %10s\n", "JOB", "BEFORE", "AFTER", "DELTA")
	fmt.Fprintln(w, strings.Repeat("-", opts.Width))

	for _, a := range after {
		b, exists := beforeMap[a.JobName]

		var beforeRate float64
		if exists {
			beforeRate = b.FailureRate
		}

		delta := a.FailureRate - beforeRate
		deltaStr := formatDelta(delta, opts.ShowDelta)
		marker := diffMarker(delta, exists)

		fmt.Fprintf(w, "%s %-28s %9.1f%% %9.1f%% %s\n",
			marker,
			truncate(a.JobName, 28),
			beforeRate*100,
			a.FailureRate*100,
			deltaStr,
		)
	}
}

func formatDelta(delta float64, show bool) string {
	if !show {
		return ""
	}
	switch {
	case delta > 0.001:
		return fmt.Sprintf("+%.1f%%", delta*100)
	case delta < -0.001:
		return fmt.Sprintf("%.1f%%", delta*100)
	default:
		return "  ±0.0%"
	}
}

func diffMarker(delta float64, existed bool) string {
	if !existed {
		return "N"
	}
	switch {
	case delta > 0.001:
		return "↑"
	case delta < -0.001:
		return "↓"
	default:
		return " "
	}
}
