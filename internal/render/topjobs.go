package render

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/cronview/internal/model"
)

// TopJobsOptions controls how the top-jobs panel is rendered.
type TopJobsOptions struct {
	MaxRows   int
	SortBy    string // "failures" or "rate"
	Width     int
	ShowRank  bool
}

// DefaultTopJobsOptions returns sensible defaults.
func DefaultTopJobsOptions() TopJobsOptions {
	return TopJobsOptions{
		MaxRows:  10,
		SortBy:   "failures",
		Width:    60,
		ShowRank: true,
	}
}

// RenderTopJobs writes a ranked list of jobs ordered by failure count or
// failure rate to w. Jobs with zero failures are omitted.
func RenderTopJobs(w io.Writer, stats []model.AggregateStats, opts TopJobsOptions) {
	if opts.MaxRows <= 0 {
		opts.MaxRows = DefaultTopJobsOptions().MaxRows
	}

	// Filter to only jobs that have at least one failure.
	var failing []model.AggregateStats
	for _, s := range stats {
		if s.FailureCount > 0 {
			failing = append(failing, s)
		}
	}

	if len(failing) == 0 {
		fmt.Fprintln(w, "  No failures recorded.")
		return
	}

	// Sort by chosen criterion (descending).
	switch opts.SortBy {
	case "rate":
		sort.Slice(failing, func(i, j int) bool {
			return failing[i].FailureRate > failing[j].FailureRate
		})
	default: // "failures"
		sort.Slice(failing, func(i, j int) bool {
			if failing[i].FailureCount != failing[j].FailureCount {
				return failing[i].FailureCount > failing[j].FailureCount
			}
			return failing[i].JobName < failing[j].JobName
		})
	}

	if len(failing) > opts.MaxRows {
		failing = failing[:opts.MaxRows]
	}

	nameWidth := 30
	if opts.Width > 50 {
		nameWidth = opts.Width - 28
	}

	header := fmt.Sprintf("  %-*s  %6s  %7s  %6s",
		nameWidth, "JOB", "RUNS", "FAIL", "RATE")
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, "  "+strings.Repeat("-", len(header)-2))

	for i, s := range failing {
		rank := ""
		if opts.ShowRank {
			rank = fmt.Sprintf("%2d. ", i+1)
			nameWidth -= 4
			if nameWidth < 8 {
				nameWidth = 8
			}
		}
		name := truncate(s.JobName, nameWidth)
		fmt.Fprintf(w, "  %s%-*s  %6d  %7d  %5.1f%%\n",
			rank, nameWidth, name, s.TotalRuns, s.FailureCount, s.FailureRate*100)
		if opts.ShowRank {
			nameWidth += 4 // restore for next iteration
		}
	}
}
