package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cronview/internal/watch"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "syslog")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestDefaultWatcherOptions_Interval(t *testing.T) {
	opts := watch.DefaultWatcherOptions("/var/log/syslog", nil)
	if opts.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", opts.Interval)
	}
}

func TestDefaultWatcherOptions_Path(t *testing.T) {
	opts := watch.DefaultWatcherOptions("/var/log/syslog", nil)
	if opts.Path != "/var/log/syslog" {
		t.Errorf("unexpected path: %s", opts.Path)
	}
}

func TestNewFileWatcher_NonExistentFile(t *testing.T) {
	opts := watch.DefaultWatcherOptions("/nonexistent/path/syslog", nil)
	fw := watch.NewFileWatcher(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := fw.Start(ctx)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestFileWatcher_DetectsChange(t *testing.T) {
	path := writeTempFile(t, "initial content\n")

	var callCount atomic.Int32
	opts := watch.WatcherOptions{
		Path:     path,
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string) { callCount.Add(1) },
	}
	fw := watch.NewFileWatcher(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go fw.Start(ctx) //nolint:errcheck

	time.Sleep(50 * time.Millisecond)
	// Touch the file to simulate a new cron log line.
	if err := os.WriteFile(path, []byte("new content\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if callCount.Load() == 0 {
		t.Error("expected OnChange to be called at least once")
	}
}

func TestFileWatcher_NoSpuriousCalls(t *testing.T) {
	path := writeTempFile(t, "static content\n")

	var callCount atomic.Int32
	opts := watch.WatcherOptions{
		Path:     path,
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string) { callCount.Add(1) },
	}
	fw := watch.NewFileWatcher(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go fw.Start(ctx) //nolint:errcheck
	time.Sleep(180 * time.Millisecond)

	if callCount.Load() != 0 {
		t.Errorf("expected no OnChange calls for unmodified file, got %d", callCount.Load())
	}
}
