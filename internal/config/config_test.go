package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronview/internal/config"
)

func TestDefault_SyslogPath(t *testing.T) {
	cfg := config.Default()
	if cfg.SyslogPath != "/var/log/syslog" {
		t.Errorf("expected /var/log/syslog, got %s", cfg.SyslogPath)
	}
}

func TestDefault_RefreshRate(t *testing.T) {
	cfg := config.Default()
	if cfg.RefreshRate != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.RefreshRate)
	}
}

func TestDefault_AlertThresholds(t *testing.T) {
	cfg := config.Default()
	if cfg.Alert.WarnThreshold != 0.25 {
		t.Errorf("expected warn=0.25, got %f", cfg.Alert.WarnThreshold)
	}
	if cfg.Alert.CriticalThreshold != 0.50 {
		t.Errorf("expected critical=0.50, got %f", cfg.Alert.CriticalThreshold)
	}
}

func TestLoad_NonExistentFileReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/tmp/cronview_no_such_file_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SyslogPath != "/var/log/syslog" {
		t.Errorf("expected default syslog path")
	}
}

func TestLoad_OverridesFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cronview.json")

	data := map[string]interface{}{
		"syslog_path":  "/var/log/messages",
		"refresh_rate": 60000000000, // 60s in nanoseconds
	}
	b, _ := json.Marshal(data)
	os.WriteFile(path, b, 0644)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SyslogPath != "/var/log/messages" {
		t.Errorf("expected /var/log/messages, got %s", cfg.SyslogPath)
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	orig := config.Default()
	orig.SyslogPath = "/custom/path"
	orig.Alert.MinRuns = 10

	if err := config.Save(path, orig); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.SyslogPath != orig.SyslogPath {
		t.Errorf("syslog path mismatch: %s != %s", loaded.SyslogPath, orig.SyslogPath)
	}
	if loaded.Alert.MinRuns != orig.Alert.MinRuns {
		t.Errorf("min_runs mismatch: %d != %d", loaded.Alert.MinRuns, orig.Alert.MinRuns)
	}
}
