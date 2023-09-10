package tcpserver

import (
	"fmt"
	"log"
	"net"
)

var activeConnections []net.Conn

func Start() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		activeConnections = append(activeConnections, conn)
		fmt.Println("New Connection: ", conn)
		if err != nil {
			log.Fatal(err)
		}
		go CheckReceivedData(conn, activeConnections)
	}
}

func readStreamLoop(conn net.Conn) {

}
