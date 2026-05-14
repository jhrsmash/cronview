package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/cronview/internal/model"
)

// LiveOptions controls the behaviour of the live/auto-refresh view.
type LiveOptions struct {
	// RefreshRate is how often the dashboard refreshes.
	RefreshRate time.Duration
	// Width is the terminal column width used for rendering.
	Width int
	// ShowSparkline controls whether the per-job sparkline row is rendered.
	ShowSparkline bool
	// ShowHistogram controls whether the failure histogram block is rendered.
	ShowHistogram bool
}

// DefaultLiveOptions returns sensible defaults for the live view.
func DefaultLiveOptions() LiveOptions {
	return LiveOptions{
		RefreshRate:   5 * time.Second,
		Width:         120,
		ShowSparkline: true,
		ShowHistogram: false,
	}
}

// LiveFrame holds all data required to render a single refresh frame.
type LiveFrame struct {
	Stats     []model.JobStats
	Summary   model.SummaryStats
	Page      int
	PageSize  int
	RenderedAt time.Time
}

// RenderLiveFrame writes a complete dashboard frame to w.
// It clears the previous content by printing ANSI escape codes so the output
// appears to update in-place when written to a real terminal.
func RenderLiveFrame(w io.Writer, frame LiveFrame, opts LiveOptions) error {
	// Move cursor to top-left and clear screen.
	fmt.Fprint(w, "\033[H\033[2J")

	// ── Header bar ──────────────────────────────────────────────────────────
	header := fmt.Sprintf(" cronview  │  %s  │  refresh every %s",
		frame.RenderedAt.Format("2006-01-02 15:04:05"),
		opts.RefreshRate,
	)
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("─", clamp(opts.Width, 40, 240)))

	// ── Summary block ───────────────────────────────────────────────────────
	if err := RenderSummary(w, frame.Summary); err != nil {
		return fmt.Errorf("live frame summary: %w", err)
	}
	fmt.Fprintln(w)

	// ── Job table (paginated) ────────────────────────────────────────────────
	pageSize := frame.PageSize
	if pageSize <= 0 {
		pageSize = DefaultPageOptions().PageSize
	}
	slice := PageSlice(frame.Stats, frame.Page, pageSize)
	if err := RenderJobTable(w, slice); err != nil {
		return fmt.Errorf("live frame table: %w", err)
	}

	// ── Pagination bar ───────────────────────────────────────────────────────
	total := TotalPages(len(frame.Stats), pageSize)
	fmt.Fprintln(w, RenderPaginationBar(frame.Page, total))

	// ── Optional sparkline ───────────────────────────────────────────────────
	if opts.ShowSparkline && len(frame.Stats) > 0 {
		fmt.Fprintln(w)
		sparkOpts := DefaultSparklineOptions()
		// Collect all entries from the current page for the sparkline.
		var entries []model.LogEntry
		for _, s := range slice {
			entries = append(entries, s.Entries...)
		}
		line := RenderSparkline(entries, sparkOpts)
		fmt.Fprintf(w, " trend: %s\n", line)
	}

	// ── Optional histogram ───────────────────────────────────────────────────
	if opts.ShowHistogram && len(frame.Stats) > 0 {
		fmt.Fprintln(w)
		histOpts := DefaultHistogramOptions()
		var entries []model.LogEntry
		for _, s := range slice {
			entries = append(entries, s.Entries...)
		}
		if err := RenderHistogram(w, entries, histOpts); err != nil {
			return fmt.Errorf("live frame histogram: %w", err)
		}
	}

	// ── Footer ───────────────────────────────────────────────────────────────
	fmt.Fprintln(w, strings.Repeat("─", clamp(opts.Width, 40, 240)))
	fmt.Fprintln(w, " q quit  ←/→ page  s sparkline  h histogram")

	return nil
}

// clamp returns v constrained to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
