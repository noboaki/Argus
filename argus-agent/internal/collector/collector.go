package collector

import (
	"time"

	"github.com/noboaki/argus-agent/domain"
)

type Metrics struct {
	Timestamp time.Time

	CPUUsage  float64
	MemUsage  float64
	DiskUsage float64
}

type Collector interface {
	Collect() (*domain.ArgusMetric, error)
}
