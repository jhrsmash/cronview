package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestTotalPages_ExactMultiple(t *testing.T) {
	if got := TotalPages(50, 25); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestTotalPages_WithRemainder(t *testing.T) {
	if got := TotalPages(51, 25); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestTotalPages_ZeroItems(t *testing.T) {
	if got := TotalPages(0, 25); got != 1 {
		t.Errorf("expected 1 for empty list, got %d", got)
	}
}

func TestTotalPages_ZeroPageSize(t *testing.T) {
	if got := TotalPages(100, 0); got != 1 {
		t.Errorf("expected 1 for zero page size, got %d", got)
	}
}

func TestPageSlice_FirstPage(t *testing.T) {
	start, end := PageSlice(60, 25, 1)
	if start != 0 || end != 25 {
		t.Errorf("expected [0,25), got [%d,%d)", start, end)
	}
}

func TestPageSlice_LastPage(t *testing.T) {
	start, end := PageSlice(60, 25, 3)
	if start != 50 || end != 60 {
		t.Errorf("expected [50,60), got [%d,%d)", start, end)
	}
}

func TestPageSlice_BeyondEnd(t *testing.T) {
	start, end := PageSlice(10, 25, 5)
	if start != 10 || end != 10 {
		t.Errorf("expected empty slice, got [%d,%d)", start, end)
	}
}

func TestRenderPageInfo_Output(t *testing.T) {
	var buf bytes.Buffer
	RenderPageInfo(&buf, 2, 5, 120)
	out := buf.String()
	if !strings.Contains(out, "Page 2 of 5") {
		t.Errorf("missing page info in: %q", out)
	}
	if !strings.Contains(out, "120 total entries") {
		t.Errorf("missing total entries in: %q", out)
	}
}

func TestRenderPaginationBar_Middle(t *testing.T) {
	var buf bytes.Buffer
	RenderPaginationBar(&buf, 3, 7)
	out := buf.String()
	if !strings.Contains(out, "[< prev]") {
		t.Errorf("expected prev link in: %q", out)
	}
	if !strings.Contains(out, "[next >]") {
		t.Errorf("expected next link in: %q", out)
	}
	if !strings.Contains(out, "Page 3/7") {
		t.Errorf("expected page indicator in: %q", out)
	}
}

func TestRenderPaginationBar_FirstPage(t *testing.T) {
	var buf bytes.Buffer
	RenderPaginationBar(&buf, 1, 4)
	out := buf.String()
	if strings.Contains(out, "[< prev]") {
		t.Errorf("should not show prev on first page: %q", out)
	}
}

func TestRenderPaginationBar_LastPage(t *testing.T) {
	var buf bytes.Buffer
	RenderPaginationBar(&buf, 4, 4)
	out := buf.String()
	if strings.Contains(out, "[next >]") {
		t.Errorf("should not show next on last page: %q", out)
	}
}
