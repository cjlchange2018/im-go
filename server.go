package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip             string
	Port           int
	OnlineUsersMap map[string]*User
	mapLock        sync.RWMutex
	Message        chan string
}

// listen message method for server
func (this *Server) Listen() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineUsersMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// broacast method for server
func (this *Server) Broadcast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- msg
}

// handle method for Server
func (this *Server) Handle(conn net.Conn) {

	//add user to map
	user := NewUser(conn)

	this.mapLock.Lock()
	this.OnlineUsersMap[user.Name] = user
	this.mapLock.Unlock()

	//boradcast the message
	this.Broadcast(user, "is online")

	//block handeler
	select {}
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

	//start listen to messages
	go this.Listen()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Listen accrpt err: ", err)
			continue
		}

		//do handle
		go this.Handle(conn)
	}
}

// function to create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:             ip,
		Port:           port,
		OnlineUsersMap: make(map[string]*User),
		Message:        make(chan string),
	}
	return server
}
