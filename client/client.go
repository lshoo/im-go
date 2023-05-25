package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
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
	}

	return client
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

	select {}
}
