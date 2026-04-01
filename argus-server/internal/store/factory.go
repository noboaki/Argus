package store

import (
	"fmt"

	"github.com/noboaki/argus-server/config"
)

func NewMetricStore(cfg *config.Config) (MetricStore, error) {
	switch cfg.StoreBackend {
	case "influxdb":
		return NewInfluxDBStore(
			cfg.InfluxDBURL,
			cfg.InfluxDBToken,
			cfg.InfluxDBOrg,
			cfg.InfluxDBBucket,
		)
	case "s3":
		return NewS3Store(
			cfg.S3Bucket,
			cfg.S3Region,
			cfg.S3Endpoint,
			cfg.AWSAccessKey,
			cfg.AWSSecretKey,
		)
	case "memory":
		return NewMemoryMetricStore(), nil
	default:
		return nil, fmt.Errorf("unknown store backend: %s", cfg.StoreBackend)
	}
}

func NewAgentStore() AgentStore {
	return NewMemoryAgentStore()
}
