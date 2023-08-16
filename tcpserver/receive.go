package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// first check if the data received is a file upload or a JWT token
func CheckReceivedData(conn net.Conn) {
	buf := new(bytes.Buffer)
	dataBuf := new(bytes.Buffer)
	for {
		var size int64
		// var sizeTwo int64
		// read size from connection which is a binary
		// &size because it needs to read into memory
		// binary.Read(conn, binary.LittleEndian, &size)
		// time.Sleep(100 * time.Millisecond)
		// binary.Write(conn, binary.LittleEndian, &sizeTwo)
		// copy from connection into buf
		// fmt.Println("Received size ")
		// fmt.Println(size)
		// fmt.Println("Second size value")
		// fmt.Println(sizeTwo)
		n, err := io.CopyN(buf, conn, int64(3))
		if err != nil {
			log.Fatal(err)
		}
		receivedBytes := buf.Bytes()
		fmt.Println(string(receivedBytes[:]))
		fmt.Printf("Received %d bytes over the network\n", n)
		binary.Read(conn, binary.LittleEndian, &size)
		// time.Sleep(100 * time.Millisecond)
		x, err := io.CopyN(dataBuf, conn, int64(size))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dataBuf.Bytes())
		fmt.Printf("Received %d bytes", x)
	}
}

// if the uploaded data is a JWT token
func VerifyToken(conn net.Conn) {

}

// if the uploaded data is a file
func HandleFile(conn net.Conn) {

}
