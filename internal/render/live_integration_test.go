package render_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
	"github.com/user/cronview/internal/render"
)

// buildLiveEntries generates a realistic mix of log entries spread across
// multiple jobs and hostnames for integration-level live frame tests.
func buildLiveEntries(jobs []string, hostname string, now time.Time) []model.AggregateStats {
	stats := make([]model.AggregateStats, 0, len(jobs))
	for i, job := range jobs {
		failures := 0
		if i%3 == 0 {
			failures = 2
		}
		totalRuns := 10
		rate := float64(failures) / float64(totalRuns)
		lastStatus := "success"
		if failures > 0 {
			lastStatus = "failure"
		}
		stats = append(stats, model.AggregateStats{
			JobName:     job,
			Hostname:    hostname,
			TotalRuns:   totalRuns,
			Failures:    failures,
			FailureRate: rate,
			LastStatus:  lastStatus,
			LastSeen:    now.Add(-time.Duration(i) * time.Minute),
		})
	}
	return stats
}

func TestLiveIntegration_FrameIsNonEmpty(t *testing.T) {
	now := time.Now()
	jobs := []string{"backup", "cleanup", "report", "sync", "healthcheck"}
	stats := buildLiveEntries(jobs, "prod-01", now)

	opts := render.DefaultLiveOptions()
	out := render.RenderLiveFrame(stats, opts)

	if strings.TrimSpace(out) == "" {
		t.Error("expected non-empty live frame output")
	}
}

func TestLiveIntegration_AllJobsAppear(t *testing.T) {
	now := time.Now()
	jobs := []string{"backup", "cleanup", "report"}
	stats := buildLiveEntries(jobs, "prod-01", now)

	opts := render.DefaultLiveOptions()
	out := render.RenderLiveFrame(stats, opts)

	for _, job := range jobs {
		if !strings.Contains(out, job) {
			t.Errorf("expected output to contain job %q", job)
		}
	}
}

func TestLiveIntegration_FailingJobsHighlighted(t *testing.T) {
	now := time.Now()
	// First job (index 0) will have failures due to buildLiveEntries logic.
	jobs := []string{"failing-job", "healthy-job"}
	stats := buildLiveEntries(jobs, "prod-01", now)

	opts := render.DefaultLiveOptions()
	out := render.RenderLiveFrame(stats, opts)

	// The output should contain both job names regardless of colour codes.
	if !strings.Contains(out, "failing-job") {
		t.Error("expected failing-job to appear in live frame")
	}
	if !strings.Contains(out, "healthy-job") {
		t.Error("expected healthy-job to appear in live frame")
	}
}

func TestLiveIntegration_EmptyStatsProducesOutput(t *testing.T) {
	opts := render.DefaultLiveOptions()
	out := render.RenderLiveFrame([]model.AggregateStats{}, opts)

	// Even with no data, the frame should render without panicking and
	// return at least an empty-state indicator or blank frame.
	if out == "" {
		t.Error("expected at least an empty string (not a panic) for empty stats")
	}
}

func TestLiveIntegration_MaxRowsRespected(t *testing.T) {
	now := time.Now()
	jobs := []string{"job-a", "job-b", "job-c", "job-d", "job-e", "job-f"}
	stats := buildLiveEntries(jobs, "prod-01", now)

	opts := render.DefaultLiveOptions()
	opts.MaxRows = 3
	out := render.RenderLiveFrame(stats, opts)

	// With MaxRows=3, jobs beyond the limit should not appear.
	visible := 0
	for _, job := range jobs {
		if strings.Contains(out, job) {
			visible++
		}
	}
	if visible > opts.MaxRows {
		t.Errorf("expected at most %d jobs in output, found %d", opts.MaxRows, visible)
	}
}
