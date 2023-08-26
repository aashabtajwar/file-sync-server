package tcpserver

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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
	fileDir := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"]

	if _, err := os.Stat(versionCarrierPath); err == nil {
		data, er := os.ReadFile(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		versionInString := string(data[:])
		version, er := strconv.Atoi(versionInString)
		version += 1
		newFileName := fileDir + "_" + strconv.Itoa(version) + "." + metadata["mimetype"]

		// save updated file version number
		f, er := os.OpenFile(versionCarrierPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if er != nil {
			log.Fatal(er)
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(version))
		// save updated file
		file, er := os.Create(newFileName)
		if er != nil {
			log.Fatal("Error Creating Updated File: ", er)
		}
		defer file.Close()

		n, er := file.Write(fileData.Bytes())
		if er != nil {
			log.Fatal(er)
		}
		fmt.Println("Writed: ", n)

	} else if errors.Is(err, os.ErrNotExist) {
		// create file that carries file version number
		versionCarrierFile, er := os.Create(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		defer versionCarrierFile.Close()
		versionInString := strconv.Itoa(1)
		versionInBytes := []byte(versionInString)

		_, e := versionCarrierFile.Write(versionInBytes)
		if er != nil {
			log.Fatal(e)
		}
		filePath := fileDir + "_1." + metadata["mimetype"]
		file, er := os.Create(filePath)
		if er != nil {
			log.Fatal("File Saving Error: ", er)
		}
		defer file.Close()

		_, errr := file.Write(fileData.Bytes())
		if errr != nil {
			log.Fatal(er)
		}

		// add to database
		db, er := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()
		if er != nil {
			fmt.Println("Error when saving info on DB: \n", err)
		}
		insert := fmt.Sprintf("INSERT INTO workspace_files (filename, workspace_id, user_id, version) VALUES ('%s', '%s', '%s', %d)", fileDir, metadata["workspace_id"], metadata["user_id"], 1)
		res, err := db.Query(insert)
		if err != nil {
			fmt.Println("Insert Error:\n", err)
		}
		fmt.Println(res)
	}

	// check if this workspace have connected users
	// if so send the file to them

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
		binary.Read(conn, binary.LittleEndian, &size)
		_, err := io.CopyN(dataBuf, conn, int64(size))
		if err != nil {
			log.Fatal(err)
		}
		c = c + 1
		fmt.Println(dataBuf)
		var mappedData map[string]string

		if c%2 != 0 {
			// raw file data received
			// store the data in another variable
			fileData.Write(dataBuf.Bytes())
			dataBuf.Reset()

		} else if c%2 == 0 {
			// file metadata received
			data := dataBuf.Bytes()
			dataString := string(data[:])
			if err := json.Unmarshal([]byte(dataString), &mappedData); err != nil {
				fmt.Println("Error: ", err)
			}
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
