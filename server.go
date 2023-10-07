package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Server represents a chat server.
type Server struct {
	Ip             string
	Port           int
	OnlineUsersMap map[string]*User
	mapLock        sync.RWMutex
	Message        chan string
}

// Listen continuously listens for incoming messages and broadcasts them to all users.
func (server *Server) Listen() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		for _, client := range server.OnlineUsersMap {
			client.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// Broadcast sends a message to all connected users.
func (server *Server) Broadcast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- msg
}

// Handle manages a user's connection and broadcasts their messages.
func (server *Server) Handle(conn net.Conn) {
	user := NewUser(conn, server)
	user.Login()

	// Goroutine to broadcast messages
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			// Connection is closed
			if n == 0 {
				user.Logout()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err", err)
				return
			}

			// Remove newline character and broadcast message
			msg := string(buf[:n-1])
			user.SendMessage(msg)

			// Check for logout command
			if msg == "/logout" {
				user.Logout()
				return
			}
		}
	}()

	// Block indefinitely
	select {}
}

// Start starts the server and listens for incoming connections.
func (server *Server) Start() {
	// Socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	//Close listen socket
	defer listener.Close()

	// Start listening to messages
	go server.Listen()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Listen accept err: ", err)
			continue
		}

		// Handle the connection in a goroutine
		go server.Handle(conn)
	}
}

// NewServer creates a new chat server instance.
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:             ip,
		Port:           port,
		OnlineUsersMap: make(map[string]*User),
		Message:        make(chan string),
	}
	return server
}
