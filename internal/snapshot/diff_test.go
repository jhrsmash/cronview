package snapshot_test

import (
	"testing"

	"github.com/cronview/internal/model"
	"github.com/cronview/internal/snapshot"
)

func makeSnapshot(stats []model.JobStats) snapshot.Snapshot {
	return snapshot.New(stats, "")
}

func TestCompare_RateDelta(t *testing.T) {
	old := makeSnapshot([]model.JobStats{
		{JobName: "backup", FailureRate: 0.2},
	})
	cur := makeSnapshot([]model.JobStats{
		{JobName: "backup", FailureRate: 0.5},
	})
	deltas := snapshot.Compare(old, cur)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].RateDelta < 0 {
		t.Errorf("expected positive delta for worsening job, got %f", deltas[0].RateDelta)
	}
}

func TestCompare_ImprovingJob(t *testing.T) {
	old := makeSnapshot([]model.JobStats{{JobName: "cleanup", FailureRate: 0.8}})
	cur := makeSnapshot([]model.JobStats{{JobName: "cleanup", FailureRate: 0.1}})
	deltas := snapshot.Compare(old, cur)
	if deltas[0].RateDelta >= 0 {
		t.Errorf("expected negative delta for improving job, got %f", deltas[0].RateDelta)
	}
}

func TestCompare_OnlyInNew(t *testing.T) {
	old := makeSnapshot([]model.JobStats{})
	cur := makeSnapshot([]model.JobStats{{JobName: "newjob", FailureRate: 0.3}})
	deltas := snapshot.Compare(old, cur)
	if len(deltas) != 1 || !deltas[0].OnlyInNew {
		t.Errorf("expected OnlyInNew=true for new job")
	}
}

func TestCompare_OnlyInOld(t *testing.T) {
	old := makeSnapshot([]model.JobStats{{JobName: "oldjob", FailureRate: 0.1}})
	cur := makeSnapshot([]model.JobStats{})
	deltas := snapshot.Compare(old, cur)
	if len(deltas) != 1 || !deltas[0].OnlyInOld {
		t.Errorf("expected OnlyInOld=true for removed job")
	}
}

func TestCompare_MultipleJobs(t *testing.T) {
	old := makeSnapshot([]model.JobStats{
		{JobName: "a", FailureRate: 0.1},
		{JobName: "b", FailureRate: 0.4},
	})
	cur := makeSnapshot([]model.JobStats{
		{JobName: "a", FailureRate: 0.2},
		{JobName: "c", FailureRate: 0.0},
	})
	deltas := snapshot.Compare(old, cur)
	if len(deltas) != 3 {
		t.Errorf("expected 3 deltas (a, b-old, c-new), got %d", len(deltas))
	}
}

func TestCompare_EmptyBothSnapshots(t *testing.T) {
	old := makeSnapshot([]model.JobStats{})
	cur := makeSnapshot([]model.JobStats{})
	deltas := snapshot.Compare(old, cur)
	if len(deltas) != 0 {
		t.Errorf("expected 0 deltas for empty snapshots, got %d", len(deltas))
	}
}
