package server

import (
	"argus/proto"
	"log"
	"net"

	grpcHandler "github.com/noboaki/argus-server/internal/grpc"
	"github.com/noboaki/argus-server/internal/store"
	"google.golang.org/grpc"
)

func main() {
	port := ":50051"
	s := store.NewMemoryStore()

	handler := grpcHandler.NewHandler(s)

	grpcServer := grpc.NewServer()
	proto.RegisterMetricServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Argus Server listening on %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
