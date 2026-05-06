package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/user/cronview/internal/model"
)

// CSVExporter writes job stats to a CSV format.
type CSVExporter struct {
	w io.Writer
}

// NewCSVExporter creates a new CSVExporter writing to w.
func NewCSVExporter(w io.Writer) *CSVExporter {
	return &CSVExporter{w: w}
}

// WriteStats writes a slice of JobStats as CSV rows.
func (e *CSVExporter) WriteStats(stats []model.JobStats) error {
	cw := csv.NewWriter(e.w)

	header := []string{"job_name", "hostname", "total_runs", "failures", "failure_rate_pct", "last_status", "last_run"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("csv: write header: %w", err)
	}

	for _, s := range stats {
		lastRun := ""
		if !s.LastRun.IsZero() {
			lastRun = s.LastRun.Format(time.RFC3339)
		}
		row := []string{
			s.JobName,
			s.Hostname,
			fmt.Sprintf("%d", s.TotalRuns),
			fmt.Sprintf("%d", s.Failures),
			fmt.Sprintf("%.2f", s.FailureRate*100),
			s.LastStatus,
			lastRun,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("csv: write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
