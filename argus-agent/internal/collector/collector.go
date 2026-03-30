package collector

import "time"

type Metrics struct {
	Timestamp time.Time

	CPUUsage  float64
	MemUsage  float64
	DiskUsage float64
}

type Collector interface {
	Collect() (float64, error)
	Name() string
}
