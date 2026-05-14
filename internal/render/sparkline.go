package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// sparklineBlocks are unicode block characters ordered from empty to full.
var sparklineBlocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// SparklineOptions controls how a sparkline is rendered.
type SparklineOptions struct {
	// Buckets is the number of time buckets (columns) in the sparkline.
	Buckets int
	// Window is the total time range to cover.
	Window time.Duration
}

// DefaultSparklineOptions returns sensible defaults: 14 daily buckets.
func DefaultSparklineOptions() SparklineOptions {
	return SparklineOptions{
		Buckets: 14,
		Window:  14 * 24 * time.Hour,
	}
}

// RenderSparkline builds a sparkline string showing failure counts per time
// bucket for the given job entries. The most recent bucket is on the right.
func RenderSparkline(entries []model.LogEntry, opts SparklineOptions) string {
	if opts.Buckets <= 0 {
		return ""
	}

	now := time.Now()
	bucketDur := opts.Window / time.Duration(opts.Buckets)
	counts := make([]int, opts.Buckets)

	for _, e := range entries {
		if e.Status != model.StatusFailure {
			continue
		}
		age := now.Sub(e.Timestamp)
		if age < 0 || age >= opts.Window {
			continue
		}
		idx := int(age / bucketDur)
		// age==0 maps to bucket opts.Buckets-1 (rightmost = most recent)
		revIdx := opts.Buckets - 1 - idx
		if revIdx >= 0 && revIdx < opts.Buckets {
			counts[revIdx]++
		}
	}

	max := 0
	for _, c := range counts {
		if c > max {
			max = c
		}
	}

	var sb strings.Builder
	for _, c := range counts {
		if max == 0 {
			sb.WriteRune(sparklineBlocks[0])
			continue
		}
		blockIdx := int(float64(c) / float64(max) * float64(len(sparklineBlocks)-1))
		sb.WriteRune(sparklineBlocks[blockIdx])
	}

	return fmt.Sprintf("[%s]", sb.String())
}
