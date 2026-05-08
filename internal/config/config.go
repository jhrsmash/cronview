package config

import (
	"encoding/json"
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
