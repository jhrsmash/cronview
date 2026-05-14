package render

import (
	"strings"
	"testing"
	"time"
)

func makeStatusBarData() StatusBarData {
	return StatusBarData{
		TotalJobs:    20,
		FilteredJobs: 8,
		ActiveFilter: "backup",
		LastRefresh:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
		Page:         2,
		TotalPages:   4,
	}
}

func TestRenderStatusBar_ContainsJobCounts(t *testing.T) {
	var buf strings.Builder
	RenderStatusBar(&buf, makeStatusBarData(), DefaultStatusBarOptions())
	out := buf.String()
	if !strings.Contains(out, "8/20") {
		t.Errorf("expected job count 8/20 in output, got: %s", out)
	}
}

func TestRenderStatusBar_ContainsFilter(t *testing.T) {
	var buf strings.Builder
	opts := DefaultStatusBarOptions()
	opts.ShowFilter = true
	RenderStatusBar(&buf, makeStatusBarData(), opts)
	out := buf.String()
	if !strings.Contains(out, "filter:backup") {
		t.Errorf("expected filter label in output, got: %s", out)
	}
}

func TestRenderStatusBar_HidesFilterWhenDisabled(t *testing.T) {
	var buf strings.Builder
	opts := DefaultStatusBarOptions()
	opts.ShowFilter = false
	RenderStatusBar(&buf, makeStatusBarData(), opts)
	out := buf.String()
	if strings.Contains(out, "filter:") {
		t.Errorf("expected no filter label when ShowFilter=false, got: %s", out)
	}
}

func TestRenderStatusBar_ContainsPagination(t *testing.T) {
	var buf strings.Builder
	RenderStatusBar(&buf, makeStatusBarData(), DefaultStatusBarOptions())
	out := buf.String()
	if !strings.Contains(out, "page:2/4") {
		t.Errorf("expected page info in output, got: %s", out)
	}
}

func TestRenderStatusBar_HidesPaginationOnSinglePage(t *testing.T) {
	var buf strings.Builder
	data := makeStatusBarData()
	data.Page = 1
	data.TotalPages = 1
	RenderStatusBar(&buf, data, DefaultStatusBarOptions())
	out := buf.String()
	if strings.Contains(out, "page:") {
		t.Errorf("expected no page info for single page, got: %s", out)
	}
}

func TestRenderStatusBar_ContainsTimestamp(t *testing.T) {
	var buf strings.Builder
	opts := DefaultStatusBarOptions()
	opts.ShowTime = true
	RenderStatusBar(&buf, makeStatusBarData(), opts)
	out := buf.String()
	if !strings.Contains(out, "14:30:00") {
		t.Errorf("expected timestamp in output, got: %s", out)
	}
}

func TestRenderStatusBar_EmptyFilter(t *testing.T) {
	var buf strings.Builder
	data := makeStatusBarData()
	data.ActiveFilter = ""
	RenderStatusBar(&buf, data, DefaultStatusBarOptions())
	out := buf.String()
	if strings.Contains(out, "filter:") {
		t.Errorf("expected no filter label when ActiveFilter is empty, got: %s", out)
	}
}

func TestRenderStatusBar_OutputEndsWithNewline(t *testing.T) {
	var buf strings.Builder
	RenderStatusBar(&buf, makeStatusBarData(), DefaultStatusBarOptions())
	out := buf.String()
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("expected output to end with newline, got: %q", out)
	}
}
