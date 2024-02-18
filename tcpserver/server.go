package tcpserver

import (
	"fmt"
	"log"
	"net"
)

var server []net.Listener
var activeConnections []net.Conn

var conns []net.Conn

func updatedConnections() []net.Conn {
	return conns
}

func SetupConn() net.Listener {
	if len(server) == 1 {
		return server[0]
	}
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Turning on TCP Server\n", ln)
	server = append(server, ln)
	return ln
}

func Start() {
	// ln, err := net.Listen("tcp", ":3030")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// msg := make(chan string)
	ln := SetupConn()
	for {
		conn, err := ln.Accept()
		activeConnections = append(activeConnections, conn)
		conns = append(conns, conn)
		for _, e := range activeConnections {
			fmt.Println("connection: ", e)
		}
		fmt.Println("New Connection: ", conn)
		fmt.Println("CC = ", connections)
		if err != nil {
			log.Fatal(err)
		}
		go CheckReceivedData(conn, activeConnections)
	}
}
