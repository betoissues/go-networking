package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/betoissues/go-networking/grpc_server/chat"
)

func main() {
    lis, err := net.Listen("tcp", ":8080")

    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := chat.Server{
        Conns: []*chat.Connection{},
    }

    grpcServer := grpc.NewServer()

    chat.RegisterChatServiceServer(grpcServer, &s)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
