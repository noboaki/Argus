package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemCollector struct{}

func (m *MemCollector) Collect() (float64, error) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("memory collect error: %v", err)
	}
	return stat.UsedPercent, nil
}

func (m *MemCollector) Name() string {
	return "memory"
}
