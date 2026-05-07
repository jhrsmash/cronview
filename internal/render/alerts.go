package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronview/internal/alert"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

// RenderAlerts writes a formatted alert summary to w.
func RenderAlerts(w io.Writer, alerts []alert.Alert, useColor bool) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "No alerts triggered.")
		return
	}

	fmt.Fprintf(w, "%sALERTS (%d)%s\n", colorBold, len(alerts), colorReset)
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, a := range alerts {
		prefix, suffix := "", ""
		if useColor {
			switch a.Severity {
			case alert.SeverityCritical:
				prefix = colorRed + colorBold
				suffix = colorReset
			case alert.SeverityWarn:
				prefix = colorYellow
				suffix = colorReset
			}
		}
		fmt.Fprintf(w, "%s%-10s%s  %-20s  %-16s  %.1f%%\n",
			prefix,
			"["+string(a.Severity)+"]",
			suffix,
			truncate(a.JobName, 20),
			truncate(a.Hostname, 16),
			a.FailureRate*100,
		)
	}
	fmt.Fprintln(w, strings.Repeat("-", 60))
}
