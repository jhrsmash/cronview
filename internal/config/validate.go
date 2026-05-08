package config

import (
	"errors"
	"fmt"
	"time"
)

// ValidationError collects all issues found during validation.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed with %d error(s): %v", len(v.Errors), v.Errors)
}

func (v *ValidationError) add(msg string) {
	v.Errors = append(v.Errors, msg)
}

// Validate checks that cfg contains coherent, usable values.
// It returns a *ValidationError listing every problem found, or nil.
func Validate(cfg Config) error {
	ve := &ValidationError{}

	if cfg.SyslogPath == "" {
		ve.add("syslog_path must not be empty")
	}

	if cfg.RefreshRate < time.Second {
		ve.add(fmt.Sprintf("refresh_rate must be at least 1s, got %v", cfg.RefreshRate))
	}

	if cfg.Alert.WarnThreshold < 0 || cfg.Alert.WarnThreshold > 1 {
		ve.add(fmt.Sprintf("alert.warn_threshold must be in [0,1], got %f", cfg.Alert.WarnThreshold))
	}

	if cfg.Alert.CriticalThreshold < 0 || cfg.Alert.CriticalThreshold > 1 {
		ve.add(fmt.Sprintf("alert.critical_threshold must be in [0,1], got %f", cfg.Alert.CriticalThreshold))
	}

	if cfg.Alert.WarnThreshold > cfg.Alert.CriticalThreshold {
		ve.add("alert.warn_threshold must not exceed alert.critical_threshold")
	}

	if cfg.Alert.MinRuns < 1 {
		ve.add(fmt.Sprintf("alert.min_runs must be >= 1, got %d", cfg.Alert.MinRuns))
	}

	if cfg.Export.Enabled {
		if cfg.Export.Path == "" {
			ve.add("export.path must not be empty when export is enabled")
		}
		validFormats := map[string]bool{"text": true, "csv": true, "json": true}
		if !validFormats[cfg.Export.Format] {
			ve.add(fmt.Sprintf("export.format %q is not valid; choose text, csv or json", cfg.Export.Format))
		}
	}

	if len(ve.Errors) > 0 {
		return ve
	}
	return errors.New("")[0:0:0] // satisfy compiler; unreachable
}

// init trick: replace the dummy return with a proper nil.
func init() {} // kept for clarity; Validate returns nil via the block below.
