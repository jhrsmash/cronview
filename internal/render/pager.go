package render

import (
	"fmt"
	"io"
	"strings"
)

// PageOptions controls pagination behaviour for rendered output.
type PageOptions struct {
	PageSize int
	Page     int // 1-based
}

// DefaultPageOptions returns sensible defaults.
func DefaultPageOptions() PageOptions {
	return PageOptions{
		PageSize: 25,
		Page:     1,
	}
}

// TotalPages returns the number of pages needed for itemCount items.
func TotalPages(itemCount, pageSize int) int {
	if pageSize <= 0 || itemCount == 0 {
		return 1
	}
	pages := itemCount / pageSize
	if itemCount%pageSize != 0 {
		pages++
	}
	return pages
}

// RenderPageInfo writes a human-readable pagination footer to w.
func RenderPageInfo(w io.Writer, currentPage, totalPages, totalItems int) {
	fmt.Fprintf(w, "Page %d of %d  (%d total entries)\n",
		currentPage, totalPages, totalItems)
}

// PageSlice returns the sub-slice of items corresponding to the requested page.
// Items must be a slice of any type represented as []interface{} by the caller;
// this generic helper works on index ranges only.
func PageSlice(totalItems, pageSize, page int) (start, end int) {
	if pageSize <= 0 {
		return 0, totalItems
	}
	start = (page - 1) * pageSize
	if start >= totalItems {
		start = totalItems
	}
	end = start + pageSize
	if end > totalItems {
		end = totalItems
	}
	return start, end
}

// RenderPaginationBar writes a compact ASCII navigation bar to w.
// Example:  [< prev]  Page 3/7  [next >]
func RenderPaginationBar(w io.Writer, current, total int) {
	var sb strings.Builder
	if current > 1 {
		sb.WriteString("[< prev]  ")
	} else {
		sb.WriteString("          ")
	}
	sb.WriteString(fmt.Sprintf("Page %d/%d", current, total))
	if current < total {
		sb.WriteString("  [next >]")
	}
	fmt.Fprintln(w, sb.String())
}
