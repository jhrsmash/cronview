package filter

import (
	"testing"
	"time"
)

func TestOptions_IsEmpty_Default(t *testing.T) {
	var o Options
	if !o.IsEmpty() {
		t.Error("expected default Options to be empty")
	}
}

func TestOptions_IsEmpty_WithJobName(t *testing.T) {
	o := Options{}.WithJobName("backup")
	if o.IsEmpty() {
		t.Error("expected Options with JobName to not be empty")
	}
}

func TestOptions_IsEmpty_WithHostname(t *testing.T) {
	o := Options{}.WithHostname("web01")
	if o.IsEmpty() {
		t.Error("expected Options with Hostname to not be empty")
	}
}

func TestOptions_IsEmpty_WithStatus(t *testing.T) {
	o := Options{}.WithStatus("failure")
	if o.IsEmpty() {
		t.Error("expected Options with Status to not be empty")
	}
}

func TestOptions_IsEmpty_WithSince(t *testing.T) {
	o := Options{}.WithSince(time.Now())
	if o.IsEmpty() {
		t.Error("expected Options with Since to not be empty")
	}
}

func TestOptions_IsEmpty_WithLimit(t *testing.T) {
	o := Options{}.WithLimit(10)
	if o.IsEmpty() {
		t.Error("expected Options with Limit to not be empty")
	}
}

func TestOptions_WithJobName_ImmutableOriginal(t *testing.T) {
	original := Options{}
	updated := original.WithJobName("cleanup")
	if original.JobName != "" {
		t.Error("expected original Options to remain unchanged")
	}
	if updated.JobName != "cleanup" {
		t.Errorf("expected updated JobName to be 'cleanup', got %q", updated.JobName)
	}
}

func TestOptions_Chaining(t *testing.T) {
	now := time.Now()
	o := Options{}.
		WithJobName("sync").
		WithHostname("db01").
		WithStatus("success").
		WithSince(now).
		WithLimit(5)

	if o.JobName != "sync" {
		t.Errorf("JobName: got %q, want %q", o.JobName, "sync")
	}
	if o.Hostname != "db01" {
		t.Errorf("Hostname: got %q, want %q", o.Hostname, "db01")
	}
	if o.Status != "success" {
		t.Errorf("Status: got %q, want %q", o.Status, "success")
	}
	if !o.Since.Equal(now) {
		t.Errorf("Since: got %v, want %v", o.Since, now)
	}
	if o.Limit != 5 {
		t.Errorf("Limit: got %d, want %d", o.Limit, 5)
	}
	if o.IsEmpty() {
		t.Error("expected chained Options to not be empty")
	}
}
