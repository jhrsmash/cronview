package model

import "time"

// Status values for log entries.
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusUnknown = "unknown"
)

// LogEntry represents a single parsed line from syslog relating to a cron job.
type LogEntry struct {
	Timestamp time.Time
	Hostname  string
	JobName   string
	Status    string
	Message   string
}

// JobStats holds aggregated statistics for a single (job, host) pair.
type JobStats struct {
	JobName     string
	Hostname    string
	TotalRuns   int
	Failures    int
	FailureRate float64
	LastStatus  string
	LastRun     time.Time
	Entries     []LogEntry
}

// AggregateStats groups a slice of LogEntry values by (JobName, Hostname) and
// returns a map keyed by "jobname@hostname" for downstream consumers.
func AggregateStats(entries []LogEntry) map[string][]LogEntry {
	result := make(map[string][]LogEntry)
	for _, e := range entries {
		key := e.JobName + "@" + e.Hostname
		result[key] = append(result[key], e)
	}
	return result
}
