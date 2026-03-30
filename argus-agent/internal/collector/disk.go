package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsage struct {
	Path string
}

func (d *DiskUsage) Collect() (float64, error) {
	path := d.Path
	if path == "" {
		path = "/"
	}

	stat, err := disk.Usage(path)
	if err != nil {
		return 0, fmt.Errorf("disk collect error: %v", err)
	}
	return stat.UsedPercent, nil
}

func (d *DiskUsage) Name() string {
	return "disk"
}
