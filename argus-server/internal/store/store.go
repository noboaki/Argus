package store

import "time"

type Metric struct {
	AgentID   string
	Hostname  string
	Timestamp time.Time
	CPUUsage  float64
	MemUsage  float64
	DiskUsage float64
}

type Store interface {
	Save(metrics Metric) error
	GetAll() []Metric
	GetByAgent(agentID string) []Metric
}
