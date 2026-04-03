package main

import (
	"log"
	"time"

	"github.com/noboaki/argus-agent/config"
	"github.com/noboaki/argus-agent/internal/pipeline"
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
	p := pipeline.NewPipeline(
		cfg.Collectors,
		cfg.Processors,
		s,
		cfg.Labels,
	)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	log.Printf("Argus Agent %s started.", s.AgentID())

	for range ticker.C {
		if err := p.Run(); err != nil {
			return err // Send 에러 시 재연결
		}
	}
	return nil
}
