package model

import "time"

// JobStats holds aggregated statistics for a single cron job.
type JobStats struct {
	JobName      string
	TotalRuns    int
	SuccessCount int
	FailureCount int
	LastRun      time.Time
	LastStatus   string
	FailureRate  float64
}

// SummaryStats holds overall statistics across all jobs.
type SummaryStats struct {
	TotalJobs    int
	TotalRuns    int
	TotalFailures int
	MostFailing  string
	LastUpdated  time.Time
}

// ComputeJobStats calculates per-job statistics from a slice of JobEntries.
func ComputeJobStats(entries []JobEntry) map[string]*JobStats {
	stats := make(map[string]*JobStats)

	for _, e := range entries {
		s, ok := stats[e.JobName]
		if !ok {
			s = &JobStats{JobName: e.JobName}
			stats[e.JobName] = s
		}

		s.TotalRuns++
		if e.Status == StatusSuccess {
			s.SuccessCount++
		} else if e.Status == StatusFailure {
			s.FailureCount++
		}

		if e.Timestamp.After(s.LastRun) {
			s.LastRun = e.Timestamp
			s.LastStatus = string(e.Status)
		}
	}

	for _, s := range stats {
		if s.TotalRuns > 0 {
			s.FailureRate = float64(s.FailureCount) / float64(s.TotalRuns) * 100.0
		}
	}

	return stats
}

// ComputeSummary builds a SummaryStats from the per-job stats map.
func ComputeSummary(stats map[string]*JobStats) SummaryStats {
	summary := SummaryStats{
		TotalJobs:   len(stats),
		LastUpdated: time.Now(),
	}

	var topRate float64
	for _, s := range stats {
		summary.TotalRuns += s.TotalRuns
		summary.TotalFailures += s.FailureCount
		if s.FailureRate > topRate {
			topRate = s.FailureRate
			summary.MostFailing = s.JobName
		}
	}

	return summary
}
