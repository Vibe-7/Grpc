package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "grps3/.pb" // Путь до вашего сгенерированного gRPC пакета

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func OneSay(client pb.EchoServiceClient) {
	// Контекст с таймаутом для одного запроса
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Вызываем SayHello для одного ответа
	response, err := client.SayHello(ctx, &pb.HelloRequest{Message: "Hello from client!"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// Печатаем ответ от сервера
	log.Printf("Server response: %s", response.GetMessage())
}

func StreamHello(client pb.EchoServiceClient) {
	// Создаём контекст с таймаутом 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Создаём запрос
	req := &pb.HelloRequest{Message: "Hello streaming!"}

	// Запускаем потоковую RPC функцию
	stream, err := client.SayHelloStream(ctx, req)
	if err != nil {
		log.Fatalf("Error starting stream: %v", err)
	}

	// Получаем ответы от сервера
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // Поток завершён
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		log.Printf("Got stream message: %s", resp.Message)
	}
}

func ClientStreamHello(client pb.EchoServiceClient) {
	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Создаём поток для отправки сообщений на сервер
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
		time.Sleep(time.Millisecond * 1500)
	}

	// Закрываем поток и получаем итоговый ответ
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error getting final response: %v", err)
	}

	log.Printf("Client got final response: %s", resp.Message)
}

func BidirectionalStreamHello(client pb.EchoServiceClient) {
	// Создаём поток для двустороннего общения
	stream, err := client.SayHelloBidirectional(context.Background())
	if err != nil {
		log.Fatalf("Error creating bidirectional stream: %v", err)
	}

	// Отправляем несколько сообщений и получаем ответы
	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("Message #%d from client", i)
		if err := stream.Send(&pb.HelloRequest{Message: msg}); err != nil {
			log.Fatalf("Error sending: %v", err)
		}

		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		log.Printf("Got bidirectional message: %s", resp.Message)
	}

	// Закрываем поток
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}
}

func main() {
	// Подключаемся к серверу
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Создаём клиента
	client := pb.NewEchoServiceClient(conn)

	// Вызываем все методы
	OneSay(client)
	StreamHello(client)
	ClientStreamHello(client)
	BidirectionalStreamHello(client)
}
