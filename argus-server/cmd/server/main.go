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
	"google.golang.org/grpc/credentials"
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

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 30 * time.Second,
		}),
	}

	if cfg.TLSEnabled == "true" {
		cred, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatalf("TLS 로드 실패: %v", err)
		}

		opts = append(opts, grpc.Creds(cred))
		log.Println("TLS 활성화")
	} else {
		log.Println("TLS 비활성화 (insecure)")
	}

	grpcServer := grpc.NewServer(
		opts...,
	)
	proto.RegisterIngestionServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Argus Server listening on %s\n", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
