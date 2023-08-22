package tcpserver

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// mimeType := http.DetectContentType(dataBuf.Bytes())

var connectedUser map[net.Conn]string

func verifyToken(token *bytes.Buffer, conn net.Conn) {
	// stringToken := string(token.Bytes()[:])
	// claims := tokenmanager.DecodeToken(stringToken)
	// user_id := claims["id"]
	// connectedUser[conn] = user_id.(string)
}

/*
--------------------------NOTES ON FILE SAVING--------------------------------------
* usually, file contents (just the changes) have to be turned into hashes and then sent by the client
* when the server receives them, it will be saved as hashes (in the min.io ideally)
* there will be many versions of the same file hash so that previous versions of the file can be restored
*/
func saveFile(fileData *bytes.Buffer, metadata map[string]string) {
	// this is not an ideal way to define storage dir
	storageDir := "/home/aashab/code/src/github.com/aashabtajwar/th-server/storage/"
	versionCarrierPath := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"] + "_currentversion.txt" // user_id should be extracted from connectedUser (or is this one okay?)
	fmt.Println(versionCarrierPath)
	if _, err := os.Stat(versionCarrierPath); err == nil {
		data, er := os.ReadFile(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		version := binary.BigEndian.Uint32(data)
		fmt.Println(version)

	} else if errors.Is(err, os.ErrNotExist) {
		versionCarrierFile, er := os.Create(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		defer versionCarrierFile.Close()
		version := make([]byte, 4)
		binary.BigEndian.PutUint32(version, 1)
		_, e := versionCarrierFile.Write(version)
		if er != nil {
			log.Fatal(e)
		}
	}

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

		if c == 2 {
			/*
				--------------------SIDE NOTES------------------------
				* here, maybe channels should be used instead of go verifiedToken or go saveFile
				* properly learn channels
			*/
			if mappedData["type"] == "token" {
				go verifyToken(fileData, conn)
			} else if mappedData["type"] == "file" {
				go saveFile(fileData, mappedData)
			}
			c = 0
		}
	}
}
