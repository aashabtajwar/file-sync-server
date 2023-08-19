package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
)

func verifyToken(token *bytes.Buffer) {
	stringToken := string(token.Bytes()[:])
	claims := tokenmanager.DecodeToken(stringToken)
	user_id := claims["id"]
	fmt.Println("User ID", user_id)

	// fmt.Println("JWT: ", stringToken)
}

func saveFile(fileData *bytes.Buffer, metadata map[string]string) {
	fmt.Println("Printing metadata value", metadata["key1"])

}

func CheckReceivedData(conn net.Conn) {
	// buf := new(bytes.Buffer)
	dataBuf := new(bytes.Buffer)

	fileData := new(bytes.Buffer)
	c := 0
	iter := 0
	for {
		var size int64
		// read size from connection which is a binary
		// &size because it needs to read into memory
		// binary.Write(conn, binary.LittleEndian, &sizeTwo)

		iter += 1
		fmt.Println("Iteration Number: ", iter)
		binary.Read(conn, binary.LittleEndian, &size)
		x, err := io.CopyN(dataBuf, conn, int64(size))
		if err != nil {
			log.Fatal(err)
		}
		c = c + 1
		fmt.Printf("Received %d bytes and Count is %d\n", x, c)
		fmt.Println(dataBuf)
		// var fileData *bytes.Bufferdata
		var mappedData map[string]string

		if c%2 != 0 {
			// file data received
			// store the data in another variable
			fmt.Println("Handling raw data")
			fileData.Write(dataBuf.Bytes())
			dataBuf.Reset()

		} else if c%2 == 0 {
			// file metadata received
			fmt.Println("Handling Metadata")
			data := dataBuf.Bytes()
			dataString := string(data[:])
			if err := json.Unmarshal([]byte(dataString), &mappedData); err != nil {
				fmt.Println("Error: ", err)
			}
			fmt.Println(mappedData)
			dataBuf.Reset()
		}
		// mimeType := http.DetectContentType(dataBuf.Bytes())
		if c == 2 {
			if mappedData["type"] == "token" {
				go verifyToken(fileData)
			} else if mappedData["type"] == "file" {
				go saveFile(fileData, mappedData)
			}
			c = 0
		}
	}
}

// if the uploaded data is a JWT token
// func VerifyToken(conn net.Conn, buffer *bytes.Buffer, size int64) {
// 	binary.Read(conn, binary.LittleEndian, &size)
// 	_, err := io.CopyN(buffer, conn, int64(size))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	token := string(buffer.Bytes()[:])
// 	fmt.Println(token)

// }

// if the uploaded data is a file
func HandleFile(conn net.Conn) {

}
