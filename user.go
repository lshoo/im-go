package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

// Create User
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
	}

	// TODO listen message
	go user.ListenMessage()

	return user
}

// User Listening channel
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel

		user.conn.Write([]byte(msg + "\n"))
	}
}
