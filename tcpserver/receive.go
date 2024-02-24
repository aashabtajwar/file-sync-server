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
	"strings"
	"time"

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

func deleteFile(workspaceName string, fileName string) {
	fileNameSplitted := strings.Split(fileName, ".")
	fileName = fileNameSplitted[0]
	fmt.Println("delete")
	check := workspaceName + "_" + fileName
	fmt.Println("pattern = ", check)
	entries, err := os.ReadDir("./storage/")
	if err != nil {
		fmt.Println("Error Reading Directory\n", err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), check) {
			// delete file
			er := os.Remove("./storage/" + e.Name())
			if er != nil {
				fmt.Println("Error Removing File")
			}
		}
	}
}

/*
--------------------------NOTES ON FILE SAVING--------------------------------------
* usually, file contents (just the changes) have to be turned into hashes and then sent by the client
* when the server receives them, it will be saved as hashes (in the min.io ideally)
* there will be many versions of the same file hash so that previous versions of the file can be restored
*/

func saveFile(fileData *bytes.Buffer, metadata map[string]string) {
	// fmt.Println("data for saving file...\n", fileData.Bytes())

	// this is not an ideal way to define storage dir
	storageDir := "/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/"

	// version carrier path needs to be updated
	// to have a better name
	versionCarrierPath := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"] + "_currentversion.txt" // user_id should be extracted from connectedUser (or is this one okay?)
	// fileDir := storageDir + metadata["workspace"] + "_" + metadata["user_id"] + "_" + metadata["name"]
	fileDir := storageDir + metadata["workspace"] + "_" + metadata["name"]
	db, er := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	if er != nil {
		fmt.Println("Error when saving info on DB: \n", er)
	}
	defer db.Close()

	// fmt.Println("FOLDER = ", metadata["workspace"])

	// query workspaceId using workspace dir name
	queryId := fmt.Sprintf("SELECT workspace_id FROM workspace where name='%s'", metadata["workspace"])
	var workspaceId int
	if err := db.QueryRow(queryId).Scan(&workspaceId); err != nil {
		fmt.Println("Error Querying workspace ID", err)
	}

	fmt.Println("Now priting the workspace id: ", workspaceId)

	workspaceIdString := strconv.Itoa(workspaceId)

	metadata["workspaceId"] = workspaceIdString

	// file exists
	if _, err := os.Stat(versionCarrierPath); err == nil {
		data, er := os.ReadFile(versionCarrierPath)
		if er != nil {
			log.Fatal(er)
		}
		versionInString := string(data[:])
		version, er := strconv.Atoi(versionInString)
		version += 1
		splittedFullName := strings.Split(fileDir, "/")
		justTheFileName := splittedFullName[len(splittedFullName)-1]
		rearragedFileName := ""
		for i, e := range splittedFullName {
			if i != len(splittedFullName)-1 {
				rearragedFileName += e
				rearragedFileName += "/"
			}
		}
		splittedJustFileName := strings.Split(justTheFileName, ".")

		currentTime := time.Now().Unix()
		newFileName := rearragedFileName + splittedJustFileName[0] + "_" + metadata["user_id"] + "_" + strconv.Itoa(version) + "_" + strconv.FormatInt(currentTime, 10) + "_." + splittedJustFileName[1]

		// newFileName := fileDir + "_" + strconv.Itoa(version) + "." + metadata["mimetype"]

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

		_, er = file.Write(fileData.Bytes())
		if er != nil {
			log.Fatal(er)
		}
		// fmt.Println("Written: ", n)

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
		// fmt.Println("came here...")

		_, e := versionCarrierFile.Write(versionInBytes)
		if er != nil {
			log.Fatal(e)
		}

		spliitedFullName := strings.Split(fileDir, "/")
		justTheFileName := spliitedFullName[len(spliitedFullName)-1] // gave this var name because usually the file name consists of the whole dir with it -_- (ik I have to change it)
		rearrangedFileName := ""
		for i, e := range spliitedFullName {
			if i != len(spliitedFullName)-1 {
				rearrangedFileName += e
				rearrangedFileName += "/"
			}
		}
		splittedJustFileName := strings.Split(justTheFileName, ".") // separate mimetype for now
		currentTime := time.Now().Unix()
		newFileName := splittedJustFileName[0] + "_" + metadata["user_id"] + "_1_" + strconv.FormatInt(currentTime, 10) + "_." + splittedJustFileName[1]
		filePath := rearrangedFileName + newFileName
		// filePath := splittedFileDir[0] + "_v1" + splittedFileDir[1]

		// filePath := fileDir + "_1." + metadata["mimetype"]
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

func (s *TcpFServer) CheckReceivedData(conn net.Conn, connections []net.Conn) {
	// buf := new(bytes.Buffer)
	fmt.Println("CURRENT CONNECTIONS\n", connections)
	dataBuf := new(bytes.Buffer)
	// newDataBuf := make(chan bytes.Buffer)
	var dataString string
	fileData := new(bytes.Buffer)
	c := 0
	iter := 0
	fmt.Println("Update 2 = ", connections)

	// Try:
ReadLoop:
	for {
		select {
		case <-s.quit:
			fmt.Println("ending...")
			return
		default:
			conn.SetDeadline(time.Now().Add(200 * time.Second))
			var size int64
			// read size from connection which is a binary
			// &size because it needs to read into memory
			// binary.Write(conn, binary.LittleEndian, &sizeTwo)
			// continue Try
			iter += 1
			binary.Read(conn, binary.LittleEndian, &size)
			n, err := io.CopyN(dataBuf, conn, int64(size))
			if err != nil {
				// log.Fatal(err)
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ReadLoop

				} else if err != io.EOF {
					log.Println("read error", err)
					return
				} else {
					log.Println(err)
					fmt.Println("unknown error...")
					return
				}
			}
			if n == 0 {
				return
			}
			c = c + 1
			var mappedData map[string]string

			// fmt.Println("received data from receiver :\n", dataBuf.Bytes())
			if c%2 != 0 {
				fmt.Println("Update 4.1 = ", connections)
				// raw file data received
				// store the data in another variable
				// fmt.Println("RRRR File data...")
				fileData.Write(dataBuf.Bytes())
				dataBuf.Reset()

			} else if c%2 == 0 {
				fmt.Println("Update 4.2 = ", connections)
				// file metadata received
				// fmt.Println("Received Meta :\n", dataBuf.Bytes())
				data := dataBuf.Bytes()
				dataString = string(data[:])
				if err := json.Unmarshal([]byte(dataString), &mappedData); err != nil {
					fmt.Println("Error: ", err)
				}
				// fmt.Println(mappedData)
				dataBuf.Reset()
			}

			if c == 2 {
				fmt.Println("Update 5 = ", connections)
				/*
					--------------------SIDE NOTES------------------------
					* here, maybe channels should be used instead of go verifiedToken or go saveFile
					* properly learn channels
				*/
				if mappedData["type"] == "token" {
					go verifyToken(fileData, conn)
				} else if mappedData["type"] == "file" {
					BroadCastToUsers(fileData, connectedUser, mappedData, conn, dataString, conns)
					if mappedData["isDeleted"] == "Yes" {
						deleteFile(mappedData["workspace"], mappedData["fileName"])
					} else {
						conns = updatedConnections()
						go saveFile(fileData, mappedData)
					}
				}
				c = 0
				continue ReadLoop
			}

		}

	}
}
