package model

import (
	"testing"
	"time"
)

func sampleEntries() []JobEntry {
	now := time.Now()
	return []JobEntry{
		{JobName: "backup", Status: StatusSuccess, Timestamp: now.Add(-2 * time.Hour)},
		{JobName: "backup", Status: StatusFailure, Timestamp: now.Add(-1 * time.Hour)},
		{JobName: "backup", Status: StatusFailure, Timestamp: now},
		{JobName: "cleanup", Status: StatusSuccess, Timestamp: now.Add(-30 * time.Minute)},
		{JobName: "cleanup", Status: StatusSuccess, Timestamp: now.Add(-10 * time.Minute)},
	}
}

func TestComputeJobStats_TotalRuns(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	if stats["backup"].TotalRuns != 3 {
		t.Errorf("expected 3 runs for backup, got %d", stats["backup"].TotalRuns)
	}
	if stats["cleanup"].TotalRuns != 2 {
		t.Errorf("expected 2 runs for cleanup, got %d", stats["cleanup"].TotalRuns)
	}
}

func TestComputeJobStats_FailureRate(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	expected := 2.0 / 3.0 * 100.0
	if stats["backup"].FailureRate != expected {
		t.Errorf("expected failure rate %.2f, got %.2f", expected, stats["backup"].FailureRate)
	}
	if stats["cleanup"].FailureRate != 0 {
		t.Errorf("expected 0 failure rate for cleanup, got %.2f", stats["cleanup"].FailureRate)
	}
}

func TestComputeJobStats_LastStatus(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	if stats["backup"].LastStatus != string(StatusFailure) {
		t.Errorf("expected last status failure for backup, got %s", stats["backup"].LastStatus)
	}
}

func TestComputeJobStats_LastStatus_Cleanup(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	if stats["cleanup"].LastStatus != string(StatusSuccess) {
		t.Errorf("expected last status success for cleanup, got %s", stats["cleanup"].LastStatus)
	}
}

func TestComputeSummary_MostFailing(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	summary := ComputeSummary(stats)
	if summary.MostFailing != "backup" {
		t.Errorf("expected most failing job to be backup, got %s", summary.MostFailing)
	}
}

func TestComputeSummary_Totals(t *testing.T) {
	stats := ComputeJobStats(sampleEntries())
	summary := ComputeSummary(stats)
	if summary.TotalJobs != 2 {
		t.Errorf("expected 2 total jobs, got %d", summary.TotalJobs)
	}
	if summary.TotalRuns != 5 {
		t.Errorf("expected 5 total runs, got %d", summary.TotalRuns)
	}
	if summary.TotalFailures != 2 {
		t.Errorf("expected 2 total failures, got %d", summary.TotalFailures)
	}
}

func TestComputeJobStats_Empty(t *testing.T) {
	stats := ComputeJobStats([]JobEntry{})
	if len(stats) != 0 {
		t.Errorf("expected empty stats map, got %d entries", len(stats))
	}
}

func TestComputeSummary_Empty(t *testing.T) {
	stats := ComputeJobStats([]JobEntry{})
	summary := ComputeSummary(stats)
	if summary.TotalJobs != 0 {
		t.Errorf("expected 0 total jobs, got %d", summary.TotalJobs)
	}
	if summary.TotalRuns != 0 {
		t.Errorf("expected 0 total runs, got %d", summary.TotalRuns)
	}
	if summary.MostFailing != "" {
		t.Errorf("expected empty most failing job, got %s", summary.MostFailing)
	}
}
