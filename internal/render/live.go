package render

import (
	"fmt"
	"io"
	"time"

	"github.com/user/cronview/internal/model"
)

// LiveOptions configures the live dashboard frame renderer.
type DefaultLiveOptions struct{}

// LiveFrameOptions controls what is shown in a live dashboard frame.
type LiveFrameOptions struct {
	Width       int
	PageSize    int
	ShowSpark   bool
	ShowAlerts  bool
	ShowStatus  bool
	CurrentPage int
	Filter      string
}

// DefaultLiveFrameOptions returns sensible defaults for a live frame.
func DefaultLiveOptions() LiveFrameOptions {
	return LiveFrameOptions{
		Width:       100,
		PageSize:    20,
		ShowSpark:   true,
		ShowAlerts:  true,
		ShowStatus:  true,
		CurrentPage: 1,
	}
}

// RenderLiveFrame writes a complete terminal dashboard frame to w.
func RenderLiveFrame(w io.Writer, stats []model.JobStats, opts LiveFrameOptions) {
	total := len(stats)
	tpages := TotalPages(total, opts.PageSize)
	page := clamp(opts.CurrentPage, 1, max(1, tpages))

	visible := PageSlice(stats, page, opts.PageSize)

	fmt.Fprintln(w, "╔═ cronview ═══════════════════════════════════════════════════════╗")

	RenderJobTable(w, visible)

	if opts.ShowStatus {
		barData := StatusBarData{
			TotalJobs:    total,
			FilteredJobs: len(visible),
			ActiveFilter: opts.Filter,
			LastRefresh:  time.Now(),
			Page:         page,
			TotalPages:   tpages,
		}
		barOpts := DefaultStatusBarOptions()
		barOpts.Width = opts.Width
		RenderStatusBar(w, barData, barOpts)
	}

	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════════════╝")
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
