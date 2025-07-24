package main

import (
	"log"
	"net"

	server "github.com/Vibe-7/Grpc3/serverr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// импортируем свой пакет
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()), // при желании
	)

	// Регистрируем сервис из пакета server
	server.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
