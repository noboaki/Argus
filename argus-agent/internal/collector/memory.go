package collector

import (
	"fmt"
	"time"

	"github.com/noboaki/argus-agent/domain"
	"github.com/shirou/gopsutil/v4/mem"
)

type MemCollector struct{}

func (m *MemCollector) Collect() (*domain.ArgusMetric, error) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("memory collect error: %v", err)
	}

	return &domain.ArgusMetric{
		Name:      "memory",
		Value:     stat.UsedPercent,
		Timestamp: time.Now(),
	}, nil
}
