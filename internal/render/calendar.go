package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// CalendarOptions controls the appearance of the calendar heatmap view.
type CalendarOptions struct {
	Weeks      int
	CellWidth  int
	Critical   float64
	Degraded   float64
}

// DefaultCalendarOptions returns sensible defaults for the calendar view.
func DefaultCalendarOptions() CalendarOptions {
	return CalendarOptions{
		Weeks:     8,
		CellWidth: 2,
		Critical:  0.5,
		Degraded:  0.2,
	}
}

// RenderCalendar renders a week-by-week calendar heatmap for a single job,
// showing failure density per day over the past N weeks.
func RenderCalendar(job string, entries []model.LogEntry, opts CalendarOptions) string {
	if opts.Weeks <= 0 || opts.CellWidth <= 0 {
		return ""
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	total := opts.Weeks * 7

	// bucket[i] = (failures, runs) for day i (0 = oldest)
	type bucket struct{ failures, runs int }
	buckets := make([]bucket, total)

	for _, e := range entries {
		if job != "" && e.JobName != job {
			continue
		}
		day := int(now.Sub(e.Time.UTC().Truncate(24 * time.Hour)).Hours() / 24)
		idx := total - 1 - day
		if idx < 0 || idx >= total {
			continue
		}
		buckets[idx].runs++
		if e.Status == model.StatusFailed {
			buckets[idx].failures++
		}
	}

	var sb strings.Builder
	cell := strings.Repeat(" ", opts.CellWidth)

	// Day-of-week header (Mon–Sun)
	dayLabels := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	sb.WriteString("    ")
	for _, lbl := range dayLabels {
		sb.WriteString(fmt.Sprintf("%-*s", opts.CellWidth, lbl))
	}
	sb.WriteByte('\n')

	for w := 0; w < opts.Weeks; w++ {
		weekStart := now.AddDate(0, 0, -(opts.Weeks-1-w)*7)
		sb.WriteString(fmt.Sprintf("%-4s", weekStart.Format("01/02")))
		for d := 0; d < 7; d++ {
			idx := w*7 + d
			b := buckets[idx]
			var ch string
			switch {
			case b.runs == 0:
				ch = "·"
			case float64(b.failures)/float64(b.runs) >= opts.Critical:
				ch = "\033[31m█\033[0m"
			case float64(b.failures)/float64(b.runs) >= opts.Degraded:
				ch = "\033[33m▒\033[0m"
			default:
				ch = "\033[32m░\033[0m"
			}
			sb.WriteString(ch + cell[1:])
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
