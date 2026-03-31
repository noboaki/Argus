package store

import (
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.RWMutex
	metrics map[string][]Metric
	agents  map[string]*AgentInfo
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string][]Metric),
		agents:  make(map[string]*AgentInfo),
	}
}

func (s *MemoryStore) Save(metric Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[metric.AgentID] = append(s.metrics[metric.AgentID], metric)
	return nil
}

func (s *MemoryStore) GetByAgent(agentID string) []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Metric{}, s.metrics[agentID]...)
}

func (s *MemoryStore) GetLatestMetric(agentID string) (*Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := s.metrics[agentID]
	if len(metrics) == 0 {
		return nil, fmt.Errorf("no metrics for agent %s", agentID)
	}

	latest := metrics[len(metrics)-1]
	return &latest, nil
}

func (s *MemoryStore) RegisterAgent(info AgentInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 이미 등록된 Agent면 재접속으로 처리 (ConnectedAt 갱신)
	if existing, ok := s.agents[info.AgentID]; ok {
		existing.IsOnline = true
		existing.ConnectedAt = info.ConnectedAt
		existing.LastSeenAt = info.ConnectedAt
		return nil
	}

	s.agents[info.AgentID] = &info
	return nil
}

func (s *MemoryStore) UnregisterAgent(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[agentID]
	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	agent.IsOnline = false
	return nil
}

func (s *MemoryStore) UpdateLastSeen(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[agentID]
	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	agent.LastSeenAt = time.Now()
	return nil
}

func (s *MemoryStore) GetAgents() []AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]AgentInfo, 0, len(s.agents))
	for _, a := range s.agents {
		result = append(result, *a)
	}
	return result
}

func (s *MemoryStore) GetAgentById(agentID string) (*AgentInfo, error) {
	s.mu.RLock()
	defer s.mu.RLock()

	agent, ok := s.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}

	copied := *agent

	return &copied, nil
}
