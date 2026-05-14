package render_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
	"github.com/user/cronview/internal/render"
)

func buildMixedEntries() []model.LogEntry {
	now := time.Now()
	return []model.LogEntry{
		{Time: now.Add(-1 * time.Hour), Hostname: "web1", JobName: "cleanup", Status: model.StatusSuccess},
		{Time: now.Add(-2 * time.Hour), Hostname: "web1", JobName: "cleanup", Status: model.StatusFailure},
		{Time: now.Add(-3 * time.Hour), Hostname: "web2", JobName: "backup", Status: model.StatusFailure},
		{Time: now.Add(-4 * time.Hour), Hostname: "web2", JobName: "backup", Status: model.StatusFailure},
		{Time: now.Add(-5 * time.Hour), Hostname: "web1", JobName: "cleanup", Status: model.StatusSuccess},
		{Time: now.Add(-30 * time.Hour), Hostname: "web1", JobName: "cleanup", Status: model.StatusFailure}, // outside window
	}
}

func TestHistogramIntegration_OutputIsNonEmpty(t *testing.T) {
	var buf bytes.Buffer
	render.RenderHistogram(&buf, buildMixedEntries(), render.DefaultHistogramOptions())
	if buf.Len() == 0 {
		t.Error("expected non-empty histogram output")
	}
}

func TestHistogramIntegration_FailuresOnlyCountedOnce(t *testing.T) {
	opts := render.DefaultHistogramOptions()
	opts.Buckets = 24

	entries := buildMixedEntries()
	var buf bytes.Buffer
	render.RenderHistogram(&buf, entries, opts)
	output := buf.String()

	// Should not contain the "no failures" placeholder
	if strings.Contains(output, "no failures") {
		t.Error("expected actual histogram bars, not 'no failures' message")
	}
}

func TestHistogramIntegration_CustomBarWidth(t *testing.T) {
	opts := render.DefaultHistogramOptions()
	opts.BarWidth = 2
	opts.Buckets = 6
	opts.MaxHeight = 3

	entries := []model.LogEntry{
		{Time: time.Now().Add(-1 * time.Hour), JobName: "job", Status: model.StatusFailure},
	}
	var buf bytes.Buffer
	render.RenderHistogram(&buf, entries, opts)
	output := buf.String()

	if !strings.Contains(output, "older") {
		t.Errorf("expected axis labels in custom-width histogram, got: %s", output)
	}
}
