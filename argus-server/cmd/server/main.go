package main

import (
	"argus/proto"
	"log"
	"net"
	"time"

	"github.com/noboaki/argus-server/config"
	grpcHandler "github.com/noboaki/argus-server/internal/grpc"
	"github.com/noboaki/argus-server/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	cfg := config.Load()

	metricStore, err := store.NewMetricStore(cfg)
	if err != nil {
		log.Fatalf("metric store init error: %v", err)
	}

	agentStore := store.NewAgentStore()

	handler := grpcHandler.NewHandler(agentStore, metricStore)

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(1000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 30 * time.Second,
		}),
	)
	proto.RegisterMetricServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Argus Server listening on %s\n", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
