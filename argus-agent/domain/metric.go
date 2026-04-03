package domain

import (
	"time"
)

type ArgusMetric struct {
	Name      string
	Value     float64
	Timestamp time.Time
	Labels    Labels
}

func NewArgusMetric(name string, value float64) *ArgusMetric {
	return &ArgusMetric{
		Name:      name,
		Value:     value,
		Timestamp: time.Now(),
		Labels:    make(Labels),
	}
}

func (m *ArgusMetric) WithLabels(labels Labels) *ArgusMetric {
	m.Labels.Merge(labels)
	return m
}

func (m *ArgusMetric) WithTimestamp(t time.Time) *ArgusMetric {
	m.Timestamp = t
	return m
}
