package main

import "net"

// User represents a connected user.
type User struct {
	Name   string      // User's name
	Addr   string      // User's network address
	C      chan string // Channel for communication with the user
	conn   net.Conn    // User's network connection
	server *Server     // Server that user associated with
}

// Listen listens for messages on the user's channel and sends them to their connection.
func (user *User) Listen() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}

// NewUser creates a new User instance.
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// Start listening to the user's channel
	go user.Listen()

	return user
}

func (user *User) Login() {
	// Add user to the map
	user.server.mapLock.Lock()
	user.server.OnlineUsersMap[user.Name] = user
	user.server.mapLock.Unlock()

	// Broadcast the "is online" message
	user.server.Broadcast(user, "is online")
}

func (user *User) Logout() {
	// Remove user from the map
	user.server.mapLock.Lock()
	delete(user.server.OnlineUsersMap, user.Name)
	user.server.mapLock.Unlock()

	// Broadcast the "is online" message
	user.server.Broadcast(user, "is offline")
}

func (user *User) SendMessage(msg string) {
	user.server.Broadcast(user, msg)
}
