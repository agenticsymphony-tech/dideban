package collector

import (
	"context"

	"github.com/dideban-project/dideban/internal/model"
)

// Collector gathers one family of metrics from the local machine.
type Collector interface {
	// Name returns a short stable name, e.g. "memory".
	Name() string
	// Collect reads the current values and returns them as metrics.
	Collect(ctx context.Context) ([]model.Metric, error)
}
