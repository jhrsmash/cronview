package watch

import (
	"context"
	"os"
	"time"
)

// FileWatcher monitors a file for changes by polling its modification time.
type FileWatcher struct {
	path     string
	interval time.Duration
	lastMod  time.Time
	OnChange func(path string)
}

// WatcherOptions configures a FileWatcher.
type WatcherOptions struct {
	Path     string
	Interval time.Duration
	OnChange func(path string)
}

// DefaultWatcherOptions returns sensible defaults for watching a syslog file.
func DefaultWatcherOptions(path string, onChange func(string)) WatcherOptions {
	return WatcherOptions{
		Path:     path,
		Interval: 5 * time.Second,
		OnChange: onChange,
	}
}

// NewFileWatcher creates a FileWatcher from the given options.
func NewFileWatcher(opts WatcherOptions) *FileWatcher {
	return &FileWatcher{
		path:     opts.Path,
		interval: opts.Interval,
		OnChange: opts.OnChange,
	}
}

// Start begins polling the file for changes until ctx is cancelled.
func (fw *FileWatcher) Start(ctx context.Context) error {
	info, err := os.Stat(fw.path)
	if err != nil {
		return err
	}
	fw.lastMod = info.ModTime()

	ticker := time.NewTicker(fw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if fw.poll() && fw.OnChange != nil {
				fw.OnChange(fw.path)
			}
		}
	}
}

// poll checks if the file has been modified since the last check.
func (fw *FileWatcher) poll() bool {
	info, err := os.Stat(fw.path)
	if err != nil {
		return false
	}
	if info.ModTime().After(fw.lastMod) {
		fw.lastMod = info.ModTime()
		return true
	}
	return false
}
