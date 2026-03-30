package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUCollector struct{}

func (c *CPUCollector) Collect() (float64, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, fmt.Errorf("cpu collect error: %v", err)
	}
	return cpuUsage[0], nil
}

func (c *CPUCollector) Name() string {
	return "cpu"
}
