package processor

import "github.com/noboaki/argus-agent/domain"

type SimpleProcessor struct {
	labels domain.Labels
}

func NewSimpleProcessor(l domain.Labels) *SimpleProcessor {
	return &SimpleProcessor{
		labels: l,
	}
}

func (p *SimpleProcessor) Process(m *domain.ArgusMetric) {
	m.Labels.Merge(p.labels)
}
