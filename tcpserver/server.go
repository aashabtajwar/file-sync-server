package tcpserver

import (
	"fmt"
	"log"
	"net"
)

var server []net.Listener
var activeConnections []net.Conn

func SetupConn() net.Listener {
	if len(server) == 1 {
		return server[0]
	}
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	server = append(server, ln)
	return ln
}

func Start() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		activeConnections = append(activeConnections, conn)
		for _, e := range activeConnections {
			fmt.Println("connection: ", e)
		}
		fmt.Println("New Connection: ", conn)
		if err != nil {
			log.Fatal(err)
		}
		go CheckReceivedData(conn, activeConnections)
	}
}

func readStreamLoop(conn net.Conn) {

}
