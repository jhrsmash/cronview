package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeDetailStats(jobName, host string, runs, failures int) model.JobStats {
	now := time.Now()
	entries := make([]model.LogEntry, runs)
	for i := 0; i < runs; i++ {
		status := "success"
		if i < failures {
			status = "failure"
		}
		entries[i] = model.LogEntry{
			Timestamp: now.Add(time.Duration(-runs+i) * time.Minute),
			Status:    status,
			Message:   "cron job ran",
			JobName:   jobName,
			Hostname:  host,
		}
	}
	var lastStatus string
	if runs > 0 {
		lastStatus = entries[runs-1].Status
	}
	return model.JobStats{
		JobName:     jobName,
		Hostname:    host,
		TotalRuns:   runs,
		Failures:    failures,
		FailureRate: float64(failures) / float64(max(runs, 1)),
		LastStatus:  lastStatus,
		LastRun:     now,
		Entries:     entries,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestRenderDetail_ContainsJobName(t *testing.T) {
	var sb strings.Builder
	stats := makeDetailStats("backup-db", "host1", 5, 1)
	RenderDetail(&sb, stats, DefaultDetailOptions())
	if !strings.Contains(sb.String(), "backup-db") {
		t.Error("expected job name in detail output")
	}
}

func TestRenderDetail_ContainsHostname(t *testing.T) {
	var sb strings.Builder
	stats := makeDetailStats("cleanup", "web-01", 3, 0)
	RenderDetail(&sb, stats, DefaultDetailOptions())
	if !strings.Contains(sb.String(), "web-01") {
		t.Error("expected hostname in detail output")
	}
}

func TestRenderDetail_FailureRate(t *testing.T) {
	var sb strings.Builder
	stats := makeDetailStats("sync", "host2", 10, 3)
	RenderDetail(&sb, stats, DefaultDetailOptions())
	out := sb.String()
	if !strings.Contains(out, "30.0%") {
		t.Errorf("expected 30.0%% failure rate, got:\n%s", out)
	}
}

func TestRenderDetail_LimitsHistory(t *testing.T) {
	var sb strings.Builder
	stats := makeDetailStats("prune", "host3", 20, 2)
	opts := DefaultDetailOptions()
	opts.ShowHistory = 5
	RenderDetail(&sb, stats, opts)
	out := sb.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	historyLines := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "[") {
			historyLines++
		}
	}
	if historyLines != 5 {
		t.Errorf("expected 5 history lines, got %d", historyLines)
	}
}

func TestRenderDetail_ZeroRuns(t *testing.T) {
	var sb strings.Builder
	stats := makeDetailStats("empty-job", "host4", 0, 0)
	RenderDetail(&sb, stats, DefaultDetailOptions())
	out := sb.String()
	if !strings.Contains(out, "empty-job") {
		t.Error("expected job name even with zero runs")
	}
	if !strings.Contains(out, "Total Runs   : 0") {
		t.Error("expected zero total runs")
	}
}
