package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cronview/internal/model"
	"github.com/cronview/internal/snapshot"
)

func sampleStats() []model.JobStats {
	return []model.JobStats{
		{JobName: "backup", TotalRuns: 10, Failures: 2, FailureRate: 0.2, LastStatus: "failure"},
		{JobName: "cleanup", TotalRuns: 5, Failures: 0, FailureRate: 0.0, LastStatus: "success"},
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	snap := snapshot.New(sampleStats(), "test")
	if snap.CapturedAt.Before(before) {
		t.Errorf("expected CapturedAt >= %v, got %v", before, snap.CapturedAt)
	}
}

func TestNew_SetsLabel(t *testing.T) {
	snap := snapshot.New(sampleStats(), "baseline")
	if snap.Label != "baseline" {
		t.Errorf("expected label 'baseline', got %q", snap.Label)
	}
}

func TestNew_PreservesStats(t *testing.T) {
	snap := snapshot.New(sampleStats(), "")
	if len(snap.Stats) != 2 {
		t.Fatalf("expected 2 stats, got %d", len(snap.Stats))
	}
	if snap.Stats[0].JobName != "backup" {
		t.Errorf("unexpected job name: %s", snap.Stats[0].JobName)
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New(sampleStats(), "round-trip")
	if err := snapshot.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Label != orig.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, orig.Label)
	}
	if len(loaded.Stats) != len(orig.Stats) {
		t.Errorf("stats count mismatch: got %d, want %d", len(loaded.Stats), len(orig.Stats))
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading non-existent file")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "snap.json")
	if err := snapshot.Save(path, snapshot.New(sampleStats(), "")); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
