package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/noboaki/argus-agent/domain"
)

type Config struct {
	Collectors      []string
	Processors      []string
	ArgusServerAddr string
	ArgusAgentID    string
	TLSEnabled      string
	TLSCAFile       string
	Interval        time.Duration
	Labels          domain.Labels
}

func Load() *Config {
	hostname, _ := os.Hostname()
	labels, err := loadLabelsFromEnv()
	if err != nil {
		panic(err)
	}

	collectors, err := loadListFromEnv("COLLECTORS")
	if err != nil {
		panic(err)
	}

	processors, err := loadListFromEnv("PROCESSORS")
	if err != nil {
		panic(err)
	}

	return &Config{
		Collectors:      collectors,
		Processors:      processors,
		ArgusServerAddr: getOrDefault("ARGUS_SERVER_ADDR", "localhost:50051"),
		ArgusAgentID:    getOrDefault("ARGUS_AGENT_ID", hostname),
		TLSEnabled:      getOrDefault("ARGUS_TLS_ENABLED", "false"),
		TLSCAFile:       os.Getenv("ARGUS_TLS_CA_FILE"),
		Interval:        loadInterval(),
		Labels:          labels,
	}
}

func getOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func loadInterval() time.Duration {
	raw := getOrDefault("INTERVAL", "5s")
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("invalid INTERVAL %q, using 5s", raw)
		return 5 * time.Second
	}
	return d
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

func loadListFromEnv(key string) ([]string, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return []string{}, nil
	}

	return strings.Split(raw, ","), nil
}
