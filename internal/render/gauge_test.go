package render

import (
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
)

func makeGaugeStats(job string, failures, total int) model.AggregateStats {
	rate := 0.0
	if total > 0 {
		rate = float64(failures) / float64(total)
	}
	return model.AggregateStats{
		JobName:     job,
		TotalRuns:   total,
		Failures:    failures,
		FailureRate: rate,
	}
}

func TestRenderGauge_NilReturnsEmpty(t *testing.T) {
	out := RenderGauge(nil, DefaultGaugeOptions())
	if out != "" {
		t.Errorf("expected empty string for nil stats, got %q", out)
	}
}

func TestRenderGauge_ContainsJobName(t *testing.T) {
	s := makeGaugeStats("backup.sh", 2, 10)
	out := RenderGauge(&s, DefaultGaugeOptions())
	if !strings.Contains(out, "backup.sh") {
		t.Errorf("expected job name in output, got: %s", out)
	}
}

func TestRenderGauge_ContainsPercentage(t *testing.T) {
	s := makeGaugeStats("sync", 1, 4) // 25%
	out := RenderGauge(&s, DefaultGaugeOptions())
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected '25.0%%' in output, got: %s", out)
	}
}

func TestRenderGauge_CriticalUsesSolidFill(t *testing.T) {
	s := makeGaugeStats("job", 6, 10) // 60% >= crit threshold 50%
	opts := DefaultGaugeOptions()
	out := RenderGauge(&s, opts)
	if !strings.Contains(out, "█") {
		t.Errorf("expected solid fill char for critical rate, got: %s", out)
	}
}

func TestRenderGauge_LowRateUsesLightFill(t *testing.T) {
	s := makeGaugeStats("job", 1, 20) // 5% < warn threshold
	out := RenderGauge(&s, DefaultGaugeOptions())
	if !strings.Contains(out, "░") {
		t.Errorf("expected light fill char for low rate, got: %s", out)
	}
}

func TestRenderGauge_HideLabelOption(t *testing.T) {
	s := makeGaugeStats("nightly", 0, 5)
	opts := DefaultGaugeOptions()
	opts.ShowLabel = false
	out := RenderGauge(&s, opts)
	if strings.Contains(out, "nightly") {
		t.Errorf("expected job name hidden when ShowLabel=false, got: %s", out)
	}
}

func TestRenderGaugeList_Empty(t *testing.T) {
	out := RenderGaugeList(nil, DefaultGaugeOptions(), 10)
	if !strings.Contains(out, "no jobs") {
		t.Errorf("expected 'no jobs' for empty list, got: %s", out)
	}
}

func TestRenderGaugeList_RespectsMaxRows(t *testing.T) {
	stats := []model.AggregateStats{
		makeGaugeStats("a", 1, 10),
		makeGaugeStats("b", 2, 10),
		makeGaugeStats("c", 3, 10),
	}
	out := RenderGaugeList(stats, DefaultGaugeOptions(), 2)
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines with maxRows=2, got %d: %s", len(lines), out)
	}
}

func TestRenderGaugeList_AllJobsPresent(t *testing.T) {
	stats := []model.AggregateStats{
		makeGaugeStats("alpha", 0, 5),
		makeGaugeStats("beta", 1, 5),
	}
	out := RenderGaugeList(stats, DefaultGaugeOptions(), 0)
	for _, name := range []string{"alpha", "beta"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected %q in gauge list output", name)
		}
	}
}
