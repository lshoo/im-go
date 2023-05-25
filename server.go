package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// Online user map
	Users map[string]*User
	lock  sync.RWMutex

	// Message channel
	Message chan string
}

// Create server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:      ip,
		Port:    port,
		Users:   make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// Step 1: Start Server, and server always listens connection
// Step 2: Client connected the server(ip and port)
// Step 3: Server listened the connection, and handel connection to create user, add user to online map
// Step 4: Server write message to Message channel
// Step 5: Server broadcast messge in Message channel to online user
// Step 6: Server handle client message and broadcast

// Step 1: Start server
func (server *Server) Start() {
	// Listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))

	if err != nil {
		fmt.Println("net listen error: ", err)

		return
	}

	defer listener.Close()

	// Start listen user goroutine
	go server.listen()

	// Handle for each connection
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Listener accept error: ", err)

			continue
		}

		go server.handle(conn)

	}
}

// Step 3: Define connection handle
func (server *Server) handle(conn net.Conn) {
	fmt.Println("Server handling", server)

	user := NewUser(conn, server)

	user.Online()

	// current user isLive
	isLive := make(chan bool)

	// handle client message
	go server.handleMessage(user, isLive)

	// blocking
	for {
		select {
		case <-isLive:
			// Do nothing, active the select timer

		case <-time.After(time.Second * 100):
			// timeout, close the connection
			user.SendMessage("超时,已下线 \n")

			close(user.Channel)

			conn.Close()

			return // runtime.Goexit()
		}
	}
}

// Step 4: Server broadcast connected Message
func (server *Server) Broadcast(user *User, msg string) {
	message := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- message
}

// Step 5: listen Message goroutine of server, and response to all online users
func (server *Server) listen() {
	for {
		msg := <-server.Message

		server.lock.Lock()

		for _, user := range server.Users {
			user.Channel <- msg
		}

		server.lock.Unlock()
	}
}

// Step 6: server handle client user message and broadcast to all online users
func (server *Server) handleMessage(user *User, channel chan bool) {

	buf := make([]byte, 4096)

	for {
		n, err := user.conn.Read(buf)

		if n == 0 {
			// server.broadcastConnected(user, "下线")
			user.Offline()
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("Conn Read Error!", err)
			return
		}

		// Extract user message from connection, and discard last char(\n)
		msg := string(buf[:n-1])

		fmt.Println(msg)

		// broadcast message
		// server.Broadcast(user, msg)
		user.HandleMessage(msg)

		// show the user is live if any message
		channel <- true
	}

}
