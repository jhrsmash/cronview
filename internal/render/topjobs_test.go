package render

import (
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
)

func makeTopStats() []model.AggregateStats {
	return []model.AggregateStats{
		{JobName: "backup", TotalRuns: 100, FailureCount: 40, FailureRate: 0.40},
		{JobName: "cleanup", TotalRuns: 50, FailureCount: 5, FailureRate: 0.10},
		{JobName: "sync", TotalRuns: 200, FailureCount: 80, FailureRate: 0.40},
		{JobName: "report", TotalRuns: 20, FailureCount: 0, FailureRate: 0.00},
		{JobName: "deploy", TotalRuns: 10, FailureCount: 9, FailureRate: 0.90},
	}
}

func TestRenderTopJobs_ExcludesZeroFailures(t *testing.T) {
	var sb strings.Builder
	RenderTopJobs(&sb, makeTopStats(), DefaultTopJobsOptions())
	out := sb.String()
	if strings.Contains(out, "report") {
		t.Error("expected job with zero failures to be excluded")
	}
}

func TestRenderTopJobs_SortByFailures(t *testing.T) {
	var sb strings.Builder
	opts := DefaultTopJobsOptions()
	opts.SortBy = "failures"
	RenderTopJobs(&sb, makeTopStats(), opts)
	out := sb.String()
	syncIdx := strings.Index(out, "sync")
	backupIdx := strings.Index(out, "backup")
	if syncIdx == -1 || backupIdx == -1 {
		t.Fatal("expected sync and backup in output")
	}
	if syncIdx > backupIdx {
		t.Error("expected sync (80 failures) to appear before backup (40 failures)")
	}
}

func TestRenderTopJobs_SortByRate(t *testing.T) {
	var sb strings.Builder
	opts := DefaultTopJobsOptions()
	opts.SortBy = "rate"
	RenderTopJobs(&sb, makeTopStats(), opts)
	out := sb.String()
	deployIdx := strings.Index(out, "deploy")
	syncIdx := strings.Index(out, "sync")
	if deployIdx == -1 || syncIdx == -1 {
		t.Fatal("expected deploy and sync in output")
	}
	if deployIdx > syncIdx {
		t.Error("expected deploy (90%% rate) to appear before sync (40%% rate)")
	}
}

func TestRenderTopJobs_RespectsMaxRows(t *testing.T) {
	var sb strings.Builder
	opts := DefaultTopJobsOptions()
	opts.MaxRows = 2
	RenderTopJobs(&sb, makeTopStats(), opts)
	out := sb.String()
	// Only 2 data rows should appear (header + separator + 2 jobs)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// header + separator + 2 data rows = 4
	if len(lines) != 4 {
		t.Errorf("expected 4 lines for MaxRows=2, got %d: %q", len(lines), out)
	}
}

func TestRenderTopJobs_EmptyStats(t *testing.T) {
	var sb strings.Builder
	RenderTopJobs(&sb, []model.AggregateStats{}, DefaultTopJobsOptions())
	out := sb.String()
	if !strings.Contains(out, "No failures") {
		t.Errorf("expected 'No failures' message for empty input, got: %q", out)
	}
}

func TestRenderTopJobs_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	RenderTopJobs(&sb, makeTopStats(), DefaultTopJobsOptions())
	out := sb.String()
	for _, hdr := range []string{"JOB", "RUNS", "FAIL", "RATE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}
