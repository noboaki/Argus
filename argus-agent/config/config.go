package config

import (
	"encoding/json"
	"os"

	"github.com/noboaki/argus-agent/domain"
)

type Config struct {
	ArgusServerAddr string
	ArgusAgentID    string
	Interval        string
	Labels          domain.Labels
}

func Load() *Config {
	hostname, _ := os.Hostname()
	labels, err := loadLabelsFromEnv()
	if err != nil {
		panic(err)
	}

	return &Config{
		ArgusServerAddr: getOrDefault("ARGUS_SERVER_ADDR", "localhost:50051"),
		ArgusAgentID:    getOrDefault("ARGUS_AGENT_ID", hostname),
		Interval:        getOrDefault("INTERVAL", "5s"),
		Labels:          labels,
	}
}

func getOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func loadLabelsFromEnv() (domain.Labels, error) {
	raw := os.Getenv("LABELS")
	if raw == "" {
		return domain.Labels{}, nil
	}

	var labels domain.Labels
	err := json.Unmarshal([]byte(raw), &labels)
	return labels, err
}
