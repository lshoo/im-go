package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
}

// Create User
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	// TODO listen message
	// go user.ListenMessage()

	return user
}

// User Listening channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel

		user.conn.Write([]byte(msg + "\n"))
	}
}

// User online
func (user *User) Online() {
	server := user.server

	// Add User to online map
	server.lock.Lock()
	server.Users[user.Name] = user
	server.lock.Unlock()

	// listen message
	go user.ListenMessage()

	// broadcast user login message to online map
	server.Broadcast(user, "已上线")
}

// User offline
func (user *User) Offline() {
	server := user.server

	// Add User to online map
	server.lock.Lock()
	delete(server.Users, user.Name)
	server.lock.Unlock()

	// broadcast user logout message to online map
	server.Broadcast(user, "已下线")
}

// User send message
func (user *User) BroadcastMessage(msg string) {
	user.server.Broadcast(user, msg)
}
