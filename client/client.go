package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// NewClient creates a new client connection to the server.
//
// serverIp: IP address of the server.
// serverPort: port number of the server.
// Returns a new *Client instance or nil if there was an error connecting
func NewClient(serverIp string, serverPort int) *Client {
	// connect the server
	url := fmt.Sprintf("%s:%d", serverIp, serverPort)
	fmt.Println(url)

	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("net dial error: ", err)
		return nil
	}

	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		conn:       conn,
		flag:       999,
	}

	return client
}

// Client menu mode
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 修改用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>  输入错误, 请输入有效的数字: 0-3 <<<<<< ")
		return false
	}

}
func (client *Client) PublicChat() {
	var msg string

	fmt.Println(">>>>>>  请输入消息: ")
	fmt.Scanln(&msg)

	for msg != "exit" {
		if len(msg) != 0 {
			_, err := client.conn.Write([]byte(msg + "\n"))

			if err != nil {
				fmt.Println("net write error: ", err)
				break
			}
		}

		msg = ""
		fmt.Println(">>>>>>  请输入消息: ")
		fmt.Scanln(&msg)
	}
}

func (client *Client) selectUser() {
	msg := "who\n"
	_, err := client.conn.Write([]byte(msg))

	if err != nil {
		fmt.Println("net write error: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var receiver string
	var msg string

	client.selectUser()

	fmt.Println("请输入聊天对象的名字, 输入exit退出")
	fmt.Scanln(&receiver)

	for receiver != "exit" {
		fmt.Println("请输入聊天内容, 输入exit退出")
		fmt.Scanln(&msg)

		for msg != "exit" {
			if len(msg) != 0 {
				message := "to|" + receiver + "|" + msg + "\n\n"
				_, err := client.conn.Write([]byte(message))

				if err != nil {
					fmt.Println("net write error: ", err)
					break
				}
			}

			msg = ""
			fmt.Println(">>>>>>  请输入消息, 输入exit退出: ")
			fmt.Scanln(&msg)
		}

		receiver = ""
		fmt.Println("请输入聊天对象的名字, 输入exit退出")
		fmt.Scanln(&receiver)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>  请输入用户名: ")
	fmt.Scanln(&client.Name)

	msg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(msg))

	if err != nil {
		fmt.Println("net write error: ", err)
		return false
	}

	return true
}

// handle server response message, and display in stdout
func (client *Client) HandleMessage() {
	// if client conn has data, immediately read it and output to stdout, block forever
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) chat() {
	for client.flag != 0 {
		for !client.menu() {
		}

		switch client.flag {
		case 1:
			// fmt.Println("公聊模式")
			client.PublicChat()
			break
		case 2:
			// fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 3:
			// fmt.Println("修改用户名")
			client.UpdateName()
			break
		case 0:
			fmt.Println("退出聊天室")
			break
		}

	}
}

// e.g. client -ip 127.0.0.1 -port 30000
var serverIp string
var serverPort int

// init args of cmd
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Please set server ip")
	flag.IntVar(&serverPort, "port", 30000, "Please set server port")
}

// main is the entry point of the program.
//
// No parameters.
// No return values.
func main() {

	// parse args
	flag.Parse()

	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println(">>>>>>  连接服务器失败 <<<<<< ")
		return
	}

	fmt.Println(">>>>>>  连接服务器成功 <<<<<< ")

	go client.HandleMessage()

	client.chat()

	select {}
}
