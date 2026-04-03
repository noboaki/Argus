package collector

import (
	"fmt"
	"time"

	"github.com/noboaki/argus-agent/domain"
	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUCollector struct{}

func (c *CPUCollector) Collect() ([]*domain.ArgusMetric, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("cpu collect error: %v", err)
	}

	return []*domain.ArgusMetric{
		domain.NewArgusMetric("cpu", cpuUsage[0]),
	}, nil
}
