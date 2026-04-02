package main

import (
	"log"
	"time"

	"github.com/noboaki/argus-agent/config"
	"github.com/noboaki/argus-agent/domain"
	"github.com/noboaki/argus-agent/internal/collector"
	"github.com/noboaki/argus-agent/internal/processor"
	"github.com/noboaki/argus-agent/internal/sender"
)

func main() {
	runWithRetry()
}

func runWithRetry() {
	cfg := config.Load()

	for {
		s, err := sender.New(cfg)
		if err != nil {
			log.Printf("연결 실패, 5초 후 재시도: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("서버 연결 성공: %s", cfg.ArgusServerAddr)

		if err := run(s, cfg); err != nil {
			log.Printf("스트림 에러, 재연결: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func run(s *sender.GRPCSender, cfg *config.Config) error {
	collectors := []collector.Collector{
		&collector.CPUCollector{},
		&collector.MemCollector{},
		&collector.DiskCollector{},
	}

	processors := []processor.Processor{
		processor.NewSimpleProcessor(cfg.Labels),
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Printf("Argus Agent %s started. Collecting Metrics every 5s...", s.AgentID())

	for range ticker.C {
		var metrics []*domain.ArgusMetric

		for _, c := range collectors {
			m, err := c.Collect()
			if err != nil {
				log.Printf("[%s] error: %v", m.Name, err)
				continue
			}

			for _, p := range processors {
				p.Process(m)
			}

			metrics = append(metrics, m)
		}

		if err := s.Send(metrics); err != nil {
			log.Printf("Send 에러 상세: %v (type: %T)", err, err)
			return err
		}
	}
	return nil
}
