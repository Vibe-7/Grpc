Echo Service gRPC Pet Project

Описание:
Этот репозиторий демонстрирует реализацию простого gRPC-сервера и клиента на Go, поддерживающих четыре типа взаимодействия:

**Unary RPC** – классический запрос-ответ (SayHello).

**Server Streaming RPC** – сервер отправляет серию сообщений в ответ на один запрос клиента (SayHelloStream).

**Client Streaming RPC** – клиент посылает поток сообщений, после чего сервер возвращает единичный ответ (SayHelloClientStream).

**Bidirectional Streaming RPC** – двусторонний поток сообщений между клиентом и сервером (SayHelloBidirectional).

Технологии и зависимости

**Go** (Golang) (версия 1.20+)

**gRPC**: google.golang.org/grpc

**Protocol Buffers**: syntax = "proto3", пакет protoc-gen-go для генерации Go-кода

Контекст и таймауты: **context.WithTimeout** для защиты от зависаний

Безопасность соединения: **grpc.WithTransportCredentials(insecure.NewCredentials())** (для локальной разработки)

Логирование: стандартный пакет **log**

Сетевое взаимодействие: **net.Listen** для порта сервера

**Так же я добавил вторую ветку, где разделил сервер и cmd**

