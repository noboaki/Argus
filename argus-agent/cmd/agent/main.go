package main

import (
	"fmt"
	"log"
	"time"

	"github.com/noboaki/argus-agent/internal/collector"
)

func main() {
	cpu := &collector.CPUCollector{}
	memory := &collector.MemCollector{}
	disk := &collector.DiskUsage{}

	collectors := []collector.Collector{cpu, memory, disk}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	log.Println("Argus Agent started. Collecting Metrics every 1s...")

	for range ticker.C {
		metrics := collector.Metrics{Timestamp: time.Now()}

		for _, c := range collectors {
			val, err := c.Collect()
			if err != nil {
				log.Printf("[%s] error: %v", c.Name(), err)
				continue
			}

			switch c.Name() {
			case "cpu":
				metrics.CPUUsage = val
			case "memory":
				metrics.MemUsage = val
			case "disk":
				metrics.DiskUsage = val
			}
		}

		fmt.Printf("[%s] CPU: %.1f%%  MEM: %.1f%%  DISK: %.1f%%\n",
			metrics.Timestamp.Format("2006-01-02T15:04:05.999999-07:00"),
			metrics.CPUUsage,
			metrics.MemUsage,
			metrics.DiskUsage,
		)
	}
}
