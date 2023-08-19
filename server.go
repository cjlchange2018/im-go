package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// function to create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// handler method for Server
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection Established.")
}

// start method for Server
func (this *Server) Start() {
	//socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	//close listen socket
	defer listener.Close()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Listen accrpt err: ", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
