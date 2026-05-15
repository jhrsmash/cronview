package watch

import (
	"sync"
	"time"
)

// Debouncer coalesces rapid successive calls into a single callback invocation
// after a quiet period. This prevents unnecessary re-parses when a log file is
// written in multiple small chunks.
type Debouncer struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
	action func()
}

// NewDebouncer creates a Debouncer that waits delay after the last call before
// invoking action.
func NewDebouncer(delay time.Duration, action func()) *Debouncer {
	return &Debouncer{
		delay:  delay,
		action: action,
	}
}

// Trigger schedules the debounced action. Repeated calls within the delay
// window reset the timer.
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.delay, func() {
		if d.action != nil {
			d.action()
		}
	})
}

// Stop cancels any pending invocation.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
