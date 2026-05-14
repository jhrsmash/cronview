package render

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func makeEntry(status model.Status, age time.Duration) model.LogEntry {
	return model.LogEntry{
		Timestamp: time.Now().Add(-age),
		Status:    status,
		JobName:   "test-job",
		Hostname:  "host1",
	}
}

func TestRenderSparkline_LengthMatchesBuckets(t *testing.T) {
	opts := SparklineOptions{Buckets: 10, Window: 10 * 24 * time.Hour}
	result := RenderSparkline(nil, opts)
	// result is "[" + 10 runes + "]"
	runes := []rune(result)
	if len(runes) != 12 {
		t.Errorf("expected length 12 (brackets + 10 blocks), got %d", len(runes))
	}
}

func TestRenderSparkline_EmptyEntries(t *testing.T) {
	opts := DefaultSparklineOptions()
	result := RenderSparkline([]model.LogEntry{}, opts)
	if !strings.HasPrefix(result, "[") || !strings.HasSuffix(result, "]") {
		t.Errorf("expected bracketed sparkline, got %q", result)
	}
	// All buckets should be lowest block character
	inner := []rune(result)[1 : len([]rune(result))-1]
	for _, r := range inner {
		if r != sparklineBlocks[0] {
			t.Errorf("expected all low blocks for empty input, got rune %q", r)
		}
	}
}

func TestRenderSparkline_RecentFailureInLastBucket(t *testing.T) {
	opts := SparklineOptions{Buckets: 7, Window: 7 * 24 * time.Hour}
	entries := []model.LogEntry{
		makeEntry(model.StatusFailure, 1*time.Hour),
	}
	result := RenderSparkline(entries, opts)
	runes := []rune(result)
	// Last bucket (index Buckets) should be highest block
	lastBlock := runes[len(runes)-2]
	if lastBlock != sparklineBlocks[len(sparklineBlocks)-1] {
		t.Errorf("expected full block in last bucket, got %q", lastBlock)
	}
}

func TestRenderSparkline_SuccessEntriesIgnored(t *testing.T) {
	opts := SparklineOptions{Buckets: 5, Window: 5 * 24 * time.Hour}
	entries := []model.LogEntry{
		makeEntry(model.StatusSuccess, 1*time.Hour),
		makeEntry(model.StatusSuccess, 2*time.Hour),
	}
	result := RenderSparkline(entries, opts)
	inner := []rune(result)[1 : len([]rune(result))-1]
	for _, r := range inner {
		if r != sparklineBlocks[0] {
			t.Errorf("success entries should not affect sparkline, got rune %q", r)
		}
	}
}

func TestRenderSparkline_ZeroBucketsReturnsEmpty(t *testing.T) {
	opts := SparklineOptions{Buckets: 0, Window: 24 * time.Hour}
	result := RenderSparkline(nil, opts)
	if result != "" {
		t.Errorf("expected empty string for zero buckets, got %q", result)
	}
}

func TestRenderSparkline_OutOfWindowEntriesIgnored(t *testing.T) {
	opts := SparklineOptions{Buckets: 7, Window: 7 * 24 * time.Hour}
	entries := []model.LogEntry{
		makeEntry(model.StatusFailure, 30*24*time.Hour), // 30 days ago, outside window
	}
	result := RenderSparkline(entries, opts)
	inner := []rune(result)[1 : len([]rune(result))-1]
	for _, r := range inner {
		if r != sparklineBlocks[0] {
			t.Errorf("out-of-window entry should be ignored, got rune %q", r)
		}
	}
}
