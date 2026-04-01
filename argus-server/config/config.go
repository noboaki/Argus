package config

import "os"

type Config struct {
	Port string

	StoreBackend string

	// InfluxDB
	InfluxDBURL    string
	InfluxDBToken  string
	InfluxDBOrg    string
	InfluxDBBucket string

	// S3 | MinIO
	S3Bucket     string
	S3Region     string
	S3Endpoint   string
	AWSAccessKey string
	AWSSecretKey string
}

func Load() *Config {
	return &Config{
		// Server
		Port: ":" + getOrDefault("ARGUS_SERVER_PORT", "50051"),

		// Store
		StoreBackend: getOrDefault("ARGUS_STORE_BACKEND", "memory"),

		// InfluxDB
		InfluxDBURL:    os.Getenv("INFLUXDB_URL"),
		InfluxDBToken:  os.Getenv("INFLUXDB_TOKEN"),
		InfluxDBOrg:    os.Getenv("INFLUXDB_ORG"),
		InfluxDBBucket: os.Getenv("INFLUXDB_BUCKET"),

		// S3 / MinIO
		S3Bucket:     os.Getenv("AWS_BUCKET"),
		S3Region:     getOrDefault("AWS_REGION", "ap-northeast-2"),
		S3Endpoint:   os.Getenv("S3_ENDPOINT"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}

func getOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
