package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// HistogramOptions configures how the failure histogram is rendered.
type HistogramOptions struct {
	Buckets   int
	BarWidth  int
	MaxHeight int
	Window    time.Duration
}

// DefaultHistogramOptions returns sensible defaults for histogram rendering.
func DefaultHistogramOptions() HistogramOptions {
	return HistogramOptions{
		Buckets:   24,
		BarWidth:  3,
		MaxHeight: 8,
		Window:    24 * time.Hour,
	}
}

// RenderHistogram writes an ASCII bar chart of failure counts per time bucket.
func RenderHistogram(w io.Writer, entries []model.LogEntry, opts HistogramOptions) {
	if opts.Buckets <= 0 || opts.MaxHeight <= 0 {
		return
	}

	counts := bucketFailures(entries, opts.Buckets, opts.Window)

	maxVal := 0
	for _, c := range counts {
		if c > maxVal {
			maxVal = c
		}
	}

	if maxVal == 0 {
		fmt.Fprintln(w, "  [no failures in window]")
		return
	}

	for row := opts.MaxHeight; row >= 1; row-- {
		line := ""
		for _, c := range counts {
			barHeight := 0
			if maxVal > 0 {
				barHeight = (c * opts.MaxHeight) / maxVal
			}
			if barHeight >= row {
				line += strings.Repeat("█", opts.BarWidth)
			} else {
				line += strings.Repeat(" ", opts.BarWidth)
			}
			line += " "
		}
		fmt.Fprintln(w, " "+strings.TrimRight(line, " "))
	}

	// X-axis
	fmt.Fprintln(w, " "+strings.Repeat("-", opts.Buckets*(opts.BarWidth+1)))
	fmt.Fprintf(w, "  %-*s%s\n", opts.Buckets*(opts.BarWidth+1)/2, "older", "newer")
}

func bucketFailures(entries []model.LogEntry, buckets int, window time.Duration) []int {
	counts := make([]int, buckets)
	if len(entries) == 0 {
		return counts
	}

	now := time.Now()
	start := now.Add(-window)
	bucketDur := window / time.Duration(buckets)

	for _, e := range entries {
		if e.Status != model.StatusFailure {
			continue
		}
		if e.Time.Before(start) {
			continue
		}
		idx := int(e.Time.Sub(start) / bucketDur)
		if idx >= buckets {
			idx = buckets - 1
		}
		counts[idx]++
	}
	return counts
}
