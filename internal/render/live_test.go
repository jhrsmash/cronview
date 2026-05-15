package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

// makeLiveStats returns a slice of AggregateStats for live rendering tests.
func makeLiveStats(n int) []model.AggregateStats {
	now := time.Now()
	stats := make([]model.AggregateStats, n)
	for i := 0; i < n; i++ {
		stats[i] = model.AggregateStats{
			JobName:     fmt.Sprintf("job-%d", i+1),
			Hostname:    "host1",
			TotalRuns:   10,
			Failures:    i % 3,
			FailureRate: float64(i%3) / 10.0,
			LastStatus:  "SUCCESS",
			LastRun:     now.Add(-time.Duration(i) * time.Minute),
		}
	}
	return stats
}

func TestRenderLiveFrame_NonEmpty(t *testing.T) {
	stats := makeLiveStats(3)
	opts := DefaultLiveOptions()
	out := RenderLiveFrame(stats, opts)
	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty output from RenderLiveFrame")
	}
}

func TestRenderLiveFrame_ContainsJobName(t *testing.T) {
	stats := makeLiveStats(2)
	opts := DefaultLiveOptions()
	out := RenderLiveFrame(stats, opts)
	if !strings.Contains(out, "job-1") {
		t.Errorf("expected output to contain job name 'job-1', got:\n%s", out)
	}
}

func TestRenderLiveFrame_EmptyStats(t *testing.T) {
	opts := DefaultLiveOptions()
	out := RenderLiveFrame([]model.AggregateStats{}, opts)
	// Should not panic and should return something (e.g. empty state message or header)
	if out == "" {
		t.Error("expected at least an empty-state string, got empty output")
	}
}

func TestRenderLiveFrame_RespectsMaxRows(t *testing.T) {
	stats := makeLiveStats(20)
	opts := DefaultLiveOptions()
	opts.MaxRows = 5
	out := RenderLiveFrame(stats, opts)
	// job-6 through job-20 should not appear
	if strings.Contains(out, "job-6") {
		t.Errorf("expected output to respect MaxRows=5, but found 'job-6'")
	}
	if !strings.Contains(out, "job-1") {
		t.Errorf("expected output to contain 'job-1' within MaxRows=5")
	}
}

func TestRenderLiveFrame_ShowsTimestamp(t *testing.T) {
	stats := makeLiveStats(1)
	opts := DefaultLiveOptions()
	opts.ShowTimestamp = true
	out := RenderLiveFrame(stats, opts)
	// Timestamp format contains the current year at minimum
	year := fmt.Sprintf("%d", time.Now().Year())
	if !strings.Contains(out, year) {
		t.Errorf("expected output to contain current year %s when ShowTimestamp=true", year)
	}
}

func TestRenderLiveFrame_HidesTimestampWhenDisabled(t *testing.T) {
	stats := makeLiveStats(1)
	opts := DefaultLiveOptions()
	opts.ShowTimestamp = false
	out := RenderLiveFrame(stats, opts)
	year := fmt.Sprintf("%d", time.Now().Year())
	// Without timestamp, the year string should not appear in the frame
	if strings.Contains(out, year) {
		t.Errorf("expected output to hide timestamp when ShowTimestamp=false")
	}
}

func TestClamp_BelowMin(t *testing.T) {
	if got := clamp(3, 5, 10); got != 5 {
		t.Errorf("clamp(3,5,10) = %d, want 5", got)
	}
}

func TestClamp_AboveMax(t *testing.T) {
	if got := clamp(15, 5, 10); got != 10 {
		t.Errorf("clamp(15,5,10) = %d, want 10", got)
	}
}

func TestClamp_WithinRange(t *testing.T) {
	if got := clamp(7, 5, 10); got != 7 {
		t.Errorf("clamp(7,5,10) = %d, want 7", got)
	}
}
