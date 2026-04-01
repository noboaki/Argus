package store

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBStore struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewInfluxDBStore(url, token, org, bucket string) (*InfluxDBStore, error) {
	client := influxdb2.NewClient(url, token)
	return &InfluxDBStore{
		client: client,
		org:    org,
		bucket: bucket,
	}, nil
}

func (s *InfluxDBStore) Save(metrics Metric) error {
	writeAPI := s.client.WriteAPIBlocking(s.org, s.bucket)

	point := influxdb2.NewPointWithMeasurement("metrics").
		AddTag("agent_id", metrics.AgentID).
		AddTag("hostname", metrics.Hostname).
		AddField("cpu_usage", metrics.CPUUsage).
		AddField("mem_usage", metrics.MemUsage).
		AddField("disk_usage", metrics.DiskUsage).
		SetTime(metrics.Timestamp)

	return writeAPI.WritePoint(context.Background(), point)
}

func (s *InfluxDBStore) GetByAgent(agentID string) []Metric {
	return nil
}

func (s *InfluxDBStore) GetLatestMetric(agentID string) (*Metric, error) {
	return nil, nil
}
