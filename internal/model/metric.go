// Package model holds the data types shared between the agent, server and storage.
package model

import "time"

// Metric is a single measurement of one named quantity at one point in time.
type Metric struct {
	Name   string            // e.g. "memory_used_bytes"
	Value  float64           // the measured value
	Time   time.Time         // when it was measured
	Labels map[string]string // optional, e.g. {"iface": "eth0"}; may be nil
}
