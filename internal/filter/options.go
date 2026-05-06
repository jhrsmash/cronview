package filter

import "time"

// Options holds the filtering criteria for cron job entries.
type Options struct {
	// JobName filters entries by job name substring match (case-insensitive).
	JobName string

	// Hostname filters entries to only those from the specified host.
	Hostname string

	// Status filters entries by their classified status (e.g. "success", "failure").
	Status string

	// Since excludes entries older than this time. Zero value means no lower bound.
	Since time.Time

	// Until excludes entries newer than this time. Zero value means no upper bound.
	Until time.Time

	// MinFailureRate excludes jobs whose failure rate is below this threshold (0.0–1.0).
	MinFailureRate float64

	// Limit caps the number of entries returned. 0 means no limit.
	Limit int
}

// IsEmpty reports whether no filtering criteria have been set.
func (o Options) IsEmpty() bool {
	return o.JobName == "" &&
		o.Hostname == "" &&
		o.Status == "" &&
		o.Since.IsZero() &&
		o.Until.IsZero() &&
		o.MinFailureRate == 0 &&
		o.Limit == 0
}

// WithJobName returns a copy of Options with JobName set.
func (o Options) WithJobName(name string) Options {
	o.JobName = name
	return o
}

// WithHostname returns a copy of Options with Hostname set.
func (o Options) WithHostname(host string) Options {
	o.Hostname = host
	return o
}

// WithStatus returns a copy of Options with Status set.
func (o Options) WithStatus(status string) Options {
	o.Status = status
	return o
}

// WithSince returns a copy of Options with Since set.
func (o Options) WithSince(t time.Time) Options {
	o.Since = t
	return o
}

// WithLimit returns a copy of Options with Limit set.
func (o Options) WithLimit(n int) Options {
	o.Limit = n
	return o
}
