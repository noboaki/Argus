package store

import (
	"argus/proto"
	"fmt"
	"sync"
)

type MemoryMetricStore struct {
	mu      sync.RWMutex
	metrics map[string]map[string][]*proto.Metric
}

func NewMemoryMetricStore() *MemoryMetricStore {
	return &MemoryMetricStore{
		metrics: make(map[string]map[string][]*proto.Metric),
	}
}

func (s *MemoryMetricStore) Save(batch *proto.MetricBatch) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.metrics[batch.AgentId]; !ok {
		s.metrics[batch.AgentId] = make(map[string][]*proto.Metric)
	}

	for _, m := range batch.GetMetrics() {
		s.metrics[batch.AgentId][m.Name] = append(
			s.metrics[batch.AgentId][m.Name], m,
		)
	}

	return nil
}

func (s *MemoryMetricStore) GetByAgent(agentID string) map[string][]*proto.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string][]*proto.Metric)
	for name, metrics := range s.metrics[agentID] {
		result[name] = append([]*proto.Metric{}, metrics...)
	}
	return result
}

func (s *MemoryMetricStore) GetLatestMetric(agentID, metricName string) (*proto.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics, ok := s.metrics[agentID][metricName]
	if !ok || len(metrics) == 0 {
		return nil, fmt.Errorf("no metric %s for agent %s", metricName, agentID)
	}

	return metrics[len(metrics)-1], nil
}
