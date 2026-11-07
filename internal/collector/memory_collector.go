package collector

import (
	"context"
	"time"

	"github.com/dideban-project/dideban/internal/model"
)

// MemoryCollector exposes ReadMemory through the Collector interface.
type MemoryCollector struct{}

var _ Collector = MemoryCollector{}

// Name implements Collector.
func (MemoryCollector) Name() string { return "memory" }

// Collect implements Collector.
func (MemoryCollector) Collect(ctx context.Context) ([]model.Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m, err := ReadMemory()
	if err != nil {
		return nil, err
	}
	return memoryMetrics(m, time.Now()), nil
}

// memoryMetrics converts Memory into metrics. It is separate from Collect
// so it can be tested without reading /proc.
func memoryMetrics(m Memory, ts time.Time) []model.Metric {
	return []model.Metric{
		{Name: "memory_total_bytes", Value: float64(m.TotalBytes), Time: ts},
		{Name: "memory_available_bytes", Value: float64(m.AvailableBytes), Time: ts},
		{Name: "memory_used_bytes", Value: float64(m.UsedBytes), Time: ts},
	}
}
