package alert

import (
	"fmt"
	"strings"

	"github.com/user/cronview/internal/model"
)

// Severity represents the alert level.
type Severity string

const (
	SeverityWarn     Severity = "WARN"
	SeverityCritical Severity = "CRITICAL"
)

// Alert represents a triggered alert for a job.
type Alert struct {
	JobName      string
	Hostname     string
	FailureRate  float64
	TotalRuns    int
	Severity     Severity
	Message      string
}

func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s on %s — failure rate %.1f%% (%d runs): %s",
		a.Severity, a.JobName, a.Hostname, a.FailureRate*100, a.TotalRuns, a.Message)
}

// Config holds thresholds for alert evaluation.
type Config struct {
	WarnThreshold     float64 // e.g. 0.25 for 25%
	CriticalThreshold float64 // e.g. 0.50 for 50%
	MinRuns           int     // minimum runs before alerting
}

// DefaultConfig returns sensible default alert thresholds.
func DefaultConfig() Config {
	return Config{
		WarnThreshold:     0.25,
		CriticalThreshold: 0.50,
		MinRuns:           3,
	}
}

// Evaluate checks job stats against thresholds and returns triggered alerts.
func Evaluate(stats []model.JobStats, cfg Config) []Alert {
	var alerts []Alert
	for _, s := range stats {
		if s.TotalRuns < cfg.MinRuns {
			continue
		}
		var sev Severity
		switch {
		case s.FailureRate >= cfg.CriticalThreshold:
			sev = SeverityCritical
		case s.FailureRate >= cfg.WarnThreshold:
			sev = SeverityWarn
		default:
			continue
		}
		alerts = append(alerts, Alert{
			JobName:     s.JobName,
			Hostname:    s.Hostname,
			FailureRate: s.FailureRate,
			TotalRuns:   s.TotalRuns,
			Severity:    sev,
			Message:     buildMessage(s, sev),
		})
	}
	return alerts
}

func buildMessage(s model.JobStats, sev Severity) string {
	parts := []string{
		fmt.Sprintf("%.1f%% failure rate", s.FailureRate*100),
		fmt.Sprintf("%d/%d runs failed", s.FailedRuns, s.TotalRuns),
	}
	if sev == SeverityCritical {
		parts = append(parts, "immediate attention required")
	}
	return strings.Join(parts, "; ")
}
