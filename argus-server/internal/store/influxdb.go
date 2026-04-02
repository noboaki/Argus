package store

import (
	"argus/proto"
	"context"
	"fmt"
	"time"

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

func (s *InfluxDBStore) Save(batch *proto.MetricBatch) error {
	writeAPI := s.client.WriteAPIBlocking(s.org, s.bucket)

	for _, m := range batch.GetMetrics() {
		point := influxdb2.NewPointWithMeasurement("metrics").
			AddTag("agent_id", batch.AgentId).
			AddTag("hostname", batch.Hostname).
			AddField(m.Name, m.Value).
			SetTime(time.Unix(m.Timestamp, 0))

		for k, v := range m.Labels {
			point.AddTag(k, v)
		}

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			return fmt.Errorf("write point error (metric: %s): %w", m.Name, err)
		}
	}

	return nil
}

func (s *InfluxDBStore) GetByAgent(agentID string) map[string][]*proto.Metric {
	return nil
}

func (s *InfluxDBStore) GetLatestMetric(agentID, metricName string) (*proto.Metric, error) {
	return nil, nil
}
