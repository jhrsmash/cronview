package render

import (
	"fmt"
	"strings"

	"github.com/user/cronview/internal/model"
)

// SparkbarOptions controls rendering of the inline sparkbar widget.
type SparkbarOptions struct {
	Width      int
	ShowLabel  bool
	ShowRate   bool
	FillChar   rune
	EmptyChar  rune
	CritColor  string
	WarnColor  string
	OKColor    string
	ResetColor string
}

// DefaultSparkbarOptions returns sensible defaults.
func DefaultSparkbarOptions() SparkbarOptions {
	return SparkbarOptions{
		Width:      20,
		ShowLabel:  true,
		ShowRate:   true,
		FillChar:   '█',
		EmptyChar:  '░',
		CritColor:  "\033[31m",
		WarnColor:  "\033[33m",
		OKColor:    "\033[32m",
		ResetColor: "\033[0m",
	}
}

// RenderSparkbar renders a single horizontal bar representing the failure rate
// of a job, with optional label and percentage annotation.
func RenderSparkbar(stat *model.AggregateStats, opts SparkbarOptions) string {
	if stat == nil {
		return ""
	}

	rate := stat.FailureRate
	filled := int(rate * float64(opts.Width))
	if filled > opts.Width {
		filled = opts.Width
	}
	empty := opts.Width - filled

	color := opts.OKColor
	switch {
	case rate >= 0.5:
		color = opts.CritColor
	case rate >= 0.2:
		color = opts.WarnColor
	}

	bar := color +
		strings.Repeat(string(opts.FillChar), filled) +
		strings.Repeat(string(opts.EmptyChar), empty) +
		opts.ResetColor

	var sb strings.Builder
	if opts.ShowLabel {
		sb.WriteString(fmt.Sprintf("%-24s ", truncate(stat.JobName, 24)))
	}
	sb.WriteString("[")
	sb.WriteString(bar)
	sb.WriteString("]")
	if opts.ShowRate {
		sb.WriteString(fmt.Sprintf(" %5.1f%%", rate*100))
	}
	return sb.String()
}

// RenderSparkbarList renders a sparkbar for each stat, one per line.
func RenderSparkbarList(stats []*model.AggregateStats, opts SparkbarOptions) string {
	if len(stats) == 0 {
		return "(no jobs)"
	}
	lines := make([]string, 0, len(stats))
	for _, s := range stats {
		lines = append(lines, RenderSparkbar(s, opts))
	}
	return strings.Join(lines, "\n")
}
