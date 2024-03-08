package main

import (
	"bufio"
	"log/slog"
	"net"
	"os"
	"sync"
)

const (
	SERVER_PORT = "1337"
)

var (
	inbox = make(chan MsgClient)
	users = make(map[net.Conn]bool)
	mu    sync.RWMutex
)

type MsgClient struct {
	msg  string
	conn net.Conn
}

func main() {
	slog.Info("starting server...")

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+SERVER_PORT)

	if err != nil {
		slog.Error("error resolving address: " + err.Error())
		os.Exit(1)
	}

	server, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		slog.Error("error starting server: " + err.Error())
		os.Exit(1)
	}

	defer server.Close()

	go broadcast()

	slog.Info("waiting for connection...")
	for {
		conn, err := server.Accept()

		if err != nil {
			slog.Warn("error receiving connection: " + err.Error())
		}

		slog.Info("connection established")
		go connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()

	mu.Lock()
	users[conn] = true
	mu.Unlock()

	buf := bufio.NewScanner(conn)
	for buf.Scan() {

		msg := buf.Text()

		if len(msg) == 0 {
			continue
		}

		inbox <- MsgClient{
			conn: conn,
			msg:  msg,
		}
	}

	mu.Lock()
	delete(users, conn)
	mu.Unlock()
}

func broadcast() {
	for {
		select {
		case cur := <-inbox:
			for i := range users {
				if i != cur.conn {
					i.Write([]byte("msg: " + cur.msg + "\n"))
				}
			}
		}
	}
}
