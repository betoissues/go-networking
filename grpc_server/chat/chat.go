package chat

import (
    "log"
    "context"
)

type Server struct {
    UnimplementedChatServiceServer
}

func (s *Server) SendMessage(ctx context.Context, in *Message) (*Message, error){
    log.Printf("Message received: %s\n", in.Body)
    return &Message{Body: "Hello from the server"}, nil
}
