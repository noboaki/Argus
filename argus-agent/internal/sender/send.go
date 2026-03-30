package sender

import (
	"argus/proto"
	"context"
	"os"

	"github.com/noboaki/argus-agent/internal/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSender struct {
	stream   proto.MetricService_StreamMetricsClient
	agentID  string
	hostname string
}

func New(serverAddr string) (*GRPCSender, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewMetricServiceClient(conn)

	stream, err := client.StreamMetrics(context.Background())
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	return &GRPCSender{
		stream:   stream,
		agentID:  hostname,
		hostname: hostname,
	}, nil
}

func (s *GRPCSender) Send(metrics collector.Metrics) error {
	payload := &proto.MetricPayload{
		AgentId:   s.agentID,
		Hostname:  s.hostname,
		Timestamp: metrics.Timestamp.Unix(),
		CpuUsage:  metrics.CPUUsage,
		MemUsage:  metrics.MemUsage,
		DiskUsage: metrics.DiskUsage,
	}
	return s.stream.Send(payload)
}
