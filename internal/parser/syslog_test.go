package parser

import (
	"strings"
	"testing"
)

const sampleLog = `Jan  2 15:04:05 myhost CRON[1234]: (root) CMD (/usr/bin/backup.sh)
Jan  2 15:04:10 myhost CRON[1234]: (root) END (/usr/bin/backup.sh)
Jan  3 08:00:01 myhost CRON[5678]: (deploy) CMD (/opt/cleanup.sh)
Jan  3 08:00:05 myhost CRON[5678]: (deploy) EXIT CODE 1 (/opt/cleanup.sh)
Jan  3 09:15:00 myhost CRON[9999]: (www-data) CMD (/var/www/cron.php)
`

func TestParseSyslog_Count(t *testing.T) {
	entries, err := ParseSyslog(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(entries))
	}
}

func TestParseSyslog_Hostname(t *testing.T) {
	entries, _ := ParseSyslog(strings.NewReader(sampleLog))
	for _, e := range entries {
		if e.Hostname != "myhost" {
			t.Errorf("expected hostname 'myhost', got %q", e.Hostname)
		}
	}
}

func TestParseSyslog_StatusClassification(t *testing.T) {
	entries, _ := ParseSyslog(strings.NewReader(sampleLog))

	tests := []struct {
		index    int
		wantStatus EntryStatus
	}{
		{0, StatusStart},
		{1, StatusSuccess},
		{2, StatusStart},
		{3, StatusFailure},
		{4, StatusStart},
	}

	for _, tt := range tests {
		got := entries[tt.index].Status
		if got != tt.wantStatus {
			t.Errorf("entry[%d]: expected status %v, got %v", tt.index, tt.wantStatus, got)
		}
	}
}

func TestParseSyslog_JobName(t *testing.T) {
	entries, _ := ParseSyslog(strings.NewReader(sampleLog))

	expected := "/usr/bin/backup.sh"
	if entries[0].JobName != expected {
		t.Errorf("expected job name %q, got %q", expected, entries[0].JobName)
	}
}

func TestParseSyslog_EmptyInput(t *testing.T) {
	entries, err := ParseSyslog(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseSyslog_NonCronLinesSkipped(t *testing.T) {
	input := `Jan  2 15:04:05 myhost sshd[999]: Accepted password for user
Jan  2 15:04:06 myhost CRON[1]: (root) CMD (/bin/true)
`
	entries, _ := ParseSyslog(strings.NewReader(input))
	if len(entries) != 1 {
		t.Errorf("expected 1 cron entry, got %d", len(entries))
	}
}
