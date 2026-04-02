package grpc

import (
	proto "argus/proto"
	"io"
	"log"
	"time"

	"github.com/noboaki/argus-server/internal/store"
)

type Handler struct {
	proto.UnimplementedIngestionServiceServer
	agentStore  store.AgentStore
	metricStore store.MetricStore
}

func NewHandler(agentStore store.AgentStore, metricStore store.MetricStore) *Handler {
	return &Handler{agentStore: agentStore, metricStore: metricStore}
}

func (h *Handler) SendMetrics(stream proto.IngestionService_SendMetricsServer) error {
	var agentID string

	defer func() {
		if agentID != "" {
			if err := h.agentStore.UnregisterAgent(agentID); err != nil {
				log.Printf("[%s] unregister error: %v", agentID, err)
			}
			log.Printf("[%s] disconnected", agentID)
		}
	}()

	for {
		payload, err := stream.Recv()

		// Agent가 정상적으로 Stream을 종료했을 경우
		if err == io.EOF {
			return stream.SendAndClose(&proto.Ack{Success: true})
		}

		if err != nil {
			return err
		}

		if agentID == "" {
			agentID = payload.AgentId
			if err := h.agentStore.RegisterAgent(store.AgentInfo{
				AgentMetadata: store.AgentMetadata{
					AgentID:  payload.AgentId,
					Hostname: payload.Hostname,
				},
				ConnectedAt: time.Now(),
				LastSeenAt:  time.Now(),
				IsOnline:    true,
			}); err != nil {
				log.Printf("[%s] register error: %v", payload.AgentId, err)
			}
			log.Printf("[%s] connected (hostname: %s)", payload.AgentId, payload.Hostname)
		}

		if err := h.metricStore.Save(payload); err != nil {
			log.Printf("[%s] save error: %v", payload.AgentId, err)
		}

		if err := h.agentStore.UpdateLastSeen(payload.AgentId); err != nil {
			log.Printf("[%s] update last seen error: %v", payload.AgentId, err)
		}

		log.Printf("[%s] received %d metrics", payload.AgentId, len(payload.Metrics))

		for _, m := range payload.Metrics {
			log.Printf("  %-15s value=%.2f labels=%v timestamp=%s",
				m.Name,
				m.Value,
				m.Labels,
				time.Unix(m.Timestamp, 0).Format("15:04:05"),
			)
		}
	}
}
