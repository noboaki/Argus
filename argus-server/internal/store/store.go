package store

import "time"

type AgentMetadata struct {
	AgentID  string
	Hostname string
}

type Metric struct {
	AgentMetadata
	Timestamp time.Time
	CPUUsage  float64
	MemUsage  float64
	DiskUsage float64
}

type AgentInfo struct {
	AgentMetadata
	ConnectedAt time.Time
	LastSeenAt  time.Time
	IsOnline    bool
}

type Store interface {
	Save(metrics Metric) error
	GetByAgent(agentID string) []Metric
	GetLatestMetric(agentID string) (*Metric, error)

	RegisterAgent(info AgentInfo) error
	UnregisterAgent(agentID string) error
	UpdateLastSeen(agentID string) error
	GetAgents() []AgentInfo
	GetAgentById(agentID string) (*AgentInfo, error)
}
