package model

import "time"

// Status represents the outcome of a cron job execution.
type Status int

const (
	StatusUnknown Status = iota
	StatusSuccess
	StatusFailure
)

// String returns a human-readable label for a Status.
func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusFailure:
		return "failure"
	default:
		return "unknown"
	}
}

// LogEntry represents a single parsed cron log line.
type LogEntry struct {
	Timestamp time.Time
	Hostname  string
	JobName   string
	Message   string
	Status    Status
}

// IsFailure returns true if the entry represents a failed job run.
func (e LogEntry) IsFailure() bool {
	return e.Status == StatusFailure
}

// IsSuccess returns true if the entry represents a successful job run.
func (e LogEntry) IsSuccess() bool {
	return e.Status == StatusSuccess
}
