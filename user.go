package main

import (
	"net"
	"strings"
)

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

// User handle message
func (user *User) HandleMessage(msg string) {
	if msg == "who" {
		user.server.lock.Lock()
		for _, u := range user.server.Users {
			msg := "[" + u.Addr + "]" + u.Name + ":" + "在线中\n"
			user.SendMessage(msg)
		}
		user.server.lock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// rename format like: rename|James
		name := strings.Split(msg, "|")[1]

		// check name exists
		_, ok := user.server.Users[name]
		if ok {
			user.SendMessage("用户名已存在\n")
		} else {
			server := user.server
			server.lock.Lock()
			delete(server.Users, user.Name)
			server.Users[name] = user
			user.server.lock.Unlock()

			user.Name = name
			user.SendMessage("用户名成功修改为：" + name + "\n")
		}
	} else {
		user.server.Broadcast(user, msg)
	}
}

// User send message
func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg))
}
