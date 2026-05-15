package watch_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cronview/internal/watch"
)

func TestDebouncer_CallsActionAfterDelay(t *testing.T) {
	var count atomic.Int32
	d := watch.NewDebouncer(50*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	time.Sleep(120 * time.Millisecond)

	if count.Load() != 1 {
		t.Errorf("expected 1 call, got %d", count.Load())
	}
}

func TestDebouncer_CoalesceRapidTriggers(t *testing.T) {
	var count atomic.Int32
	d := watch.NewDebouncer(60*time.Millisecond, func() { count.Add(1) })

	for i := 0; i < 5; i++ {
		d.Trigger()
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)

	if count.Load() != 1 {
		t.Errorf("expected exactly 1 coalesced call, got %d", count.Load())
	}
}

func TestDebouncer_Stop_PreventsAction(t *testing.T) {
	var count atomic.Int32
	d := watch.NewDebouncer(80*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	d.Stop()
	time.Sleep(150 * time.Millisecond)

	if count.Load() != 0 {
		t.Errorf("expected 0 calls after Stop, got %d", count.Load())
	}
}

func TestDebouncer_MultipleRoundsWork(t *testing.T) {
	var count atomic.Int32
	d := watch.NewDebouncer(40*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	time.Sleep(100 * time.Millisecond)
	d.Trigger()
	time.Sleep(100 * time.Millisecond)

	if count.Load() != 2 {
		t.Errorf("expected 2 calls across two rounds, got %d", count.Load())
	}
}

func TestNewDebouncer_NilAction_DoesNotPanic(t *testing.T) {
	d := watch.NewDebouncer(20*time.Millisecond, nil)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Trigger with nil action panicked: %v", r)
		}
	}()
	d.Trigger()
	time.Sleep(60 * time.Millisecond)
}
