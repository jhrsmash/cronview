package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeCalendarEntry(job string, daysAgo int, status model.Status) model.LogEntry {
	return model.LogEntry{
		JobName:  job,
		Hostname: "host1",
		Time:     time.Now().UTC().AddDate(0, 0, -daysAgo),
		Status:   status,
	}
}

func TestRenderCalendar_NonEmpty(t *testing.T) {
	entries := []model.LogEntry{
		makeCalendarEntry("backup", 0, model.StatusSuccess),
		makeCalendarEntry("backup", 1, model.StatusFailed),
	}
	out := RenderCalendar("backup", entries, DefaultCalendarOptions())
	if out == "" {
		t.Fatal("expected non-empty calendar output")
	}
}

func TestRenderCalendar_ContainsDayLabels(t *testing.T) {
	out := RenderCalendar("backup", nil, DefaultCalendarOptions())
	for _, lbl := range []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"} {
		if !strings.Contains(out, lbl) {
			t.Errorf("expected day label %q in output", lbl)
		}
	}
}

func TestRenderCalendar_EmptyEntries_ShowsDots(t *testing.T) {
	out := RenderCalendar("backup", nil, DefaultCalendarOptions())
	if !strings.Contains(out, "·") {
		t.Error("expected dot characters for empty days")
	}
}

func TestRenderCalendar_CriticalDayHighlighted(t *testing.T) {
	// Fill today with all failures to trigger critical cell
	var entries []model.LogEntry
	for i := 0; i < 5; i++ {
		entries = append(entries, makeCalendarEntry("sync", 0, model.StatusFailed))
	}
	out := RenderCalendar("sync", entries, DefaultCalendarOptions())
	// Critical colour escape should appear
	if !strings.Contains(out, "\033[31m") {
		t.Error("expected red ANSI colour for critical failure rate")
	}
}

func TestRenderCalendar_SuccessOnlyShowsGreen(t *testing.T) {
	entries := []model.LogEntry{
		makeCalendarEntry("cleanup", 0, model.StatusSuccess),
		makeCalendarEntry("cleanup", 0, model.StatusSuccess),
	}
	out := RenderCalendar("cleanup", entries, DefaultCalendarOptions())
	if !strings.Contains(out, "\033[32m") {
		t.Error("expected green ANSI colour for healthy day")
	}
}

func TestRenderCalendar_FiltersOtherJobs(t *testing.T) {
	entries := []model.LogEntry{
		makeCalendarEntry("other", 0, model.StatusFailed),
		makeCalendarEntry("other", 0, model.StatusFailed),
		makeCalendarEntry("other", 0, model.StatusFailed),
	}
	opts := DefaultCalendarOptions()
	out := RenderCalendar("target", entries, opts)
	// No failures for "target", so no red cells expected
	if strings.Contains(out, "\033[31m") {
		t.Error("expected no red cells when job has no matching entries")
	}
}

func TestRenderCalendar_ZeroWeeks_ReturnsEmpty(t *testing.T) {
	opts := DefaultCalendarOptions()
	opts.Weeks = 0
	out := RenderCalendar("backup", nil, opts)
	if out != "" {
		t.Errorf("expected empty output for zero weeks, got %q", out)
	}
}

func TestRenderCalendar_RowCount(t *testing.T) {
	opts := DefaultCalendarOptions()
	opts.Weeks = 4
	out := RenderCalendar("backup", nil, opts)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	// header line + 4 week rows
	if len(lines) != 5 {
		t.Errorf("expected 5 lines (1 header + 4 weeks), got %d", len(lines))
	}
}
