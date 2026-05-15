package render

import (
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
)

func makeSparkbarStat(job string, failures, total int) *model.AggregateStats {
	rate := 0.0
	if total > 0 {
		rate = float64(failures) / float64(total)
	}
	return &model.AggregateStats{
		JobName:     job,
		TotalRuns:   total,
		Failures:    failures,
		FailureRate: rate,
	}
}

func TestRenderSparkbar_NilReturnsEmpty(t *testing.T) {
	opts := DefaultSparkbarOptions()
	result := RenderSparkbar(nil, opts)
	if result != "" {
		t.Errorf("expected empty string for nil stat, got %q", result)
	}
}

func TestRenderSparkbar_ContainsJobName(t *testing.T) {
	opts := DefaultSparkbarOptions()
	stat := makeSparkbarStat("backup-db", 1, 10)
	result := RenderSparkbar(stat, opts)
	if !strings.Contains(result, "backup-db") {
		t.Errorf("expected job name in output, got: %s", result)
	}
}

func TestRenderSparkbar_ContainsPercentage(t *testing.T) {
	opts := DefaultSparkbarOptions()
	stat := makeSparkbarStat("cleanup", 5, 10)
	result := RenderSparkbar(stat, opts)
	if !strings.Contains(result, "50.0%") {
		t.Errorf("expected 50.0%% in output, got: %s", result)
	}
}

func TestRenderSparkbar_ZeroFailures(t *testing.T) {
	opts := DefaultSparkbarOptions()
	stat := makeSparkbarStat("healthcheck", 0, 20)
	result := RenderSparkbar(stat, opts)
	if !strings.Contains(result, "0.0%") {
		t.Errorf("expected 0.0%% in output, got: %s", result)
	}
	if !strings.Contains(result, string(opts.EmptyChar)) {
		t.Errorf("expected empty chars in bar for zero failure rate")
	}
}

func TestRenderSparkbar_FullFailures(t *testing.T) {
	opts := DefaultSparkbarOptions()
	stat := makeSparkbarStat("broken-job", 10, 10)
	result := RenderSparkbar(stat, opts)
	if !strings.Contains(result, "100.0%") {
		t.Errorf("expected 100.0%% in output, got: %s", result)
	}
	if !strings.Contains(result, string(opts.FillChar)) {
		t.Errorf("expected fill chars in bar for full failure rate")
	}
}

func TestRenderSparkbar_HidesLabelWhenDisabled(t *testing.T) {
	opts := DefaultSparkbarOptions()
	opts.ShowLabel = false
	stat := makeSparkbarStat("myjob", 2, 10)
	result := RenderSparkbar(stat, opts)
	if strings.Contains(result, "myjob") {
		t.Errorf("expected job name to be hidden, got: %s", result)
	}
}

func TestRenderSparkbarList_EmptyReturnsPlaceholder(t *testing.T) {
	opts := DefaultSparkbarOptions()
	result := RenderSparkbarList(nil, opts)
	if result != "(no jobs)" {
		t.Errorf("expected placeholder for empty list, got: %q", result)
	}
}

func TestRenderSparkbarList_OneLinePerStat(t *testing.T) {
	opts := DefaultSparkbarOptions()
	stats := []*model.AggregateStats{
		makeSparkbarStat("job-a", 1, 5),
		makeSparkbarStat("job-b", 3, 5),
		makeSparkbarStat("job-c", 0, 5),
	}
	result := RenderSparkbarList(stats, opts)
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
