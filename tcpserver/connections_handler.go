package tcpserver

import (
	"fmt"
	"net"
)

// var connections []net.Conn

var connections = make(map[string]net.Conn)

func AddConnection(user_id string, conn net.Conn) {
	connections[user_id] = conn
}

func ReturnConnection(user_id string) net.Conn {
	connection, found := connections[user_id]
	if found {
		fmt.Println("Found Connection = ", connection)
		return connections[user_id]
	}
	return nil
}
