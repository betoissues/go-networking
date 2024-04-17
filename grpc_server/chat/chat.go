package chat;

import (
    "log"
    "context"
    "sync"
)

type Server struct {
    UnimplementedChatServiceServer
    Conns []*Connection
}

type Connection struct {
    UnimplementedChatServiceServer
    stream ChatService_ConnectServer
    id string
    active bool
    error chan error
}

func (p *Server) Connect(pconn *ConnectionRequest, stream ChatService_ConnectServer) (error) {
    conn := &Connection{
        stream: stream,
        id: pconn.User.Id,
        active: true,
    }
    log.Printf("Connection received: %v", conn.id)

    p.Conns = append(p.Conns, conn)
    return <-conn.error
}

func (s *Server) SendMessage(ctx context.Context, in *Message) (*Message, error){
    log.Printf("Message received: %s\n", in.Body)
    return &Message{Body: "Hello from the server"}, nil
}

func (s *Server) BroadcastMessage(ctx context.Context, in *Message) (*Message, error){
    log.Printf("Broadcast message received: %s\n", in.Body)
    wait := sync.WaitGroup{}
    return &Message{Body: "Broadcast message received"}, nil
}
