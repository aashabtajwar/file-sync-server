package tcpserver

import "net"

// var connections []net.Conn

var connections map[string]net.Conn

func AddConnection(user_id string, conn net.Conn) {
	connections[user_id] = conn
}

func ReturnConnection(user_id string) net.Conn {
	return connections[user_id]
}
