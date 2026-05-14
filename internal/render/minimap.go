package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronview/internal/model"
)

// MinimapOptions controls rendering of the job health minimap.
type MinimapOptions struct {
	// Width is the number of columns in the minimap grid.
	Width int
	// ShowLegend controls whether a legend is appended below the grid.
	ShowLegend bool
}

// DefaultMinimapOptions returns sensible defaults for the minimap.
func DefaultMinimapOptions() MinimapOptions {
	return MinimapOptions{
		Width:      10,
		ShowLegend: true,
	}
}

// RenderMinimap writes a compact grid showing the health of each job.
// Each cell represents one job: '█' for healthy, '▒' for degraded (≥ warn
// threshold), '░' for critical (≥ crit threshold).
func RenderMinimap(w io.Writer, stats []model.JobStats, opts MinimapOptions) {
	if len(stats) == 0 {
		fmt.Fprintln(w, "(no jobs)")
		return
	}

	const (
		cellHealthy  = "█"
		cellDegraded = "▒"
		cellCritical = "░"

		warnThreshold = 0.10
		critThreshold = 0.30
	)

	width := opts.Width
	if width <= 0 {
		width = DefaultMinimapOptions().Width
	}

	var sb strings.Builder
	for i, s := range stats {
		if i > 0 && i%width == 0 {
			sb.WriteString("\n")
		}
		switch {
		case s.FailureRate >= critThreshold:
			sb.WriteString(cellCritical)
		case s.FailureRate >= warnThreshold:
			sb.WriteString(cellDegraded)
		default:
			sb.WriteString(cellHealthy)
		}
	}
	fmt.Fprintln(w, sb.String())

	if opts.ShowLegend {
		fmt.Fprintf(w, "%s healthy  %s degraded(≥10%%)  %s critical(≥30%%)\n",
			cellHealthy, cellDegraded, cellCritical)
	}
}
