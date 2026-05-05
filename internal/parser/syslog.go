package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// CronEntry represents a single parsed cron job log entry.
type CronEntry struct {
	Timestamp time.Time
	Hostname  string
	JobName   string
	PID       string
	Message   string
	Status    EntryStatus
}

// EntryStatus represents the outcome of a cron job execution.
type EntryStatus int

const (
	StatusUnknown EntryStatus = iota
	StatusStart
	StatusSuccess
	StatusFailure
)

// syslogLineRe matches lines like:
// Jan  2 15:04:05 hostname CRON[1234]: (user) CMD (command)
var syslogLineRe = regexp.MustCompile(
	`^(\w+\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+CRON\[(\d+)\]:\s+(.+)$`,
)

var currentYear = time.Now().Year()

// ParseSyslog reads cron entries from a syslog-format reader.
func ParseSyslog(r io.Reader) ([]CronEntry, error) {
	var entries []CronEntry
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseLine(line)
		if err != nil {
			continue // skip non-matching lines
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return entries, nil
}

func parseLine(line string) (CronEntry, error) {
	matches := syslogLineRe.FindStringSubmatch(line)
	if matches == nil {
		return CronEntry{}, fmt.Errorf("no match")
	}

	timeStr := fmt.Sprintf("%d %s", currentYear, matches[1])
	ts, err := time.Parse("2006 Jan  2 15:04:05", timeStr)
	if err != nil {
		ts, err = time.Parse("2006 Jan _2 15:04:05", timeStr)
		if err != nil {
			return CronEntry{}, fmt.Errorf("parse time: %w", err)
		}
	}

	msg := matches[4]
	status := classifyMessage(msg)

	return CronEntry{
		Timestamp: ts,
		Hostname:  matches[2],
		PID:       matches[3],
		Message:   msg,
		JobName:   extractJobName(msg),
		Status:    status,
	}, nil
}

func classifyMessage(msg string) EntryStatus {
	upper := strings.ToUpper(msg)
	switch {
	case strings.Contains(upper, "CMD"):
		return StatusStart
	case strings.Contains(upper, "EXIT CODE 0"), strings.Contains(upper, "(root) END"):
		return StatusSuccess
	case strings.Contains(upper, "EXIT CODE"), strings.Contains(upper, "ERROR"), strings.Contains(upper, "FAILED"):
		return StatusFailure
	default:
		return StatusUnknown
	}
}

func extractJobName(msg string) string {
	if idx := strings.Index(msg, "CMD ("); idx != -1 {
		cmd := msg[idx+5:]
		if end := strings.LastIndex(cmd, ")"); end != -1 {
			return strings.TrimSpace(cmd[:end])
		}
	}
	return msg
}
