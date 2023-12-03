package tcpserver

/*

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
				binary.Write(co)

			}

		}
	}
}

*/
