package render

import (
	"fmt"
	"strings"

	"github.com/user/cronview/internal/model"
)

// GaugeOptions controls the appearance of the failure-rate gauge.
type GaugeOptions struct {
	Width     int
	WarnRate  float64
	CritRate  float64
	ShowLabel bool
}

// DefaultGaugeOptions returns sensible defaults for RenderGauge.
func DefaultGaugeOptions() GaugeOptions {
	return GaugeOptions{
		Width:     30,
		WarnRate:  0.25,
		CritRate:  0.50,
		ShowLabel: true,
	}
}

// RenderGauge renders a single-line ASCII gauge showing the failure rate
// for the given job stats. Returns an empty string when stats is nil.
func RenderGauge(stats *model.AggregateStats, opts GaugeOptions) string {
	if stats == nil {
		return ""
	}

	width := opts.Width
	if width < 4 {
		width = 4
	}

	rate := stats.FailureRate
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}

	filled := int(rate * float64(width))
	empty := width - filled

	var fillChar, emptyChar string
	switch {
	case rate >= opts.CritRate:
		fillChar = "█"
	case rate >= opts.WarnRate:
		fillChar = "▓"
	default:
		fillChar = "░"
	}
	emptyChar = "·"

	bar := strings.Repeat(fillChar, filled) + strings.Repeat(emptyChar, empty)

	if opts.ShowLabel {
		return fmt.Sprintf("[%s] %.1f%%  %s", bar, rate*100, stats.JobName)
	}
	return fmt.Sprintf("[%s] %.1f%%", bar, rate*100)
}

// RenderGaugeList renders a vertical list of gauges for multiple jobs,
// sorted by the order provided. Returns at most maxRows rows.
func RenderGaugeList(statsList []model.AggregateStats, opts GaugeOptions, maxRows int) string {
	if len(statsList) == 0 {
		return "(no jobs)"
	}

	var sb strings.Builder
	count := 0
	for _, s := range statsList {
		if maxRows > 0 && count >= maxRows {
			break
		}
		copy := s
		sb.WriteString(RenderGauge(&copy, opts))
		sb.WriteByte('\n')
		count++
	}
	return strings.TrimRight(sb.String(), "\n")
}
