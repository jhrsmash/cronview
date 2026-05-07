package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/cronview/internal/model"
)

// JSONExporter writes job stats as JSON.
type JSONExporter struct {
	w io.Writer
}

// jsonRecord is the serializable form of a single job's stats.
type jsonRecord struct {
	JobName     string  `json:"job_name"`
	Hostname    string  `json:"hostname"`
	TotalRuns   int     `json:"total_runs"`
	Failures    int     `json:"failures"`
	FailureRate float64 `json:"failure_rate"`
	LastStatus  string  `json:"last_status"`
	LastSeen    string  `json:"last_seen"`
}

// jsonOutput is the top-level JSON envelope.
type jsonOutput struct {
	GeneratedAt string       `json:"generated_at"`
	Count       int          `json:"count"`
	Jobs        []jsonRecord `json:"jobs"`
}

// NewJSONExporter creates a JSONExporter that writes to w.
func NewJSONExporter(w io.Writer) *JSONExporter {
	return &JSONExporter{w: w}
}

// Write serialises stats to JSON and writes them to the underlying writer.
func (e *JSONExporter) Write(stats []model.JobStats) error {
	records := make([]jsonRecord, 0, len(stats))
	for _, s := range stats {
		records = append(records, jsonRecord{
			JobName:     s.JobName,
			Hostname:    s.Hostname,
			TotalRuns:   s.TotalRuns,
			Failures:    s.Failures,
			FailureRate: roundRate(s.FailureRate),
			LastStatus:  s.LastStatus,
			LastSeen:    s.LastSeen.Format(time.RFC3339),
		})
	}

	out := jsonOutput{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Count:       len(records),
		Jobs:        records,
	}

	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("json export: %w", err)
	}
	return nil
}

// roundRate rounds a failure rate to two decimal places.
func roundRate(r float64) float64 {
	return float64(int(r*100+0.5)) / 100
}
