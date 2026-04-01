package store

import (
	"fmt"
	"sync"
)

type MemoryMetricStore struct {
	mu      sync.RWMutex
	metrics map[string][]Metric
}

func NewMemoryMetricStore() *MemoryMetricStore {
	return &MemoryMetricStore{
		metrics: make(map[string][]Metric),
	}
}

func (s *MemoryMetricStore) Save(metric Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[metric.AgentID] = append(s.metrics[metric.AgentID], metric)
	return nil
}

func (s *MemoryMetricStore) GetByAgent(agentID string) []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Metric{}, s.metrics[agentID]...)
}

func (s *MemoryMetricStore) GetLatestMetric(agentID string) (*Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := s.metrics[agentID]
	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics for agent %s", agentID)
	}

	latest := metrics[len(metrics)-1]
	return &latest, nil
}
