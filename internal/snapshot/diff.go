package snapshot

import "github.com/cronview/internal/model"

// Delta describes the change in a single job's failure rate between two snapshots.
type Delta struct {
	JobName     string
	OldRate     float64
	NewRate     float64
	RateDelta   float64 // positive = worsening, negative = improving
	OnlyInOld   bool
	OnlyInNew   bool
}

// Compare returns a slice of Deltas between a baseline and a current snapshot.
// Jobs present in only one snapshot are included with the appropriate flag set.
func Compare(baseline, current Snapshot) []Delta {
	oldMap := indexByJob(baseline.Stats)
	newMap := indexByJob(current.Stats)

	seen := make(map[string]bool)
	var deltas []Delta

	for name, ns := range newMap {
		seen[name] = true
		if os, ok := oldMap[name]; ok {
			deltas = append(deltas, Delta{
				JobName:   name,
				OldRate:   os.FailureRate,
				NewRate:   ns.FailureRate,
				RateDelta: ns.FailureRate - os.FailureRate,
			})
		} else {
			deltas = append(deltas, Delta{
				JobName:   name,
				NewRate:   ns.FailureRate,
				OnlyInNew: true,
			})
		}
	}

	for name, os := range oldMap {
		if !seen[name] {
			deltas = append(deltas, Delta{
				JobName:   name,
				OldRate:   os.FailureRate,
				OnlyInOld: true,
			})
		}
	}
	return deltas
}

func indexByJob(stats []model.JobStats) map[string]model.JobStats {
	m := make(map[string]model.JobStats, len(stats))
	for _, s := range stats {
		m[s.JobName] = s
	}
	return m
}
