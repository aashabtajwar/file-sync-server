package tcpserver

import (
	"log"
	"net"
)

func Start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal()
	}
	// fmt.Println(ln)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal()
		}
		go ReceiveFiles(conn)
	}
}
