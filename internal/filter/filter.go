package filter

import (
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// Options holds filtering criteria for cron job entries.
type Options struct {
	JobName  string
	Hostname string
	Status   string
	Since    time.Time
	Until    time.Time
}

// Apply filters a slice of LogEntry values according to the provided Options.
// Only entries matching ALL specified (non-zero) criteria are returned.
func Apply(entries []model.LogEntry, opts Options) []model.LogEntry {
	result := make([]model.LogEntry, 0, len(entries))
	for _, e := range entries {
		if opts.JobName != "" && !strings.EqualFold(e.JobName, opts.JobName) {
			continue
		}
		if opts.Hostname != "" && !strings.EqualFold(e.Hostname, opts.Hostname) {
			continue
		}
		if opts.Status != "" && !strings.EqualFold(string(e.Status), opts.Status) {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// UniqueHostnames returns a deduplicated, sorted list of hostnames present in entries.
func UniqueHostnames(entries []model.LogEntry) []string {
	seen := make(map[string]struct{})
	for _, e := range entries {
		seen[e.Hostname] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for h := range seen {
		out = append(out, h)
	}
	return out
}
