package tcpserver

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func SendFiles(workspaceName string, workspaceId string) {
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
				fmt.Println(pwd)
			}

		}
	}
}
