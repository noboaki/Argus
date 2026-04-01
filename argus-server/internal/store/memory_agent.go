package store

import (
	"fmt"
	"sync"
	"time"
)

type MemoryAgentStore struct {
	mu     sync.RWMutex
	agents map[string]*AgentInfo
}

func NewMemoryAgentStore() *MemoryAgentStore {
	return &MemoryAgentStore{
		agents: make(map[string]*AgentInfo),
	}
}

func (s *MemoryAgentStore) RegisterAgent(info AgentInfo) error {
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

func (s *MemoryAgentStore) UnregisterAgent(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[agentID]
	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	agent.IsOnline = false
	return nil
}

func (s *MemoryAgentStore) UpdateLastSeen(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[agentID]
	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	agent.LastSeenAt = time.Now()
	return nil
}

func (s *MemoryAgentStore) GetAgents() []AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]AgentInfo, 0, len(s.agents))
	for _, a := range s.agents {
		result = append(result, *a)
	}
	return result
}

func (s *MemoryAgentStore) GetAgentById(agentID string) (*AgentInfo, error) {
	s.mu.RLock()
	defer s.mu.RLock()

	agent, ok := s.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}

	copied := *agent

	return &copied, nil
}
