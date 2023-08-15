package tcpserver

import (
	"log"
	"net"
)

func Start() {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go CheckReceivedData(conn)
	}
}

func readStreamLoop(conn net.Conn) {

}
