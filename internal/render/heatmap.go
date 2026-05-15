package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// HeatmapOptions controls rendering of the failure heatmap.
type HeatmapOptions struct {
	Days    int    // number of days to show (columns)
	Width   int    // total display width
	JobName string // label prefix
}

// DefaultHeatmapOptions returns sensible defaults.
func DefaultHeatmapOptions() HeatmapOptions {
	return HeatmapOptions{
		Days:  14,
		Width: 60,
	}
}

// heatCell maps a failure rate to a display character with ANSI colour.
func heatCell(rate float64) string {
	switch {
	case rate == 0:
		return "\033[32m░\033[0m" // green – no failures
	case rate < 0.25:
		return "\033[33m▒\033[0m" // yellow – low
	case rate < 0.60:
		return "\033[33m▓\033[0m" // bright yellow – moderate
	default:
		return "\033[31m█\033[0m" // red – high
	}
}

// RenderHeatmap writes a day-by-day failure heatmap for a single job.
// Each column represents one day; the cell shade reflects the failure rate
// for entries that fall within that day.
func RenderHeatmap(w io.Writer, stats model.AggregateStats, opts HeatmapOptions) {
	now := time.Now().UTC().Truncate(24 * time.Hour)

	// bucket entries by day offset (0 = today, Days-1 = oldest)
	buckets := make([][]model.LogEntry, opts.Days)
	for i := range buckets {
		buckets[i] = []model.LogEntry{}
	}

	cutoff := now.AddDate(0, 0, -(opts.Days - 1))
	for _, e := range stats.Entries {
		if e.Time.Before(cutoff) {
			continue
		}
		day := int(now.Sub(e.Time.UTC().Truncate(24 * time.Hour)).Hours() / 24)
		if day >= 0 && day < opts.Days {
			buckets[opts.Days-1-day] = append(buckets[opts.Days-1-day], e)
		}
	}

	// build cell row (oldest → newest, left → right)
	var cells strings.Builder
	for _, bucket := range buckets {
		if len(bucket) == 0 {
			cells.WriteString("\033[90m·\033[0m")
			continue
		}
		failed := 0
		for _, e := range bucket {
			if e.Status == "failure" {
				failed++
			}
		}
		cells.WriteString(heatCell(float64(failed) / float64(len(bucket))))
	}

	label := stats.JobName
	if opts.JobName != "" {
		label = opts.JobName
	}
	fmt.Fprintf(w, "%-20s %s  [%d days]\n", truncate(label, 20), cells.String(), opts.Days)
}
