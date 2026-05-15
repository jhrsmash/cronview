package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
	"github.com/user/cronview/internal/parser"
)

func makeTimelineStats(entries []parser.LogEntry) model.AggregateStats {
	return model.AggregateStats{Entries: entries}
}

var timelineNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestRenderTimeline_LengthMatchesWidth(t *testing.T) {
	stats := makeTimelineStats(nil)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	// output is "[" + width chars + "]"
	expected := opts.Width + 2
	if len([]rune(out)) != expected {
		t.Errorf("expected length %d, got %d: %q", expected, len([]rune(out)), out)
	}
}

func TestRenderTimeline_EmptyEntries_AllDots(t *testing.T) {
	stats := makeTimelineStats(nil)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	inner := strings.Trim(out, "[]")
	for _, r := range inner {
		if r != opts.EmptyChar {
			t.Errorf("expected all empty chars, got %q in %q", r, out)
			break
		}
	}
}

func TestRenderTimeline_RecentFailureAppearsInLastSlot(t *testing.T) {
	recent := timelineNow.Add(-1 * time.Minute)
	entries := []parser.LogEntry{
		{Time: recent, Status: "failure", JobName: "backup"},
	}
	stats := makeTimelineStats(entries)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	inner := []rune(strings.Trim(out, "[]"))
	last := inner[len(inner)-1]
	if last != opts.FailureChar {
		t.Errorf("expected last slot to be failure char %q, got %q", opts.FailureChar, last)
	}
}

func TestRenderTimeline_SuccessEntry_ShowsSuccessChar(t *testing.T) {
	recent := timelineNow.Add(-30 * time.Minute)
	entries := []parser.LogEntry{
		{Time: recent, Status: "success", JobName: "backup"},
	}
	stats := makeTimelineStats(entries)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	if !strings.ContainsRune(out, opts.SuccessChar) {
		t.Errorf("expected success char %q in output %q", opts.SuccessChar, out)
	}
}

func TestRenderTimeline_OldEntries_Ignored(t *testing.T) {
	old := timelineNow.Add(-48 * time.Hour)
	entries := []parser.LogEntry{
		{Time: old, Status: "failure", JobName: "backup"},
	}
	stats := makeTimelineStats(entries)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	if strings.ContainsRune(out, opts.FailureChar) {
		t.Errorf("old entry should not appear in timeline: %q", out)
	}
}

func TestRenderTimeline_ZeroWidth_ReturnsEmpty(t *testing.T) {
	stats := makeTimelineStats(nil)
	opts := DefaultTimelineOptions()
	opts.Width = 0
	out := RenderTimeline(stats, timelineNow, opts)
	if out != "" {
		t.Errorf("expected empty string for zero width, got %q", out)
	}
}

func TestRenderTimeline_FailureOverridesSuccess(t *testing.T) {
	slotTime := timelineNow.Add(-30 * time.Minute)
	entries := []parser.LogEntry{
		{Time: slotTime, Status: "success", JobName: "job"},
		{Time: slotTime.Add(time.Second), Status: "failure", JobName: "job"},
	}
	stats := makeTimelineStats(entries)
	opts := DefaultTimelineOptions()
	out := RenderTimeline(stats, timelineNow, opts)
	if !strings.ContainsRune(out, opts.FailureChar) {
		t.Errorf("expected failure char when failure follows success in same slot: %q", out)
	}
}
