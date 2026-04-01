package main

import (
	"log"
	"os"
	"time"

	"github.com/noboaki/argus-agent/internal/collector"
	"github.com/noboaki/argus-agent/internal/sender"
)

func main() {
	serverAddr := resolveServerAddr()
	runWithRetry(serverAddr)
}

func runWithRetry(serverAddr string) {
	for {
		s, err := sender.New(serverAddr)
		if err != nil {
			log.Printf("연결 실패, 5초 후 재시도: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("서버 연결 성공: %s", serverAddr)

		if err := run(s); err != nil {
			log.Printf("스트림 에러, 재연결: %v", err)
		}
	}
}

func run(s *sender.GRPCSender) error {
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
			return err
		}
	}
	return nil
}

func resolveServerAddr() string {
	if addr := os.Getenv("ARGUS_SERVER_ADDR"); addr != "" {
		return addr
	}
	return "localhost:50051"
}
