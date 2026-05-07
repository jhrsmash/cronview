package export

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/user/cronview/internal/model"
)

func jsonSampleStats() []model.JobStats {
	return []model.JobStats{
		{
			JobName:     "backup",
			Hostname:    "web-01",
			TotalRuns:   10,
			Failures:    2,
			FailureRate: 0.20,
			LastStatus:  "success",
			LastSeen:    time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			JobName:     "cleanup",
			Hostname:    "db-01",
			TotalRuns:   5,
			Failures:    5,
			FailureRate: 1.00,
			LastStatus:  "failure",
			LastSeen:    time.Date(2024, 6, 2, 8, 0, 0, 0, time.UTC),
		},
	}
}

func TestJSONExporter_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	exporter := NewJSONExporter(&buf)
	if err := exporter.Write(jsonSampleStats()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestJSONExporter_Count(t *testing.T) {
	var buf bytes.Buffer
	_ = NewJSONExporter(&buf).Write(jsonSampleStats())

	var out jsonOutput
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Count != 2 {
		t.Errorf("expected count 2, got %d", out.Count)
	}
}

func TestJSONExporter_JobFields(t *testing.T) {
	var buf bytes.Buffer
	_ = NewJSONExporter(&buf).Write(jsonSampleStats())

	var out jsonOutput
	_ = json.Unmarshal(buf.Bytes(), &out)

	if out.Jobs[0].JobName != "backup" {
		t.Errorf("expected job_name 'backup', got %q", out.Jobs[0].JobName)
	}
	if out.Jobs[1].LastStatus != "failure" {
		t.Errorf("expected last_status 'failure', got %q", out.Jobs[1].LastStatus)
	}
}

func TestJSONExporter_FailureRate(t *testing.T) {
	var buf bytes.Buffer
	_ = NewJSONExporter(&buf).Write(jsonSampleStats())

	var out jsonOutput
	_ = json.Unmarshal(buf.Bytes(), &out)

	if out.Jobs[0].FailureRate != 0.20 {
		t.Errorf("expected failure_rate 0.20, got %v", out.Jobs[0].FailureRate)
	}
}

func TestJSONExporter_EmptyStats(t *testing.T) {
	var buf bytes.Buffer
	if err := NewJSONExporter(&buf).Write([]model.JobStats{}); err != nil {
		t.Fatalf("unexpected error on empty stats: %v", err)
	}
	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if out.Count != 0 {
		t.Errorf("expected count 0, got %d", out.Count)
	}
}
