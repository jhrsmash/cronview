package export_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cronview/internal/export"
	"github.com/user/cronview/internal/model"
)

var testTime = time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)

func sampleStats() []model.JobStats {
	return []model.JobStats{
		{
			JobName:     "backup",
			Hostname:    "host1",
			TotalRuns:   10,
			Failures:    2,
			FailureRate: 0.2,
			LastStatus:  "failure",
			LastRun:     testTime,
		},
		{
			JobName:     "cleanup",
			Hostname:    "host2",
			TotalRuns:   5,
			Failures:    0,
			FailureRate: 0.0,
			LastStatus:  "success",
			LastRun:     testTime,
		},
	}
}

func TestCSVExporter_Header(t *testing.T) {
	var buf bytes.Buffer
	e := export.NewCSVExporter(&buf)
	if err := e.WriteStats(sampleStats()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}
	header := lines[0]
	for _, col := range []string{"job_name", "hostname", "total_runs", "failure_rate_pct", "last_status"} {
		if !strings.Contains(header, col) {
			t.Errorf("header missing column %q", col)
		}
	}
}

func TestCSVExporter_RowCount(t *testing.T) {
	var buf bytes.Buffer
	e := export.NewCSVExporter(&buf)
	_ = e.WriteStats(sampleStats())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 2 data rows
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestCSVExporter_FailureRate(t *testing.T) {
	var buf bytes.Buffer
	e := export.NewCSVExporter(&buf)
	_ = e.WriteStats(sampleStats())
	if !strings.Contains(buf.String(), "20.00") {
		t.Error("expected failure rate 20.00 in output")
	}
}

func TestCSVExporter_EmptyStats(t *testing.T) {
	var buf bytes.Buffer
	e := export.NewCSVExporter(&buf)
	if err := e.WriteStats([]model.JobStats{}); err != nil {
		t.Fatalf("unexpected error on empty stats: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
