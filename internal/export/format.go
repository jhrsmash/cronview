package export

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronview/internal/model"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatText Format = "text"
)

// ParseFormat converts a string to a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case string(FormatCSV):
		return FormatCSV, nil
	case string(FormatText):
		return FormatText, nil
	default:
		return "", fmt.Errorf("export: unknown format %q (supported: csv, text)", s)
	}
}

// Write exports stats in the given format to w.
func Write(w io.Writer, stats []model.JobStats, format Format) error {
	switch format {
	case FormatCSV:
		return NewCSVExporter(w).WriteStats(stats)
	case FormatText:
		return writeText(w, stats)
	default:
		return fmt.Errorf("export: unsupported format %q", format)
	}
}

func writeText(w io.Writer, stats []model.JobStats) error {
	for _, s := range stats {
		line := fmt.Sprintf("%-30s %-15s runs=%-4d failures=%-4d rate=%.1f%% status=%s\n",
			s.JobName, s.Hostname, s.TotalRuns, s.Failures, s.FailureRate*100, s.LastStatus)
		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("export: write text: %w", err)
		}
	}
	return nil
}
