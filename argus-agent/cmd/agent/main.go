package main

import (
	"log"
	"os"
	"time"

	"github.com/noboaki/argus-agent/internal/collector"
	"github.com/noboaki/argus-agent/internal/sender"
)

func main() {
	serverAddr := os.Getenv("ARGUS_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}

	s, err := sender.New(serverAddr)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	collectors := []collector.Collector{
		&collector.CPUCollector{},
		&collector.MemCollector{},
		&collector.DiskCollector{},
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Printf("Argus Agent %s started. Collecting Metrics every 5s...", s.AgentID())

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

		if err := s.Send(metrics); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}
