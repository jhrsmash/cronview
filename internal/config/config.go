package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the top-level cronview configuration.
type Config struct {
	SyslogPath  string        `json:"syslog_path"`
	RefreshRate time.Duration `json:"refresh_rate"`
	Alert       AlertConfig   `json:"alert"`
	Export      ExportConfig  `json:"export"`
}

// AlertConfig mirrors alert threshold settings.
type AlertConfig struct {
	WarnThreshold     float64 `json:"warn_threshold"`
	CriticalThreshold float64 `json:"critical_threshold"`
	MinRuns           int     `json:"min_runs"`
}

// ExportConfig controls automatic export behaviour.
type ExportConfig struct {
	Enabled bool   `json:"enabled"`
	Format  string `json:"format"`
	Path    string `json:"path"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		SyslogPath:  "/var/log/syslog",
		RefreshRate: 30 * time.Second,
		Alert: AlertConfig{
			WarnThreshold:     0.25,
			CriticalThreshold: 0.50,
			MinRuns:           3,
		},
		Export: ExportConfig{
			Enabled: false,
			Format:  "text",
			Path:    "cronview_export.txt",
		},
	}
}

// Validate checks that the configuration values are logically consistent.
// It returns an error describing the first invalid field encountered.
func (c Config) Validate() error {
	if c.RefreshRate <= 0 {
		return fmt.Errorf("config: refresh_rate must be positive, got %v", c.RefreshRate)
	}
	if c.Alert.WarnThreshold < 0 || c.Alert.WarnThreshold > 1 {
		return fmt.Errorf("config: warn_threshold must be between 0 and 1, got %v", c.Alert.WarnThreshold)
	}
	if c.Alert.CriticalThreshold < 0 || c.Alert.CriticalThreshold > 1 {
		return fmt.Errorf("config: critical_threshold must be between 0 and 1, got %v", c.Alert.CriticalThreshold)
	}
	if c.Alert.WarnThreshold > c.Alert.CriticalThreshold {
		return fmt.Errorf("config: warn_threshold (%v) must not exceed critical_threshold (%v)", c.Alert.WarnThreshold, c.Alert.CriticalThreshold)
	}
	if c.Alert.MinRuns < 0 {
		return fmt.Errorf("config: min_runs must be non-negative, got %d", c.Alert.MinRuns)
	}
	return nil
}

// Load reads a JSON config file from path. Missing fields retain defaults.
func Load(path string) (Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save writes the config as JSON to path.
func Save(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
