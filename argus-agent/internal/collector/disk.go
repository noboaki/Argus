package collector

import (
	"fmt"
	"time"

	"github.com/noboaki/argus-agent/domain"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskCollector struct {
	Path string
}

func (d *DiskCollector) Collect() (*domain.ArgusMetric, error) {
	path := d.Path
	if path == "" {
		path = "/"
	}

	stat, err := disk.Usage(path)
	if err != nil {
		return nil, fmt.Errorf("disk collect error: %v", err)
	}

	return &domain.ArgusMetric{
		Name:      "disk",
		Value:     stat.UsedPercent,
		Timestamp: time.Now(),
	}, nil
}
