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

	"github.com/aashabtajwar/th-server/app/tokenmanager"
)

// mimeType := http.DetectContentType(dataBuf.Bytes())

var connectedUser = make(map[string]net.Conn)

func verifyToken(token *bytes.Buffer, conn net.Conn) {
	fmt.Println("raw token data\n", token)
	stringToken := string(token.Bytes()[:])
	fmt.Println("Printing --> ", stringToken)
	claims := tokenmanager.DecodeToken(stringToken)
	user_id := claims["id"]
	connectedUser[user_id.(string)] = conn
	AddConnection(user_id.(string), conn)
	fmt.Println("New user connected: ", user_id.(string))
}

/*
--------------------------NOTES ON FILE SAVING--------------------------------------
* usually, file contents (just the changes) have to be turned into hashes and then sent by the client
* when the server receives them, it will be saved as hashes (in the min.io ideally)
* there will be many versions of the same file hash so that previous versions of the file can be restored
*/

func saveFile(fileData *bytes.Buffer, metadata map[string]string) {
	fmt.Println("data for saving file...\n", fileData.Bytes())

	// this is not an ideal way to define storage dir
	storageDir := "/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/"
	versionCarrierPath := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"] + "_currentversion.txt" // user_id should be extracted from connectedUser (or is this one okay?)
	fileDir := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"]

	db, er := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	defer db.Close()
	if er != nil {
		fmt.Println("Error when saving info on DB: \n", er)
	}

	// query workspaceId using workspace dir name
	queryId := fmt.Sprintf("SELECT workspace_id FROM workspace where name='%s'", metadata["workspace"])
	var workspaceId int
	if err := db.QueryRow(queryId).Scan(&workspaceId); err != nil {
		fmt.Println("Error Querying workspace ID", err)
	}

	fmt.Println("Now priting the workspace id: ", workspaceId)

	workspaceIdString := strconv.Itoa(workspaceId)

	metadata["workspaceId"] = workspaceIdString

	if _, err := os.Stat(versionCarrierPath); err == nil {

		fmt.Println("following here")
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
		fmt.Println("Written: ", n)

		// query file version info from db
		// for now, just querying using filename
		// but this is not a good idea
		// anyway, update the version

		query := fmt.Sprintf("SELECT version FROM workspace_files WHERE filename='%s'", fileDir)

		var versionNum int
		if er := db.QueryRow(query).Scan(&versionNum); er != nil {
			fmt.Println(er)
		}
		// fmt.Println("Now Printing the ")

		// for now sql version number and version number saved in file are different
		// version number should ONLY be saved in DATABASE, not in file
		versionNum += 1
		updateQuery := fmt.Sprintf("UPDATE workspace_files SET version=%d WHERE filename='%s'", versionNum, fileDir)
		_, err := db.Exec(updateQuery)
		if err != nil {
			fmt.Println("ERROR updating file verion\n", err)
		}
		fmt.Println("File Updated")

	} else if errors.Is(err, os.ErrNotExist) {
		// create file that carries file version number
		versionCarrierFile, er := os.Create(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		defer versionCarrierFile.Close()
		versionInString := strconv.Itoa(1)
		versionInBytes := []byte(versionInString)
		fmt.Println("came here...")

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
		insert := fmt.Sprintf("INSERT INTO workspace_files (filename, workspace_id, user_id, version) VALUES ('%s', '%s', '%s', %d)", fileDir, metadata["workspaceId"], metadata["user_id"], 1)
		res, err := db.Query(insert)
		if err != nil {
			fmt.Println("Insert Error:\n", err)
		}
		fmt.Println(res)
	}

	// check if this workspace have connected users
	// if so send the file to them

}

func CheckReceivedData(conn net.Conn, connections []net.Conn) {
	// buf := new(bytes.Buffer)
	dataBuf := new(bytes.Buffer)
	var dataString string
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
		var mappedData map[string]string

		fmt.Println("received data :\n", dataBuf.Bytes())
		if c%2 != 0 {
			// raw file data received
			// store the data in another variable
			fmt.Println("RRRR File data...")
			fileData.Write(dataBuf.Bytes())
			dataBuf.Reset()

		} else if c%2 == 0 {
			// file metadata received
			fmt.Println("Received Meta :\n", dataBuf.Bytes())
			data := dataBuf.Bytes()
			dataString = string(data[:])
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
				go BroadCastToUsers(fileData, connectedUser, mappedData, conn, dataString, connections)

				go saveFile(fileData, mappedData)
			}
			c = 0
		}
	}
}
