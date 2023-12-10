package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func SendFiles(workspaceName string, workspaceId string, user_id string) {
	// send files according to the workspace
	// loop over the files and check if they contain the correct workspace_workspaceid in their names
	// the ones that do, send those files

	time.Sleep(100 * time.Millisecond)
	check := workspaceName + "_" + workspaceName // check var is the workspace_workspaceid to check in file name
	// ln := SetupConn()
	entries, err := os.ReadDir("./storage/")
	if err != nil {
		log.Fatal(err)
	}
	conn := ReturnConnection(user_id)
	for _, e := range entries {
		if strings.Contains(e.Name(), check) {
			// make sure to only send .go files
			splitted := strings.Split(e.Name(), ".")
			if (splitted[len(splitted)-1]) == ".go" {
				pwd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				filePath := pwd + "/" + e.Name()
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Println("Error opening file\n", err)
				}
				fi, err := file.Stat()
				if err != nil {
					fmt.Println("Error getting File Stat\n", err)
				}
				byteData, er := os.ReadFile(filePath)
				if er != nil {
					fmt.Println("Error reading file\n", er)
				}

				metaDataString := fmt.Sprintf(`
					{
						"workspace": "%s",
						"filename": "%s",
						"mimetype": "%s",
						"type": "file",
						"name": "%s"
					}
				`, workspaceName, e.Name(), splitted[len(splitted)-1], e.Name())

				metaDataBytes := []byte(metaDataString)
				binary.Write(conn, binary.LittleEndian, int64(fi.Size()))
				n1, err := io.CopyN(conn, bytes.NewReader(byteData), int64(fi.Size()))
				if err != nil {
					fmt.Println("Error Sending File data\n", err)
				}
				fmt.Printf("Written %d bytes\n", n1)
				time.Sleep(100 * time.Millisecond)

				// send metadata
				binary.Write(conn, binary.BigEndian, int64(len(metaDataBytes)))
				n2, err := io.CopyN(conn, bytes.NewReader(metaDataBytes), int64(len(metaDataBytes)))
				if err != nil {
					fmt.Println("Error sending file metadatra\n", err)
				}
				fmt.Printf("Written %d bytes\n", n2)

			}

		}
	}
}
