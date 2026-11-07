package collector

import (
	"testing"
	"time"
)

func TestMemoryMetrics(t *testing.T) {
	ts := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	got := memoryMetrics(Memory{TotalBytes: 100, AvailableBytes: 40, UsedBytes: 60}, ts)

	want := map[string]float64{
		"memory_total_bytes":     100,
		"memory_available_bytes": 40,
		"memory_used_bytes":      60,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d metrics, want %d", len(got), len(want))
	}
	for _, m := range got {
		v, ok := want[m.Name]
		if !ok {
			t.Errorf("unexpected metric %q", m.Name)
			continue
		}
		if m.Value != v {
			t.Errorf("%s: got %v, want %v", m.Name, m.Value, v)
		}
		if !m.Time.Equal(ts) {
			t.Errorf("%s: got time %v, want %v", m.Name, m.Time, ts)
		}
	}
}
