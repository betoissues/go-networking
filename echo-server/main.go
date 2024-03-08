package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1337"
)

func main() {
	fmt.Println("starting server...")

	// Alternative `net.ListenTCP` requires `TCPListener`
	server, err := net.Listen("tcp", SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("error listening: ", err.Error())
		os.Exit(1)
	}

	// defer to makes sure connection is closed
	defer server.Close()

	fmt.Println("waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("client connected")
		go echo(connection)
	}
}

func echo(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("error receiving: ", err.Error())
			return
		}
		fmt.Println("recv: ", string(buffer[:mLen]))
		_, err = connection.Write([]byte("echo: " + string(buffer[:mLen])))

		if err != nil {
			fmt.Println("error sending: ", err.Error())
			return
		}
	}
}
