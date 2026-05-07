package alert_test

import (
	"testing"

	"github.com/user/cronview/internal/alert"
	"github.com/user/cronview/internal/model"
)

func sampleStats() []model.JobStats {
	return []model.JobStats{
		{JobName: "backup", Hostname: "srv1", TotalRuns: 10, FailedRuns: 6, FailureRate: 0.60},
		{JobName: "cleanup", Hostname: "srv1", TotalRuns: 8, FailedRuns: 2, FailureRate: 0.25},
		{JobName: "healthcheck", Hostname: "srv2", TotalRuns: 20, FailedRuns: 1, FailureRate: 0.05},
		{JobName: "sync", Hostname: "srv2", TotalRuns: 2, FailedRuns: 2, FailureRate: 1.0}, // below MinRuns
	}
}

func TestEvaluate_CriticalAlert(t *testing.T) {
	alerts := alert.Evaluate(sampleStats(), alert.DefaultConfig())
	var found bool
	for _, a := range alerts {
		if a.JobName == "backup" && a.Severity == alert.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected CRITICAL alert for 'backup'")
	}
}

func TestEvaluate_WarnAlert(t *testing.T) {
	alerts := alert.Evaluate(sampleStats(), alert.DefaultConfig())
	var found bool
	for _, a := range alerts {
		if a.JobName == "cleanup" && a.Severity == alert.SeverityWarn {
			found = true
		}
	}
	if !found {
		t.Error("expected WARN alert for 'cleanup'")
	}
}

func TestEvaluate_NoAlertBelowThreshold(t *testing.T) {
	alerts := alert.Evaluate(sampleStats(), alert.DefaultConfig())
	for _, a := range alerts {
		if a.JobName == "healthcheck" {
			t.Errorf("unexpected alert for 'healthcheck': %v", a)
		}
	}
}

func TestEvaluate_SkipsBelowMinRuns(t *testing.T) {
	alerts := alert.Evaluate(sampleStats(), alert.DefaultConfig())
	for _, a := range alerts {
		if a.JobName == "sync" {
			t.Errorf("expected 'sync' to be skipped due to MinRuns, got alert: %v", a)
		}
	}
}

func TestEvaluate_Count(t *testing.T) {
	alerts := alert.Evaluate(sampleStats(), alert.DefaultConfig())
	if len(alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestAlert_String_ContainsSeverity(t *testing.T) {
	a := alert.Alert{
		JobName:     "backup",
		Hostname:    "srv1",
		FailureRate: 0.60,
		TotalRuns:   10,
		Severity:    alert.SeverityCritical,
		Message:     "60.0% failure rate",
	}
	s := a.String()
	if len(s) == 0 {
		t.Error("expected non-empty string")
	}
	if s[:10] != "[CRITICAL]" {
		t.Errorf("expected string to start with [CRITICAL], got: %s", s)
	}
}
