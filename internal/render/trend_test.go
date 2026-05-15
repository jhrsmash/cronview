package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeTrendStats(jobName string, baselineFailures, recentFailures, windowSize int) model.AggregateStats {
	now := time.Now()
	entries := make([]model.LogEntry, windowSize*2)
	for i := 0; i < windowSize; i++ {
		status := model.StatusSuccess
		if i < baselineFailures {
			status = model.StatusFailure
		}
		entries[i] = model.LogEntry{Timestamp: now.Add(-time.Duration(windowSize*2-i) * time.Minute), Status: status, JobName: jobName}
	}
	for i := 0; i < windowSize; i++ {
		status := model.StatusSuccess
		if i < recentFailures {
			status = model.StatusFailure
		}
		entries[windowSize+i] = model.LogEntry{Timestamp: now.Add(-time.Duration(windowSize-i) * time.Minute), Status: status, JobName: jobName}
	}
	return model.AggregateStats{
		JobName:     jobName,
		FailureRate: float64(baselineFailures+recentFailures) / float64(windowSize*2),
		Entries:     entries,
	}
}

func TestComputeTrend_Stable(t *testing.T) {
	stats := makeTrendStats("backup", 1, 1, 10)
	opts := DefaultTrendOptions()
	dir := ComputeTrend(stats, opts)
	if dir != TrendStable {
		t.Errorf("expected stable, got %d", dir)
	}
}

func TestComputeTrend_Worsening(t *testing.T) {
	stats := makeTrendStats("backup", 0, 8, 10)
	opts := DefaultTrendOptions()
	dir := ComputeTrend(stats, opts)
	if dir != TrendWorsening {
		t.Errorf("expected worsening, got %d", dir)
	}
}

func TestComputeTrend_Improving(t *testing.T) {
	stats := makeTrendStats("backup", 8, 0, 10)
	opts := DefaultTrendOptions()
	dir := ComputeTrend(stats, opts)
	if dir != TrendImproving {
		t.Errorf("expected improving, got %d", dir)
	}
}

func TestComputeTrend_TooFewEntries(t *testing.T) {
	stats := model.AggregateStats{
		JobName: "tiny",
		Entries: []model.LogEntry{{Status: model.StatusFailure}},
	}
	opts := DefaultTrendOptions()
	dir := ComputeTrend(stats, opts)
	if dir != TrendStable {
		t.Errorf("expected stable for insufficient data, got %d", dir)
	}
}

func TestRenderTrendBadge_Improving(t *testing.T) {
	badge := RenderTrendBadge(TrendImproving)
	if !strings.Contains(badge, "improving") {
		t.Errorf("expected 'improving' in badge, got %q", badge)
	}
}

func TestRenderTrendBadge_Worsening(t *testing.T) {
	badge := RenderTrendBadge(TrendWorsening)
	if !strings.Contains(badge, "worsening") {
		t.Errorf("expected 'worsening' in badge, got %q", badge)
	}
}

func TestRenderTrendTable_ContainsJobName(t *testing.T) {
	stats := makeTrendStats("nightly-sync", 1, 1, 10)
	out := RenderTrendTable([]model.AggregateStats{stats}, DefaultTrendOptions())
	if !strings.Contains(out, "nightly-sync") {
		t.Errorf("expected job name in trend table output")
	}
}

func TestRenderTrendTable_ContainsHeaders(t *testing.T) {
	out := RenderTrendTable(nil, DefaultTrendOptions())
	if !strings.Contains(out, "Job") || !strings.Contains(out, "Trend") {
		t.Errorf("expected headers in trend table output")
	}
}
