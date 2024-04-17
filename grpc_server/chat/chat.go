package chat

import (
	"context"
	"log"
	sync "sync"
)

type Server struct {
	UnimplementedChatServiceServer
	Conns []*Connection
}

type Connection struct {
	UnimplementedChatServiceServer
	stream ChatService_ConnectServer
	id     string
	active bool
	error  chan error
}

func (p *Server) Connect(pconn *ConnectionRequest, stream ChatService_ConnectServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
	}
	log.Printf("Connection received: %v", conn.id)

	p.Conns = append(p.Conns, conn)
	return <-conn.error
}

func (s *Server) BroadcastMessage(ctx context.Context, in *Message) (*CloseResponse, error) {
	log.Printf("Broadcast message received: %s\n", in.Body)
	wait := sync.WaitGroup{}
	done := make(chan int)
	for _, conn := range s.Conns {
		wait.Add(1)
		go func(in *Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				if err := conn.stream.Send(in); err != nil {
					log.Printf("Error sending to client: %v", err)
					conn.active = false
					conn.error <- err
				}
			}
		}(in, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &CloseResponse{}, nil
}
