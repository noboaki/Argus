package sender

import (
	"argus/proto"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/noboaki/argus-agent/config"
	"github.com/noboaki/argus-agent/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GRPCSender struct {
	stream   proto.IngestionService_SendMetricsClient
	agentID  string
	hostname string
}

func (s *GRPCSender) Send(metrics []*domain.ArgusMetric) error {
	var protoMetrics []*proto.Metric

	for _, m := range metrics {
		protoMetrics = append(protoMetrics, &proto.Metric{
			Name:      m.Name,
			Value:     m.Value,
			Timestamp: m.Timestamp.Unix(),
			Labels:    m.Labels,
		})
	}

	payload := &proto.MetricBatch{
		AgentId:  s.AgentID(),
		Hostname: s.hostname,
		Metrics:  protoMetrics,
	}

	return s.stream.Send(payload)
}

func (s *GRPCSender) AgentID() string {
	return s.agentID
}

func New(cfg *config.Config) (*GRPCSender, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 3,
			PermitWithoutStream: true,
		}),
	}

	if cfg.TLSEnabled == "true" {
		cred, err := credentials.NewClientTLSFromFile(cfg.TLSCAFile, "")
		if err != nil {
			return nil, fmt.Errorf("TLS 로드 실패: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(
		cfg.ArgusServerAddr,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	client := proto.NewIngestionServiceClient(conn)

	stream, err := client.SendMetrics(context.Background())
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	return &GRPCSender{
		stream:   stream,
		agentID:  cfg.ArgusAgentID,
		hostname: hostname,
	}, nil
}
