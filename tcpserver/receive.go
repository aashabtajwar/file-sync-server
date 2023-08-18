package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

func CheckReceivedData(conn net.Conn) {
	// buf := new(bytes.Buffer)
	dataBuf := new(bytes.Buffer)
	c := 0
	for {
		var size int64
		// read size from connection which is a binary
		// &size because it needs to read into memory
		// binary.Write(conn, binary.LittleEndian, &sizeTwo)
		// copy from connection into buf

		// n, err := io.CopyN(buf, conn, int64(5))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// receivedBytes := buf.Bytes()

		binary.Read(conn, binary.LittleEndian, &size)
		x, err := io.CopyN(dataBuf, conn, int64(size))
		if err != nil {
			log.Fatal(err)
		}
		c = c + 1
		// fmt.Println(dataBuf.Bytes())
		fmt.Printf("Received %d bytes and Count is %d\n", x, c)

		if c%2 != 0 {
			// file data received
			// store the data in another variable
			fmt.Println("this path taken")
			fmt.Printf("Count value %d\n", c)
			fileData := dataBuf
			fmt.Println(fileData)
			dataBuf.Reset()

		} else if c%2 == 0 {
			// file metadata received
			data := dataBuf.Bytes()
			dataString := string(data[:])
			// mappedData := map[string]string{}
			var mappedData map[string]string
			if err := json.Unmarshal([]byte(dataString), &mappedData); err != nil {
				fmt.Println("Error: ", err)
			}
			fmt.Println(mappedData)
			c = 0
		}

		// mimeType := http.DetectContentType(dataBuf.Bytes())
		// fmt.Println(mimeType)
	}
}

// if the uploaded data is a JWT token
func VerifyToken(conn net.Conn, buffer *bytes.Buffer, size int64) {
	binary.Read(conn, binary.LittleEndian, &size)
	_, err := io.CopyN(buffer, conn, int64(size))
	if err != nil {
		log.Fatal(err)
	}
	token := string(buffer.Bytes()[:])
	fmt.Println(token)

}

// if the uploaded data is a file
func HandleFile(conn net.Conn) {

}
