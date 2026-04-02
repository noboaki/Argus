package store

import (
	"argus/proto"
	"time"
)

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

type MetricStore interface {
	Save(*proto.MetricBatch) error
	GetByAgent(agentID string) map[string][]*proto.Metric
	GetLatestMetric(agentID, metricName string) (*proto.Metric, error)
}

type AgentStore interface {
	RegisterAgent(info AgentInfo) error
	UnregisterAgent(agentID string) error
	UpdateLastSeen(agentID string) error
	GetAgents() []AgentInfo
	GetAgentById(agentID string) (*AgentInfo, error)
}
