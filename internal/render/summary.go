package render

import (
	"fmt"
	"io"

	"github.com/user/cronview/internal/model"
)

// RenderSummary writes a one-line summary banner to w.
func RenderSummary(w io.Writer, sum model.Summary) {
	status := "OK"
	if sum.TotalFailures > 0 {
		status = "DEGRADED"
	}
	fmt.Fprintf(w, "[%s] Jobs: %d  Runs: %d  Failures: %d  Overall fail rate: %.1f%%\n",
		status,
		sum.UniqueJobs,
		sum.TotalRuns,
		sum.TotalFailures,
		sum.OverallFailureRate*100,
	)
}
