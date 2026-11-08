package collector

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dideban-project/dideban/internal/model"
)

// Loadavg holds the 1, 5 and 15 minute load averages.
type Loadavg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// ParseLoadavg parses the contents of /proc/loadavg, e.g.
// "0.42 0.50 0.61 1/512 12345".
// It takes an io.Reader (not a file path) so it can be tested with plain strings.
func ParseLoadavg(r io.Reader) (Loadavg, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Loadavg{}, fmt.Errorf("loadavg: reading: %w", err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return Loadavg{}, fmt.Errorf("loadavg: expected at least 3 fields, got %d", len(fields))
	}

	var vals [3]float64
	for i := 0; i < 3; i++ {
		v, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return Loadavg{}, fmt.Errorf("loadavg: parsing field %d: %w", i+1, err)
		}
		vals[i] = v
	}

	return Loadavg{Load1: vals[0], Load5: vals[1], Load15: vals[2]}, nil
}

// ReadLoadavg reads the current load averages from /proc/loadavg.
func ReadLoadavg() (Loadavg, error) {
	f, err := os.Open("/proc/loadavg")
	if err != nil {
		return Loadavg{}, err
	}
	defer f.Close()
	return ParseLoadavg(f)
}

// LoadavgCollector exposes ReadLoadavg through the Collector interface.
type LoadavgCollector struct{}

var _ Collector = LoadavgCollector{}

// Name implements Collector.
func (LoadavgCollector) Name() string { return "loadavg" }

// Collect implements Collector.
func (LoadavgCollector) Collect(ctx context.Context) ([]model.Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	l, err := ReadLoadavg()
	if err != nil {
		return nil, err
	}
	return loadavgMetrics(l, time.Now()), nil
}

// loadavgMetrics converts Loadavg into metrics. It is separate from Collect
// so it can be tested without reading /proc.
func loadavgMetrics(l Loadavg, ts time.Time) []model.Metric {
	return []model.Metric{
		{Name: "load1", Value: l.Load1, Time: ts},
		{Name: "load5", Value: l.Load5, Time: ts},
		{Name: "load15", Value: l.Load15, Time: ts},
	}
}
