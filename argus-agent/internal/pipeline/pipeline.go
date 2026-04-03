package pipeline

import (
	"fmt"
	"log"

	"github.com/noboaki/argus-agent/domain"
	"github.com/noboaki/argus-agent/internal/collector"
	"github.com/noboaki/argus-agent/internal/processor"
	"github.com/noboaki/argus-agent/internal/sender"
)

type Pipeline struct {
	collectors []collector.Collector
	processors []processor.Processor
	sender     *sender.GRPCSender
}

func NewPipeline(
	collectors []string,
	processors []string,
	sender *sender.GRPCSender,
	labels domain.Labels,
) *Pipeline {
	return &Pipeline{
		collectors: buildCollectors(collectors),
		processors: buildProcessors(processors, labels),
		sender:     sender,
	}
}

func (p *Pipeline) Run() error {
	var metrics []*domain.ArgusMetric

	for _, c := range p.collectors {
		ms, err := c.Collect()
		if err != nil {
			log.Printf("collect error: %v", err)
			continue
		}
		if ms == nil { // 첫 수집 시 nil 반환 (network, disk_io)
			continue
		}

		for _, m := range ms {
			for _, proc := range p.processors {
				proc.Process(m)
			}
			metrics = append(metrics, m)
		}
	}

	if len(metrics) == 0 {
		log.Printf("no metrics collected, skipping send")
		return nil
	}

	if err := p.sender.Send(metrics); err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	return nil
}

func buildCollectors(names []string) []collector.Collector {
	var collectors []collector.Collector

	if len(names) == 0 {
		return []collector.Collector{
			&collector.CPUCollector{},
			&collector.MemCollector{},
			&collector.DiskCollector{},
			&collector.NetworkCollector{},
		}
	}

	for _, name := range names {
		switch name {
		case "cpu":
			collectors = append(collectors, &collector.CPUCollector{})
		case "memory":
			collectors = append(collectors, &collector.MemCollector{})
		case "disk":
			collectors = append(collectors, &collector.DiskCollector{})
		default:
			log.Printf("unknown collector: %s", name)
		}
	}
	return collectors
}

func buildProcessors(names []string, labels domain.Labels) []processor.Processor {
	var processors []processor.Processor

	if len(names) == 0 {
		return []processor.Processor{
			processor.NewSimpleProcessor(labels),
		}
	}

	for _, name := range names {
		switch name {
		case "simple":
			processors = append(processors, processor.NewSimpleProcessor(labels))
		default:
			log.Printf("unknown processor: %s", name)
		}
	}
	return processors
}
