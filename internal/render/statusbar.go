package render

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// StatusBarOptions controls the appearance of the status bar.
type StatusBarOptions struct {
	Width       int
	ShowTime    bool
	ShowFilter  bool
}

// DefaultStatusBarOptions returns sensible defaults for the status bar.
func DefaultStatusBarOptions() StatusBarOptions {
	return StatusBarOptions{
		Width:      80,
		ShowTime:   true,
		ShowFilter: true,
	}
}

// StatusBarData holds the runtime state to display in the status bar.
type StatusBarData struct {
	TotalJobs    int
	FilteredJobs int
	ActiveFilter string
	LastRefresh  time.Time
	Page         int
	TotalPages   int
}

// RenderStatusBar writes a single-line status bar summarising current view state.
func RenderStatusBar(w io.Writer, data StatusBarData, opts StatusBarOptions) {
	var parts []string

	if opts.ShowFilter && data.ActiveFilter != "" {
		parts = append(parts, fmt.Sprintf("filter:%s", data.ActiveFilter))
	}

	parts = append(parts, fmt.Sprintf("jobs:%d/%d", data.FilteredJobs, data.TotalJobs))

	if data.TotalPages > 1 {
		parts = append(parts, fmt.Sprintf("page:%d/%d", data.Page, data.TotalPages))
	}

	if opts.ShowTime && !data.LastRefresh.IsZero() {
		parts = append(parts, fmt.Sprintf("refreshed:%s", data.LastRefresh.Format("15:04:05")))
	}

	line := strings.Join(parts, "  ")
	if len(line) > opts.Width {
		line = line[:opts.Width]
	}

	padding := opts.Width - len(line)
	if padding < 0 {
		padding = 0
	}

	fmt.Fprintf(w, "[ %s%s ]\n", line, strings.Repeat(" ", padding))
}
