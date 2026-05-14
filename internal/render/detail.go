package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// DetailOptions controls rendering of the job detail view.
type DetailOptions struct {
	Width       int
	TimeFormat  string
	ShowHistory int // number of recent entries to show
}

// DefaultDetailOptions returns sensible defaults for the detail panel.
func DefaultDetailOptions() DetailOptions {
	return DetailOptions{
		Width:       80,
		TimeFormat:  time.RFC3339,
		ShowHistory: 10,
	}
}

// RenderDetail writes a detailed view of a single job's stats and recent
// log entries to w.
func RenderDetail(w io.Writer, stats model.JobStats, opts DetailOptions) {
	if opts.Width <= 0 {
		opts.Width = 80
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}
	if opts.ShowHistory <= 0 {
		opts.ShowHistory = 10
	}

	sep := strings.Repeat("─", opts.Width)

	fmt.Fprintf(w, "%s\n", sep)
	fmt.Fprintf(w, "  Job: %s\n", stats.JobName)
	fmt.Fprintf(w, "  Host: %s\n", stats.Hostname)
	fmt.Fprintf(w, "%s\n", sep)

	fmt.Fprintf(w, "  Total Runs   : %d\n", stats.TotalRuns)
	fmt.Fprintf(w, "  Failures     : %d\n", stats.Failures)
	fmt.Fprintf(w, "  Failure Rate : %.1f%%\n", stats.FailureRate*100)
	fmt.Fprintf(w, "  Last Status  : %s\n", stats.LastStatus)
	if !stats.LastRun.IsZero() {
		fmt.Fprintf(w, "  Last Run     : %s\n", stats.LastRun.Format(opts.TimeFormat))
	}

	fmt.Fprintf(w, "%s\n", sep)
	fmt.Fprintf(w, "  Recent Entries (last %d):\n", opts.ShowHistory)

	entries := stats.Entries
	if len(entries) > opts.ShowHistory {
		entries = entries[len(entries)-opts.ShowHistory:]
	}
	for _, e := range entries {
		fmt.Fprintf(w, "    [%s] %-8s  %s\n",
			e.Timestamp.Format(opts.TimeFormat),
			e.Status,
			truncate(e.Message, opts.Width-30),
		)
	}

	fmt.Fprintf(w, "%s\n", sep)
}
