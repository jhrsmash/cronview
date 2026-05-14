package render

import (
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
)

func makeMiniStats(rate float64) model.JobStats {
	return model.JobStats{
		JobName:     "job",
		TotalRuns:   10,
		Failures:    int(rate * 10),
		FailureRate: rate,
		LastStatus:  "success",
	}
}

func TestRenderMinimap_EmptyStats(t *testing.T) {
	var sb strings.Builder
	RenderMinimap(&sb, nil, DefaultMinimapOptions())
	if !strings.Contains(sb.String(), "no jobs") {
		t.Errorf("expected '(no jobs)' for empty input, got: %q", sb.String())
	}
}

func TestRenderMinimap_HealthyCell(t *testing.T) {
	var sb strings.Builder
	RenderMinimap(&sb, []model.JobStats{makeMiniStats(0.0)}, DefaultMinimapOptions())
	if !strings.Contains(sb.String(), "█") {
		t.Errorf("expected healthy cell '█', got: %q", sb.String())
	}
}

func TestRenderMinimap_DegradedCell(t *testing.T) {
	var sb strings.Builder
	RenderMinimap(&sb, []model.JobStats{makeMiniStats(0.15)}, DefaultMinimapOptions())
	if !strings.Contains(sb.String(), "▒") {
		t.Errorf("expected degraded cell '▒', got: %q", sb.String())
	}
}

func TestRenderMinimap_CriticalCell(t *testing.T) {
	var sb strings.Builder
	RenderMinimap(&sb, []model.JobStats{makeMiniStats(0.50)}, DefaultMinimapOptions())
	if !strings.Contains(sb.String(), "░") {
		t.Errorf("expected critical cell '░', got: %q", sb.String())
	}
}

func TestRenderMinimap_WrapsAtWidth(t *testing.T) {
	stats := make([]model.JobStats, 15)
	for i := range stats {
		stats[i] = makeMiniStats(0.0)
	}
	opts := DefaultMinimapOptions()
	opts.Width = 5
	var sb strings.Builder
	RenderMinimap(&sb, stats, opts)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	// first line should be the grid; with 15 items and width 5 we expect 3 grid rows
	gridLines := 0
	for _, l := range lines {
		if strings.Contains(l, "█") || strings.Contains(l, "▒") || strings.Contains(l, "░") {
			gridLines++
		}
	}
	if gridLines != 3 {
		t.Errorf("expected 3 grid rows, got %d; output:\n%s", gridLines, sb.String())
	}
}

func TestRenderMinimap_LegendShown(t *testing.T) {
	var sb strings.Builder
	RenderMinimap(&sb, []model.JobStats{makeMiniStats(0.0)}, DefaultMinimapOptions())
	if !strings.Contains(sb.String(), "healthy") {
		t.Errorf("expected legend in output, got: %q", sb.String())
	}
}

func TestRenderMinimap_LegendHidden(t *testing.T) {
	opts := DefaultMinimapOptions()
	opts.ShowLegend = false
	var sb strings.Builder
	RenderMinimap(&sb, []model.JobStats{makeMiniStats(0.0)}, opts)
	if strings.Contains(sb.String(), "healthy") {
		t.Errorf("expected no legend in output, got: %q", sb.String())
	}
}
