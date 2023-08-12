package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func ReceiveFiles(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		// read size from connection which is a binary
		// &size because it needs to read into memory
		binary.Read(conn, binary.LittleEndian, &size)
		// copy from connection into buf
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over the network\n", n)
	}
}
