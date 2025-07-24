package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/Vibe-7/Grpc3/.pb" // путь к пакету сгенерированных .pb

	"google.golang.org/grpc"
)

// Server реализует pb.EchoServiceServer
type Server struct {
	pb.UnimplementedEchoServiceServer
}

// Register добавляет ваш сервис в gRPC-сервер
func Register(grpcServer *grpc.Server) {
	pb.RegisterEchoServiceServer(grpcServer, &Server{})
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received message: %s", req.Message)
	return &pb.HelloResponse{Message: "Hello back from server! Got " + req.Message}, nil
}

func (s *Server) SayHelloStream(req *pb.HelloRequest, stream pb.EchoService_SayHelloStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.HelloResponse{
			Message: fmt.Sprintf("Stream message %d for: %s", i+1, req.Message),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (s *Server) SayHelloClientStream(stream pb.EchoService_SayHelloClientStreamServer) error {
	var msgs []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{
				Message: fmt.Sprintf("Server received %d messages: %v", len(msgs), msgs),
			})
		}
		if err != nil {
			return err
		}
		msgs = append(msgs, req.Message)
	}
}

func (s *Server) SayHelloBidirectional(stream pb.EchoService_SayHelloBidirectionalServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.HelloResponse{
			Message: "Server echo: " + req.Message,
		}); err != nil {
			return err
		}
	}
}
