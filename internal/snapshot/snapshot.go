package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cronview/internal/model"
)

// Snapshot holds a point-in-time capture of job statistics.
type Snapshot struct {
	CapturedAt time.Time          `json:"captured_at"`
	Label      string             `json:"label,omitempty"`
	Stats      []model.JobStats   `json:"stats"`
}

// Save writes the snapshot to a JSON file at the given path.
func Save(path string, snap Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (Snapshot, error) {
	var snap Snapshot
	f, err := os.Open(path)
	if err != nil {
		return snap, fmt.Errorf("snapshot: open: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return snap, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// New creates a Snapshot from a slice of JobStats with the current timestamp.
func New(stats []model.JobStats, label string) Snapshot {
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Label:      label,
		Stats:      stats,
	}
}
