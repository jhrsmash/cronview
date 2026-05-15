package render

import (
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
)

func makeDiffStats(name string, rate float64) model.AggregateStats {
	return model.AggregateStats{
		JobName:     name,
		FailureRate: rate,
	}
}

func TestRenderDiff_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	RenderDiff(&sb, nil, nil, DefaultDiffViewOptions())
	out := sb.String()
	if !strings.Contains(out, "JOB") || !strings.Contains(out, "BEFORE") || !strings.Contains(out, "AFTER") {
		t.Errorf("expected headers in output, got:\n%s", out)
	}
}

func TestRenderDiff_ShowsJobName(t *testing.T) {
	before := []model.AggregateStats{makeDiffStats("backup.sh", 0.1)}
	after := []model.AggregateStats{makeDiffStats("backup.sh", 0.2)}

	var sb strings.Builder
	RenderDiff(&sb, before, after, DefaultDiffViewOptions())
	out := sb.String()

	if !strings.Contains(out, "backup.sh") {
		t.Errorf("expected job name in output, got:\n%s", out)
	}
}

func TestRenderDiff_WorseningMarker(t *testing.T) {
	before := []model.AggregateStats{makeDiffStats("sync.sh", 0.05)}
	after := []model.AggregateStats{makeDiffStats("sync.sh", 0.50)}

	var sb strings.Builder
	RenderDiff(&sb, before, after, DefaultDiffViewOptions())
	out := sb.String()

	if !strings.Contains(out, "↑") {
		t.Errorf("expected worsening marker ↑ in output, got:\n%s", out)
	}
}

func TestRenderDiff_ImprovingMarker(t *testing.T) {
	before := []model.AggregateStats{makeDiffStats("cleanup.sh", 0.80)}
	after := []model.AggregateStats{makeDiffStats("cleanup.sh", 0.10)}

	var sb strings.Builder
	RenderDiff(&sb, before, after, DefaultDiffViewOptions())
	out := sb.String()

	if !strings.Contains(out, "↓") {
		t.Errorf("expected improving marker ↓ in output, got:\n%s", out)
	}
}

func TestRenderDiff_NewJobMarker(t *testing.T) {
	after := []model.AggregateStats{makeDiffStats("newjob.sh", 0.3)}

	var sb strings.Builder
	RenderDiff(&sb, nil, after, DefaultDiffViewOptions())
	out := sb.String()

	if !strings.Contains(out, "N") {
		t.Errorf("expected new job marker N in output, got:\n%s", out)
	}
}

func TestRenderDiff_DeltaHidden(t *testing.T) {
	before := []model.AggregateStats{makeDiffStats("job.sh", 0.1)}
	after := []model.AggregateStats{makeDiffStats("job.sh", 0.4)}

	opts := DefaultDiffViewOptions()
	opts.ShowDelta = false

	var sb strings.Builder
	RenderDiff(&sb, before, after, opts)
	out := sb.String()

	if strings.Contains(out, "+30.0%") {
		t.Errorf("expected delta to be hidden, but found it in output:\n%s", out)
	}
}
