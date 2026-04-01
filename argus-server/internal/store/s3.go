package store

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	client *s3.Client
	bucket string
}

func NewS3Store(bucket, region, endpoint, accessKey, secretAccessKey string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv(accessKey),
			os.Getenv(secretAccessKey),
			"",
		)),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Store{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Store) Save(metrics Metric) error {
	return nil
}

func (s *S3Store) GetByAgent(agentID string) []Metric {
	return nil
}

func (s *S3Store) GetLatestMetric(agentID string) (*Metric, error) {
	return nil, nil
}
