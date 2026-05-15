package render

import (
	"fmt"
	"strings"

	"github.com/user/cronview/internal/model"
)

// TrendDirection indicates whether a job's failure rate is improving, worsening, or stable.
type TrendDirection int

const (
	TrendStable   TrendDirection = 0
	TrendImproving TrendDirection = 1
	TrendWorsening TrendDirection = -1
)

// TrendOptions configures the trend renderer.
type TrendOptions struct {
	WindowSize  int     // number of recent entries to compare against baseline
	DeltaThresh float64 // minimum change to be considered a trend
}

// DefaultTrendOptions returns sensible defaults.
func DefaultTrendOptions() TrendOptions {
	return TrendOptions{
		WindowSize:  10,
		DeltaThresh: 0.05,
	}
}

// ComputeTrend returns the trend direction for a job based on recent vs overall failure rate.
func ComputeTrend(stats model.AggregateStats, opts TrendOptions) TrendDirection {
	entries := stats.Entries
	if len(entries) < opts.WindowSize*2 {
		return TrendStable
	}

	recent := entries[len(entries)-opts.WindowSize:]
	baseline := entries[:len(entries)-opts.WindowSize]

	recentRate := failureRate(recent)
	baselineRate := failureRate(baseline)

	delta := recentRate - baselineRate
	if delta > opts.DeltaThresh {
		return TrendWorsening
	}
	if delta < -opts.DeltaThresh {
		return TrendImproving
	}
	return TrendStable
}

// RenderTrendBadge returns a short colored string badge for a trend direction.
func RenderTrendBadge(dir TrendDirection) string {
	switch dir {
	case TrendImproving:
		return "\033[32m▼ improving\033[0m"
	case TrendWorsening:
		return "\033[31m▲ worsening\033[0m"
	default:
		return "\033[33m● stable\033[0m"
	}
}

// RenderTrendTable renders a summary table of job trends.
func RenderTrendTable(statsList []model.AggregateStats, opts TrendOptions) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-30s %-12s %s\n", "Job", "Fail Rate", "Trend")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 58))
	for _, s := range statsList {
		dir := ComputeTrend(s, opts)
		fmt.Fprintf(&sb, "%-30s %-12s %s\n",
			truncate(s.JobName, 30),
			fmt.Sprintf("%.1f%%", s.FailureRate*100),
			RenderTrendBadge(dir),
		)
	}
	return sb.String()
}

func failureRate(entries []model.LogEntry) float64 {
	if len(entries) == 0 {
		return 0
	}
	failed := 0
	for _, e := range entries {
		if e.Status == model.StatusFailure {
			failed++
		}
	}
	return float64(failed) / float64(len(entries))
}
