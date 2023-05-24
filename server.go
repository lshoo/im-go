package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int
	// Online user map
	Users map[string] *User
	lock sync.RWMutex

	// Message channel
	Message chan string
}

// Create server 
func NewServer(ip string, port int) *Server {
	server := & Server {
		Ip: ip,
		Port: port,
		Users: make(map[string] *User),
		Message: make(chan string),
	}

	return server
}

// Listen Message goroutine of server, and response to all online users
func (server *Server) Listen() {
	for {
		msg := <- server.Message

		server.lock.Lock()
		
		for _, user := range server.Users {
			user.Channel <- msg
		}

		server.lock.Unlock()
	}
}
// Server broadcast Message
func (server *Server) Broadcast(user *User, msg string) {
	message := "[" + user.Addr + "]" + user.Name + ":" + msg
	
	server.Message <- message
}

// Define connection handle
func (server *Server) Handle(conn net.Conn) {
	fmt.Println("Server handling", server)

	user := NewUser(conn)

	// Add User to online map
	server.lock.Lock()
	server.Users[user.Name] = user
	server.lock.Unlock()

	// broadcast user login message to online map
	server.Broadcast(user, "已上线")

	// blocking 
	select {

	}
}

// Start server
func (server *Server) Start() {
	// Listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))

	if err != nil {
		fmt.Println("net listen error: ", err)

		return
	}

	defer listener.Close()

	// Start listen user goroutine
	go server.Listen()

	// Handle for each connection
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Listener accept error: ", err)

			continue
		}

		go server.Handle(conn)

	}
}

