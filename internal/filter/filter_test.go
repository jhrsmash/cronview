package filter_test

import (
	"testing"
	"time"

	"github.com/user/cronview/internal/filter"
	"github.com/user/cronview/internal/model"
)

var base = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

var sampleEntries = []model.LogEntry{
	{Hostname: "web-01", JobName: "backup", Status: model.StatusSuccess, Timestamp: base},
	{Hostname: "web-01", JobName: "cleanup", Status: model.StatusFailure, Timestamp: base.Add(time.Hour)},
	{Hostname: "db-02", JobName: "backup", Status: model.StatusFailure, Timestamp: base.Add(2 * time.Hour)},
	{Hostname: "db-02", JobName: "report", Status: model.StatusSuccess, Timestamp: base.Add(3 * time.Hour)},
}

func TestApply_ByJobName(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{JobName: "backup"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_ByHostname(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Hostname: "web-01"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_ByStatus(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Status: "failure"})
	if len(result) != 2 {
		t.Fatalf("expected 2 failure entries, got %d", len(result))
	}
}

func TestApply_BySince(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Since: base.Add(2 * time.Hour)})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries after since, got %d", len(result))
	}
}

func TestApply_Combined(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{
		Hostname: "db-02",
		Status:   "failure",
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].JobName != "backup" {
		t.Errorf("expected job 'backup', got '%s'", result[0].JobName)
	}
}

func TestApply_NoMatch(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{JobName: "nonexistent"})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestUniqueHostnames(t *testing.T) {
	hosts := filter.UniqueHostnames(sampleEntries)
	if len(hosts) != 2 {
		t.Fatalf("expected 2 unique hostnames, got %d", len(hosts))
	}
}
