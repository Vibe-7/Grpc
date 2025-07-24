package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "grps3/.pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedEchoServiceServer
}

// Стандартная обработка одного запроса
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received message: %s", req.Message)
	return &pb.HelloResponse{Message: "Hello back from server! Got " + req.Message}, nil
}

// Обработка стрима с сервера
func (s *Server) SayHelloStream(req *pb.HelloRequest, stream pb.EchoService_SayHelloStreamServer) error {
	log.Printf("Received stream request: %s", req.Message)

	// Отправляем 5 сообщений через 1 секунду
	for i := 0; i < 5; i++ {
		resp := &pb.HelloResponse{
			Message: fmt.Sprintf("Stream message %d for: %s", i+1, req.Message),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// Client streaming
func (s *Server) SayHelloClientStream(stream pb.EchoService_SayHelloClientStreamServer) error {
	var messages []string

	// Принимаем сообщения от клиента
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Клиент закончил, отправляем итоговый ответ
			result := fmt.Sprintf("Server received %d messages: %v", len(messages), messages)
			return stream.SendAndClose(&pb.HelloResponse{
				Message: result,
			})
		}
		if err != nil {
			return err
		}

		log.Printf("Server got message: %s", req.Message)
		messages = append(messages, req.Message)
	}
}

func ClientStreamHello(client pb.EchoServiceClient) {
	// Контекст для установки тайм-аута
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Создаём стрим
	stream, err := client.SayHelloClientStream(ctx)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// Отправляем несколько сообщений
	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("Message #%d from client", i)
		if err := stream.Send(&pb.HelloRequest{Message: msg}); err != nil {
			log.Fatalf("Error sending: %v", err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	// Закрываем и получаем ответ
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error getting final response: %v", err)
	}

	log.Printf("Client got final response: %s", resp.Message)
}

func (s *Server) SayHelloBidirectional(stream pb.EchoService_SayHelloBidirectionalServer) error {
	log.Println("Bidirectional stream started")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Bidirectional stream closed by client")
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Server got bidirectional message: %s", req.Message)

		resp := &pb.HelloResponse{
			Message: fmt.Sprintf("Server echo: %s", req.Message),
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func main() {
	// Слушаем на порту 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервер с его методами
	pb.RegisterEchoServiceServer(grpcServer, &Server{})

	defer grpcServer.Stop()

	// Запускаем сервер
	log.Println("gRPC server listening on :50051")
	func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Теперь запускаем клиент, который будет работать с сервером
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)

	// Запускаем тестовый клиентский стрим
	ClientStreamHello(client)
}
