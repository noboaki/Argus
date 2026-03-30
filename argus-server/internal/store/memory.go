package store

import "sync"

type MemoryStore struct {
	mu      sync.RWMutex
	metrics []Metric
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Save(metric Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics = append(s.metrics, metric)
	return nil
}

func (s *MemoryStore) GetAll() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Metric{}, s.metrics...)
}

func (s *MemoryStore) GetByAgent(agentID string) []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Metric
	for _, m := range s.metrics {
		if m.AgentID == agentID {
			result = append(result, m)
		}
	}
	return result
}
