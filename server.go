package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

// Create server 
func NewServer(ip string, port int) *Server {
	server := & Server {
		Ip: ip,
		Port: port,
	}

	return server
}

// Define connection handle
func (this *Server) Handle(conn net.Conn) {
	fmt.Println("Server handling")
}

// Start server
func (this *Server) Start() {
	// Listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net listen error: ", err)

		return
	}

	defer listener.Close()

	// Handle for each connection
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Listener accept error: ", err)

			continue
		}

		go this.Handle(conn)

	}
}

