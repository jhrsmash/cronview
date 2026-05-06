package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronview/internal/model"
)

const (
	colJob      = 30
	colRuns     = 6
	colFails    = 6
	colRate     = 8
	colLast     = 10
	colHostname = 20
)

// RenderJobTable writes a formatted ASCII table of job stats to w.
func RenderJobTable(w io.Writer, stats []model.JobStats) {
	writeHeader(w)
	writeSeparator(w)
	for _, s := range stats {
		writeRow(w, s)
	}
	writeSeparator(w)
	fmt.Fprintf(w, "  %d job(s) listed\n", len(stats))
}

func writeHeader(w io.Writer) {
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %-*s  %-*s\n",
		colJob, "JOB",
		colRuns, "RUNS",
		colFails, "FAILS",
		colRate, "FAIL%",
		colLast, "LAST",
		colHostname, "HOSTNAME",
	)
}

func writeSeparator(w io.Writer) {
	total := colJob + colRuns + colFails + colRate + colLast + colHostname + 10
	fmt.Fprintln(w, strings.Repeat("-", total))
}

func writeRow(w io.Writer, s model.JobStats) {
	rate := fmt.Sprintf("%.1f%%", s.FailureRate*100)
	last := s.LastStatus
	if last == "" {
		last = "unknown"
	}
	fmt.Fprintf(w, "%-*s  %-*d  %-*d  %-*s  %-*s  %-*s\n",
		colJob, truncate(s.JobName, colJob),
		colRuns, s.TotalRuns,
		colFails, s.Failures,
		colRate, rate,
		colLast, last,
		colHostname, truncate(s.Hostname, colHostname),
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
