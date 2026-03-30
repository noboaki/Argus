package grpc

import (
	proto "argus/proto"
	"io"
	"log"
	"time"

	"github.com/noboaki/argus-server/internal/store"
)

type Handler struct {
	proto.UnimplementedMetricServiceServer
	store store.Store
}

func NewHandler(store store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) StreamMetrics(stream proto.MetricService_StreamMetricsServer) error {
	for {
		payload, err := stream.Recv()

		// Agent가 정상적으로 Stream을 종료했을 경우
		if err == io.EOF {
			return stream.SendAndClose(&proto.Ack{
				Success: true,
				Message: "stream successfully closed.",
			})
		}

		if err != nil {
			return err
		}

		metric := store.Metric{
			AgentID:   payload.AgentId,
			Hostname:  payload.Hostname,
			Timestamp: time.Unix(payload.Timestamp, 0),
			CPUUsage:  payload.CpuUsage,
			MemUsage:  payload.MemUsage,
			DiskUsage: payload.DiskUsage,
		}

		if err := h.store.Save(metric); err != nil {
			log.Printf("store error: %v", err)
			continue
		}

		log.Printf("Agent_ID: [%s]  Hostname: [%s]  CPU: %.1f%%  MEM: %.1f%%  DISK: %.1f%%",
			payload.AgentId,
			payload.Hostname,
			payload.CpuUsage,
			payload.MemUsage,
			payload.DiskUsage,
		)
	}
}
