package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeHeatmapStats(jobName string, entries []model.LogEntry) model.AggregateStats {
	return model.AggregateStats{
		JobName: jobName,
		Entries: entries,
	}
}

func TestRenderHeatmap_ContainsJobName(t *testing.T) {
	stats := makeHeatmapStats("backup-daily", nil)
	var buf bytes.Buffer
	RenderHeatmap(&buf, stats, DefaultHeatmapOptions())
	if !strings.Contains(buf.String(), "backup-daily") {
		t.Errorf("expected job name in output, got: %q", buf.String())
	}
}

func TestRenderHeatmap_EmptyEntries_ShowsDots(t *testing.T) {
	stats := makeHeatmapStats("noop", nil)
	var buf bytes.Buffer
	opts := DefaultHeatmapOptions()
	opts.Days = 7
	RenderHeatmap(&buf, stats, opts)
	// strip ANSI: count raw '·' characters
	stripped := stripANSI(buf.String())
	count := strings.Count(stripped, "·")
	if count != 7 {
		t.Errorf("expected 7 dots for empty buckets, got %d in %q", count, stripped)
	}
}

func TestRenderHeatmap_AllSuccessShowsLowShade(t *testing.T) {
	now := time.Now().UTC()
	entries := []model.LogEntry{
		{Time: now, Status: "success"},
		{Time: now, Status: "success"},
	}
	stats := makeHeatmapStats("myjob", entries)
	var buf bytes.Buffer
	RenderHeatmap(&buf, stats, DefaultHeatmapOptions())
	stripped := stripANSI(buf.String())
	// today's bucket should render as '░' (no failures)
	if !strings.Contains(stripped, "░") {
		t.Errorf("expected '░' cell for zero failures, got: %q", stripped)
	}
}

func TestRenderHeatmap_AllFailuresShowsHighShade(t *testing.T) {
	now := time.Now().UTC()
	entries := []model.LogEntry{
		{Time: now, Status: "failure"},
		{Time: now, Status: "failure"},
	}
	stats := makeHeatmapStats("broken", entries)
	var buf bytes.Buffer
	RenderHeatmap(&buf, stats, DefaultHeatmapOptions())
	stripped := stripANSI(buf.String())
	if !strings.Contains(stripped, "█") {
		t.Errorf("expected '█' cell for 100%% failures, got: %q", stripped)
	}
}

func TestRenderHeatmap_OldEntriesIgnored(t *testing.T) {
	old := time.Now().UTC().AddDate(0, 0, -30)
	entries := []model.LogEntry{
		{Time: old, Status: "failure"},
	}
	stats := makeHeatmapStats("stale", entries)
	var buf bytes.Buffer
	opts := DefaultHeatmapOptions()
	opts.Days = 7
	RenderHeatmap(&buf, stats, opts)
	stripped := stripANSI(buf.String())
	// all buckets should be dots since the entry is outside the window
	if strings.Contains(stripped, "█") || strings.Contains(stripped, "░") {
		t.Errorf("old entry should not appear in heatmap, got: %q", stripped)
	}
}

func TestRenderHeatmap_CustomLabel(t *testing.T) {
	stats := makeHeatmapStats("original", nil)
	opts := DefaultHeatmapOptions()
	opts.JobName = "override-label"
	var buf bytes.Buffer
	RenderHeatmap(&buf, stats, opts)
	if !strings.Contains(buf.String(), "override-label") {
		t.Errorf("expected custom label in output")
	}
}

// stripANSI removes ANSI escape sequences for plain-text assertions.
func stripANSI(s string) string {
	var out strings.Builder
	inEsc := false
	for _, r := range s {
		switch {
		case r == '\033':
			inEsc = true
		case inEsc && r == 'm':
			inEsc = false
		case !inEsc:
			out.WriteRune(r)
		}
	}
	return out.String()
}
