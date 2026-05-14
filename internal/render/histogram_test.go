package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeHistogramEntry(status model.Status, offset time.Duration) model.LogEntry {
	return model.LogEntry{
		Time:     time.Now().Add(-offset),
		Hostname: "host1",
		JobName:  "backup",
		Status:   status,
	}
}

func TestRenderHistogram_NoFailures(t *testing.T) {
	entries := []model.LogEntry{
		makeHistogramEntry(model.StatusSuccess, 1*time.Hour),
		makeHistogramEntry(model.StatusSuccess, 2*time.Hour),
	}
	var buf bytes.Buffer
	RenderHistogram(&buf, entries, DefaultHistogramOptions())
	if !strings.Contains(buf.String(), "no failures") {
		t.Errorf("expected 'no failures' message, got: %s", buf.String())
	}
}

func TestRenderHistogram_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	RenderHistogram(&buf, nil, DefaultHistogramOptions())
	if !strings.Contains(buf.String(), "no failures") {
		t.Errorf("expected 'no failures' message for empty input, got: %s", buf.String())
	}
}

func TestRenderHistogram_ContainsAxisLine(t *testing.T) {
	entries := []model.LogEntry{
		makeHistogramEntry(model.StatusFailure, 30*time.Minute),
		makeHistogramEntry(model.StatusFailure, 1*time.Hour),
	}
	var buf bytes.Buffer
	RenderHistogram(&buf, entries, DefaultHistogramOptions())
	output := buf.String()
	if !strings.Contains(output, "---") {
		t.Errorf("expected axis dashes in output, got: %s", output)
	}
}

func TestRenderHistogram_ContainsLabels(t *testing.T) {
	entries := []model.LogEntry{
		makeHistogramEntry(model.StatusFailure, 1*time.Hour),
	}
	var buf bytes.Buffer
	RenderHistogram(&buf, entries, DefaultHistogramOptions())
	output := buf.String()
	if !strings.Contains(output, "older") || !strings.Contains(output, "newer") {
		t.Errorf("expected 'older'/'newer' labels in output, got: %s", output)
	}
}

func TestRenderHistogram_RecentFailureInLastBucket(t *testing.T) {
	opts := DefaultHistogramOptions()
	opts.Buckets = 12
	opts.MaxHeight = 4

	entries := []model.LogEntry{
		makeHistogramEntry(model.StatusFailure, 5*time.Minute),
	}
	var buf bytes.Buffer
	RenderHistogram(&buf, entries, opts)
	output := buf.String()
	if !strings.Contains(output, "█") {
		t.Errorf("expected bar character in output for recent failure, got: %s", output)
	}
}

func TestRenderHistogram_InvalidOptions(t *testing.T) {
	opts := DefaultHistogramOptions()
	opts.Buckets = 0

	var buf bytes.Buffer
	RenderHistogram(&buf, []model.LogEntry{makeHistogramEntry(model.StatusFailure, time.Minute)}, opts)
	if buf.Len() != 0 {
		t.Errorf("expected no output for zero buckets, got: %s", buf.String())
	}
}

func TestBucketFailures_IgnoresOldEntries(t *testing.T) {
	entries := []model.LogEntry{
		makeHistogramEntry(model.StatusFailure, 48*time.Hour), // outside 24h window
		makeHistogramEntry(model.StatusFailure, 1*time.Hour),
	}
	counts := bucketFailures(entries, 24, 24*time.Hour)
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 1 {
		t.Errorf("expected 1 failure in window, got %d", total)
	}
}
