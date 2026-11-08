package collector

import (
	"strings"
	"testing"
	"time"
)

func TestParseLoadavg(t *testing.T) {
	got, err := ParseLoadavg(strings.NewReader("0.42 0.50 0.61 1/512 12345\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Loadavg{Load1: 0.42, Load5: 0.50, Load15: 0.61}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseLoadavgMalformed(t *testing.T) {
	_, err := ParseLoadavg(strings.NewReader("0.42 0.50\n"))
	if err == nil {
		t.Fatal("expected error when only 2 numbers are present")
	}
}

func TestParseLoadavgBadValue(t *testing.T) {
	_, err := ParseLoadavg(strings.NewReader("0.42 abc 0.61 1/512 12345\n"))
	if err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}

func TestParseLoadavgEmpty(t *testing.T) {
	_, err := ParseLoadavg(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestLoadavgCollectorName(t *testing.T) {
	if got := (LoadavgCollector{}).Name(); got != "loadavg" {
		t.Errorf("Name() = %q, want %q", got, "loadavg")
	}
}

func TestLoadavgMetrics(t *testing.T) {
	ts := time.Unix(1700000000, 0)
	got := loadavgMetrics(Loadavg{Load1: 1, Load5: 2, Load15: 3}, ts)

	want := map[string]float64{"load1": 1, "load5": 2, "load15": 3}
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
			t.Errorf("%s = %v, want %v", m.Name, m.Value, v)
		}
		if !m.Time.Equal(ts) {
			t.Errorf("%s time = %v, want %v", m.Name, m.Time, ts)
		}
		if m.Labels != nil {
			t.Errorf("%s has labels %v, want none", m.Name, m.Labels)
		}
	}
}
