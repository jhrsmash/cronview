package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cronview/internal/model"
	"github.com/user/cronview/internal/render"
)

var sampleStats = []model.JobStats{
	{JobName: "backup", Hostname: "srv1", TotalRuns: 10, Failures: 2, FailureRate: 0.2, LastStatus: "failure"},
	{JobName: "cleanup", Hostname: "srv2", TotalRuns: 5, Failures: 0, FailureRate: 0.0, LastStatus: "success"},
}

func TestRenderJobTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	render.RenderJobTable(&buf, sampleStats)
	out := buf.String()
	for _, hdr := range []string{"JOB", "RUNS", "FAILS", "FAIL%", "LAST", "HOSTNAME"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestRenderJobTable_ContainsJobNames(t *testing.T) {
	var buf bytes.Buffer
	render.RenderJobTable(&buf, sampleStats)
	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Error("expected job name 'backup' in output")
	}
	if !strings.Contains(out, "cleanup") {
		t.Error("expected job name 'cleanup' in output")
	}
}

func TestRenderJobTable_FailureRate(t *testing.T) {
	var buf bytes.Buffer
	render.RenderJobTable(&buf, sampleStats)
	out := buf.String()
	if !strings.Contains(out, "20.0%") {
		t.Errorf("expected failure rate '20.0%%' in output, got:\n%s", out)
	}
}

func TestRenderJobTable_EmptyList(t *testing.T) {
	var buf bytes.Buffer
	render.RenderJobTable(&buf, []model.JobStats{})
	out := buf.String()
	if !strings.Contains(out, "0 job(s) listed") {
		t.Errorf("expected '0 job(s) listed' in output, got:\n%s", out)
	}
}

func TestRenderSummary_OK(t *testing.T) {
	var buf bytes.Buffer
	sum := model.Summary{UniqueJobs: 2, TotalRuns: 15, TotalFailures: 0, OverallFailureRate: 0.0}
	render.RenderSummary(&buf, sum)
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] status, got: %s", out)
	}
}

func TestRenderSummary_Degraded(t *testing.T) {
	var buf bytes.Buffer
	sum := model.Summary{UniqueJobs: 2, TotalRuns: 15, TotalFailures: 3, OverallFailureRate: 0.2}
	render.RenderSummary(&buf, sum)
	out := buf.String()
	if !strings.Contains(out, "[DEGRADED]") {
		t.Errorf("expected [DEGRADED] status, got: %s", out)
	}
}
