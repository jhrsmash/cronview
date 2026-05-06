// Package model defines the core data structures used throughout cronview.
package model

import "time"

// Status represents the execution status of a cron job.
type Status int

const (
	// StatusUnknown indicates the job status could not be determined.
	StatusUnknown Status = iota
	// StatusSuccess indicates the job completed successfully.
	StatusSuccess
	// StatusFailure indicates the job failed during execution.
	StatusFailure
	// StatusStarted indicates the job was started (no completion record yet).
	StatusStarted
)

// String returns a human-readable representation of the status.
func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusFailure:
		return "failure"
	case StatusStarted:
		return "started"
	default:
		return "unknown"
	}
}

// JobEntry represents a single cron job log entry parsed from syslog.
type JobEntry struct {
	// Timestamp is when the log entry was recorded.
	Timestamp time.Time

	// Hostname is the server that ran the job.
	Hostname string

	// JobName is the identifier for the cron job (e.g., command or label).
	JobName string

	// Status is the classified outcome of this log entry.
	Status Status

	// RawMessage is the original unparsed log line.
	RawMessage string
}

// JobStats aggregates statistics for a named cron job across multiple entries.
type JobStats struct {
	// JobName is the identifier for the cron job.
	JobName string

	// Hostname is the server associated with these stats.
	Hostname string

	// TotalRuns is the total number of completed executions recorded.
	TotalRuns int

	// Failures is the count of runs that ended with a failure status.
	Failures int

	// LastRun is the timestamp of the most recent log entry for this job.
	LastRun time.Time

	// LastStatus is the status of the most recent execution.
	LastStatus Status
}

// FailureRate returns the percentage of runs that failed, from 0.0 to 100.0.
// Returns 0 if no runs have been recorded.
func (s *JobStats) FailureRate() float64 {
	if s.TotalRuns == 0 {
		return 0.0
	}
	return float64(s.Failures) / float64(s.TotalRuns) * 100.0
}

// AggregateStats computes per-job statistics from a slice of JobEntry records.
// Entries with StatusUnknown or StatusStarted do not count toward TotalRuns.
func AggregateStats(entries []JobEntry) []JobStats {
	type key struct {
		jobName  string
		hostname string
	}

	statsMap := make(map[key]*JobStats)

	for _, e := range entries {
		k := key{jobName: e.JobName, hostname: e.Hostname}
		if _, ok := statsMap[k]; !ok {
			statsMap[k] = &JobStats{
				JobName:  e.JobName,
				Hostname: e.Hostname,
			}
		}

		st := statsMap[k]

		// Update last run timestamp and status regardless of outcome.
		if e.Timestamp.After(st.LastRun) {
			st.LastRun = e.Timestamp
			st.LastStatus = e.Status
		}

		// Only count completed runs (success or failure) toward totals.
		if e.Status == StatusSuccess || e.Status == StatusFailure {
			st.TotalRuns++
			if e.Status == StatusFailure {
				st.Failures++
			}
		}
	}

	result := make([]JobStats, 0, len(statsMap))
	for _, st := range statsMap {
		result = append(result, *st)
	}
	return result
}
