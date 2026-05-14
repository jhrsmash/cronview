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
// If useColor is true, critical alerts are rendered in bold red and
// warnings in yellow. When there are no alerts, a short message is printed.
func RenderAlerts(w io.Writer, alerts []alert.Alert, useColor bool) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "No alerts triggered.")
		return
	}

	separator := strings.Repeat("-", 60)

	fmt.Fprintf(w, "%sALERTS (%d)%s\n", colorBold, len(alerts), colorReset)
	fmt.Fprintln(w, separator)

	for _, a := range alerts {
		prefix, suffix := alertColor(a.Severity, useColor)
		fmt.Fprintf(w, "%s%-10s%s  %-20s  %-16s  %.1f%%\n",
			prefix,
			"["+string(a.Severity)+"]",
			suffix,
			truncate(a.JobName, 20),
			truncate(a.Hostname, 16),
			a.FailureRate*100,
		)
	}
	fmt.Fprintln(w, separator)
}

// alertColor returns the ANSI escape prefix and suffix for a given severity
// level. If useColor is false, both values are empty strings.
func alertColor(severity alert.Severity, useColor bool) (prefix, suffix string) {
	if !useColor {
		return "", ""
	}
	switch severity {
	case alert.SeverityCritical:
		return colorRed + colorBold, colorReset
	case alert.SeverityWarn:
		return colorYellow, colorReset
	default:
		return "", ""
	}
}
