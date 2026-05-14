// main.go is the entry point for the cronview terminal dashboard.
// It wires together configuration, parsing, filtering, stats computation,
// rendering, and optional export/alert evaluation.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/cronview/internal/alert"
	"github.com/yourorg/cronview/internal/config"
	"github.com/yourorg/cronview/internal/export"
	"github.com/yourorg/cronview/internal/filter"
	"github.com/yourorg/cronview/internal/model"
	"github.com/yourorg/cronview/internal/parser"
	"github.com/yourorg/cronview/internal/render"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cronview: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// ── CLI flags ────────────────────────────────────────────────────────────
	cfgPath := flag.String("config", "", "path to config file (default: ~/.config/cronview/config.json)")
	jobName := flag.String("job", "", "filter by job name (substring match)")
	hostname := flag.String("host", "", "filter by hostname")
	status := flag.String("status", "", "filter by status: success|failure|unknown")
	sinceStr := flag.String("since", "", "show entries after this time (RFC3339, e.g. 2024-01-01T00:00:00Z)")
	exportFmt := flag.String("export", "", "export format: text|csv|json (prints to stdout and exits)")
	showAlerts := flag.Bool("alerts", false, "evaluate alert thresholds and print any triggered alerts")
	pageSize := flag.Int("page-size", 20, "number of rows per page")
	page := flag.Int("page", 1, "page number to display")
	flag.Parse()

	// ── Configuration ────────────────────────────────────────────────────────
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// ── Parse syslog ─────────────────────────────────────────────────────────
	f, err := os.Open(cfg.SyslogPath)
	if err != nil {
		return fmt.Errorf("opening syslog %q: %w", cfg.SyslogPath, err)
	}
	defer f.Close()

	entries, err := parser.ParseSyslog(f)
	if err != nil {
		return fmt.Errorf("parsing syslog: %w", err)
	}

	// ── Build filter options ──────────────────────────────────────────────────
	opts := filter.Options{
		JobName:  *jobName,
		Hostname: *hostname,
		Status:   *status,
	}
	if *sinceStr != "" {
		t, err := time.Parse(time.RFC3339, *sinceStr)
		if err != nil {
			return fmt.Errorf("parsing --since value %q: %w", *sinceStr, err)
		}
		opts.Since = t
	}

	// ── Filter entries ────────────────────────────────────────────────────────
	filtered := filter.Apply(entries, opts)

	// ── Compute stats ─────────────────────────────────────────────────────────
	stats := model.ComputeJobStats(filtered)
	summary := model.ComputeSummary(stats)

	// ── Export mode (non-interactive) ─────────────────────────────────────────
	if *exportFmt != "" {
		fmt := export.ParseFormat(*exportFmt)
		if err := export.Write(os.Stdout, fmt, stats); err != nil {
			return fmt.Errorf("exporting: %w", err)
		}
		return nil
	}

	// ── Alert evaluation ──────────────────────────────────────────────────────
	if *showAlerts {
		alertCfg := alert.DefaultConfig()
		// Override thresholds from loaded config if present.
		if cfg.AlertThresholds.WarnFailureRate > 0 {
			alertCfg.WarnFailureRate = cfg.AlertThresholds.WarnFailureRate
		}
		if cfg.AlertThresholds.CriticalFailureRate > 0 {
			alertCfg.CriticalFailureRate = cfg.AlertThresholds.CriticalFailureRate
		}
		alerts := alert.Evaluate(stats, alertCfg)
		render.RenderAlerts(os.Stdout, alerts)
		return nil
	}

	// ── Paginate ──────────────────────────────────────────────────────────────
	pageOpts := render.DefaultPageOptions()
	pageOpts.PageSize = *pageSize
	pageOpts.CurrentPage = *page

	totalPages := render.TotalPages(len(stats), pageOpts.PageSize)
	paged := render.PageSlice(stats, pageOpts)

	// ── Render ────────────────────────────────────────────────────────────────
	render.RenderSummary(os.Stdout, summary)
	fmt.Fprintln(os.Stdout)
	render.RenderJobTable(os.Stdout, paged)
	fmt.Fprintln(os.Stdout)
	render.RenderPaginationBar(os.Stdout, pageOpts.CurrentPage, totalPages)

	return nil
}
