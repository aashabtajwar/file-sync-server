package tcpserver

import (
	"fmt"
	"log"
	"net"
)

// var activeConnections map[string]net.Conn
var activeConnections []net.Conn

func Start() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		activeConnections = append(activeConnections, conn)
		fmt.Println("Printing Connections ", activeConnections)
		if err != nil {
			log.Fatal(err)
		}
		go CheckReceivedData(conn)
	}
}

func readStreamLoop(conn net.Conn) {

}
