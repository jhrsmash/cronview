package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// TimelineOptions controls rendering of the per-job timeline strip.
type TimelineOptions struct {
	// Width is the number of time slots to render.
	Width int
	// WindowDuration is the total time span the timeline covers.
	WindowDuration time.Duration
	// SuccessChar is the rune used for a slot with no failures.
	SuccessChar rune
	// FailureChar is the rune used for a slot with at least one failure.
	FailureChar rune
	// EmptyChar is the rune used for a slot with no activity.
	EmptyChar rune
}

// DefaultTimelineOptions returns sensible defaults for a 60-slot timeline.
func DefaultTimelineOptions() TimelineOptions {
	return TimelineOptions{
		Width:          60,
		WindowDuration: 24 * time.Hour,
		SuccessChar:    '▪',
		FailureChar:    '✗',
		EmptyChar:      '·',
	}
}

// RenderTimeline produces a single-line timeline string for the given job stats.
// Each slot represents an equal time bucket within WindowDuration.
// The rightmost slot is the most recent.
func RenderTimeline(stats model.AggregateStats, now time.Time, opts TimelineOptions) string {
	if opts.Width <= 0 {
		return ""
	}

	slotDur := opts.WindowDuration / time.Duration(opts.Width)
	if slotDur <= 0 {
		slotDur = time.Minute
	}

	type slotState int
	const (
		slotEmpty slotState = iota
		slotSuccess
		slotFailure
	)

	slots := make([]slotState, opts.Width)
	windowStart := now.Add(-opts.WindowDuration)

	for _, e := range stats.Entries {
		if e.Time.Before(windowStart) || e.Time.After(now) {
			continue
		}
		offset := e.Time.Sub(windowStart)
		idx := int(offset / slotDur)
		if idx >= opts.Width {
			idx = opts.Width - 1
		}
		if e.Status == "failure" || e.Status == "failed" {
			slots[idx] = slotFailure
		} else if slots[idx] != slotFailure {
			slots[idx] = slotSuccess
		}
	}

	var sb strings.Builder
	for _, s := range slots {
		switch s {
		case slotFailure:
			sb.WriteRune(opts.FailureChar)
		case slotSuccess:
			sb.WriteRune(opts.SuccessChar)
		default:
			sb.WriteRune(opts.EmptyChar)
		}
	}

	return fmt.Sprintf("[%s]", sb.String())
}
