package collector

import (
	"strings"
	"testing"
)

const sampleMeminfo = `MemTotal:       16061232 kB
MemFree:         6969360 kB
MemAvailable:   10957480 kB
Buffers:            2412 kB
Cached:          5381036 kB
`

func TestParseMeminfo(t *testing.T) {
	got, err := ParseMeminfo(strings.NewReader(sampleMeminfo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Memory{
		TotalBytes:     16061232 * 1024,
		AvailableBytes: 10957480 * 1024,
		UsedBytes:      (16061232 - 10957480) * 1024,
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseMeminfoMissingField(t *testing.T) {
	_, err := ParseMeminfo(strings.NewReader("MemTotal: 100 kB\n"))
	if err == nil {
		t.Fatal("expected error when MemAvailable is missing")
	}
}

func TestParseMeminfoBadValue(t *testing.T) {
	_, err := ParseMeminfo(strings.NewReader("MemTotal: abc kB\nMemAvailable: 1 kB\n"))
	if err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}
